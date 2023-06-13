package pdf

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"image"
	"log"
	"os"
	"os/exec"
	entity "pdf-generator-with-qr/entities"
	util "pdf-generator-with-qr/utils"
	"reflect"
)

// Option is a function type used for applying options to a PDFProcessor.
// It takes a pointer to a PDFProcessor and modifies it according to the provided option.
type Option func(*PdfProcessor)

// PdfProcessor provides operations related to PDF files.
type PdfProcessor struct {
	FilePath      string
	Base64Output  string
	PDFProtection bool
	*entity.OptionFilePDF
	*entity.OptionMetadataPDF
}

// NewPDFGopher constructor to retrieve struct PDFProcessor
func NewPDFGopher(filePath string, options ...Option) *PdfProcessor {
	option := &PdfProcessor{
		FilePath: filePath,
		OptionFilePDF: &entity.OptionFilePDF{
			StampPosition: "br",
		},
		OptionMetadataPDF: &entity.OptionMetadataPDF{},
	}

	for _, opt := range options {
		opt(option)
	}

	return option
}

// WithOptionMetadataPDF returns an Option function that sets the OptionMetadataPDF value.
func WithOptionMetadataPDF(value entity.OptionMetadataPDF) Option {
	return func(p *PdfProcessor) {
		p.OptionMetadataPDF = &value
	}
}

// WithOptionFilePDF returns an Option function that sets the OptionFilePDF value.
func WithOptionFilePDF(value entity.OptionFilePDF) Option {
	return func(p *PdfProcessor) {
		v := reflect.ValueOf(value)
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.String() != "" {
				reflect.ValueOf(p.OptionFilePDF).Elem().Field(i).Set(field)
			}
		}
	}
}

// ProcessFile processes the input file based on its type.
func (p *PdfProcessor) ProcessFile() error {
	fileType := util.GetFileType(p.FilePath)
	switch fileType {
	case entity.PDF:
		// Check if the PDF file has a password
		hasPassword, err := util.HasPDFPassword(p.FilePath, p.PasswordPDF)
		if err != nil {
			return err
		}

		p.PDFProtection = hasPassword

		if hasPassword {
			// Decrypt the PDF File
			err := util.Decrypted(p.FilePath, p.PasswordPDF)
			if err != nil {
				return err
			}
		}

		// Process the PDF file
		err = p.processPDF(p.FilePath, p.OptionFilePDF.QRCodePath, p.OptionFilePDF.StampPosition)
		if err != nil {
			return err
		}
	case entity.Image:
		// Convert the image file to PDF
		pdfFilePath, err := ConvertImageToPDF(p.FilePath)
		if err != nil {
			return err
		}

		// Process the converted PDF file
		err = p.processPDF(pdfFilePath, p.OptionFilePDF.QRCodePath, p.OptionFilePDF.StampPosition)
		if err != nil {
			return err
		}
	case entity.Document:
		// Convert the document file to PDF
		pdfFilePath, err := ConvertDocumentToPDF(p.FilePath)
		if err != nil {
			return err
		}

		// Process the converted PDF file
		err = p.processPDF(pdfFilePath, p.OptionFilePDF.QRCodePath, p.OptionFilePDF.StampPosition)
		if err != nil {
			return err
		}

		// Delete the temporary PDF file
		err = os.Remove(pdfFilePath)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported file type")
	}

	return nil
}

// pdfToBase64 converts a PDF file to base64 encoding.
func (p *PdfProcessor) pdfToBase64(filePath string) error {
	// Read the PDF file into memory.
	pdfFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Encode the PDF file as base64.
	base64Output := base64.StdEncoding.EncodeToString(pdfFile)
	// Set the Base64Output field of the PDFProcessor struct.
	p.Base64Output = base64Output

	return nil
}

