package qr_codes

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"os"
)

func QrCodeGen(t string, filename string) error {
	// Create the barcode
	qrCode, _ := qr.Encode(t, qr.M, qr.Auto)

	// Scale the barcode to 200x200 pixels
	qrCode, _ = barcode.Scale(qrCode, 2000, 2000)

	// Save file to different directory
	fileDir := fmt.Sprintf("../samples/images/qr/%v", filename)

	// create the output file
	file, err := os.Create(fileDir)
	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)

	return err
}
