package qr_codes

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"os"
)

// GenerateQRCodeWithIcon generate QR Code with icon in the center position.
func GenerateQRCodeWithIcon(data string, iconName string, fileName string) (string, error) {
	iconPath := fmt.Sprintf("../assets/%v", iconName)
	filePath := fmt.Sprintf("../samples/images/qr/%v", fileName)

	// Create a new QR code barcode with the given data
	qrCode, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}

	// Scale the barcode to the desired size
	qrCode, err = barcode.Scale(qrCode, 125, 125)
	if err != nil {
		return "", err
	}

	// Load the icon image
	iconFile, err := os.Open(iconPath)
	if err != nil {
		return "", err
	}
	defer iconFile.Close()

	iconImg, _, err := image.Decode(iconFile)
	if err != nil {
		return "", err
	}

	resizeIcon := image.NewRGBA(image.Rect(0, 0, 30, 30))

	draw.CatmullRom.Scale(resizeIcon, resizeIcon.Bounds(), iconImg, iconImg.Bounds(), draw.Over, nil)

	// Create a new image with transparent background
	finalImg := image.NewRGBA(qrCode.Bounds())

	// Calculate the position to place the icon in the center of the QR code
	iconX := (qrCode.Bounds().Max.X - resizeIcon.Bounds().Max.X) / 2
	iconY := (qrCode.Bounds().Max.Y - resizeIcon.Bounds().Max.Y) / 2

	// Draw the QR code onto the final image
	draw.Draw(finalImg, qrCode.Bounds().Add(image.Point{}), qrCode, image.Point{}, draw.Over)

	// Draw the icon onto the final image
	draw.Draw(finalImg, resizeIcon.Bounds().Add(image.Pt(iconX, iconY)), resizeIcon, image.Point{}, draw.Over)

	// Create a new file to save the QR code image with the icon
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Save the final image as a PNG file
	err = png.Encode(file, finalImg)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