// processPDF performs operations on the PDF file using pdfcpu-cli.
func (p *PdfProcessor) processPDF(filePath string, qrCode string, stampPosition string) error {
	// Add QR code to the PDF file
	err := AddQRCodeToPDF(filePath, qrCode, stampPosition)
	if err != nil {
		return err
	}

	//add metadata to file pdf
	if !util.IsStructEmpty(p.OptionMetadataPDF) {
		err := util.AddedMetadata(filePath, p.OptionMetadataPDF)
		if err != nil {
			return err
		}

		errs := util.AddKeywords(filePath, p.OptionMetadataPDF)
		if errs != nil {
			return errs
		}
	}

	//add protection to file pdf
	if p.PDFProtection {
		err := util.Encrypted(filePath, p.OptionFilePDF.PasswordPDF)
		if err != nil {
			return err
		}
	}

	//Convert pdf file to base64 as output file
	err = p.pdfToBase64(filePath)
	if err != nil {
		return err
	}

	return nil
}

// AddQRCodeToPDF adds a QR code to the PDF file using pdfcpu-cli.
func AddQRCodeToPDF(filePath string, qrCodeName string, stampPosition string) error {
	if qrCodeName == "" {
		return errors.New("QR Code is empty")
	}

	qrCode := fmt.Sprintf("../samples/images/qr/%v", qrCodeName)

	// Load the icon image
	iconFile, err := os.Open(qrCode)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		} else {
			return fmt.Errorf("error opening file: %s", err.Error())
		}
	}
	defer iconFile.Close()

	fileOutput := fmt.Sprintf("../samples/pdf/out/doc_out.pdf")

	command := fmt.Sprintf("pdfcpu stamp add -pages even,odd  -mode image -- '%s' 'pos:%s, rot:0, sc:.1' %s %s", iconFile.Name(), stampPosition, filePath, fileOutput)

	// Execute the command
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Command exited with a non-zero status
			fmt.Printf("Command failed with error: %s\n", exitError.Error())
			if len(exitError.Stderr) > 0 {
				return fmt.Errorf("error output:%s", string(exitError.Stderr))
			}
		} else {
			// Other execution error
			return fmt.Errorf("error executing pdfcpu command: %s", err.Error())
		}
	}

	return nil
}

// ConvertDocumentToPDF converts a document file to PDF using pdfcpu-cli.
func ConvertDocumentToPDF(doxcName string) (string, error) {
	docxPath := fmt.Sprintf("../samples/images/doc/%s", doxcName)
	pdfPath := fmt.Sprintf("../samples/pdf/out/sample_docx.pdf")

	command := fmt.Sprintf("gs -sDEVICE=pdfwrite -o %s %s", pdfPath, docxPath)

	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to convert DOCX to PDF: %v\n%s", err, output)
	}

	log.Println("Conversion to PDF completed successfully.")

	return pdfPath, nil
}

// ConvertImageToPDF converts an image file to PDF using package gofpdf.
func ConvertImageToPDF(imageFileName string) (string, error) {
	imageFilePath := fmt.Sprintf("../samples/images/doc/%s", imageFileName)
	imageFileOutput := fmt.Sprintf("../samples/pdf/origin/%s", imageFileName)

	// Open the input image file
	outputFile := fmt.Sprintf(util.ChangeFileExtension(imageFileOutput, "pdf"))
	file, err := os.Open(imageFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading image file: %v", err)
	}
	defer file.Close()

	// Read the image file
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set compression to true to reduce the file size.
	pdf.SetCompression(true)

	// Add a new page
	pdf.AddPage()

	// Change producer name
	pdf.SetProducer("privyid", true)

	// Calculate the aspect ratio of the image
	aspectRatio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())

	// Set the image size to fit the page width
	pageWidth, pageHeight := pdf.GetPageSize()
	imageWidth := pageWidth - 20
	imageHeight := (imageWidth / aspectRatio) - 10

	// Calculate the vertical position to center the image
	imageY := (pageHeight - imageHeight) / 2

	// Add the image to the PDF
	pdf.ImageOptions(imageFilePath, 10, imageY, imageWidth, imageHeight, false, gofpdf.ImageOptions{}, 0, "")

	// Save the PDF to the output file
	err = pdf.OutputFileAndClose(outputFile)
	if err != nil {
		return "", fmt.Errorf("error converting image to PDF: %v", err)
	}

	return outputFile, nil
}
