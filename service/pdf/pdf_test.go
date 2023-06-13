package pdf_test

import (
	"fmt"
	entity "pdf-generator-with-qr/entities"
	"pdf-generator-with-qr/service/pdf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertImageToPDF(t *testing.T) {
	doc, err := pdf.ConvertImageToPDF("doc.jpg")

	assert.NoError(t, err)
	assert.NotEmpty(t, doc)
}

func TestConvertDocumentToPDF(t *testing.T) {
	doc, err := pdf.ConvertDocumentToPDF("sample.docx")

	assert.NoError(t, err)
	assert.NotEmpty(t, doc)
}

func TestAddQRCodeToPDF(t *testing.T) {
	err := pdf.AddQRCodeToPDF("../../results/pdf/origin/doc.pdf", "qrcode-with-icon.png", "br")

	assert.NoError(t, err)
}

func TestProcessPDF(t *testing.T) {
	pdfProcess := pdf.NewPDFGopher("../../results/pdf/out/doc_out.pdf",
		pdf.WithOptionMetadataPDF(entity.OptionMetadataPDF{Title: "Hero life in You", Author: "Me as Author", Subject: "You as Subject", Keywords: "Kopi Luwak"}),
		pdf.WithOptionFilePDF(entity.OptionFilePDF{QRCodePath: "qrcode-with-icon.png", StampPosition: "tl"}),
	)

	// Process file
	err := pdfProcess.ProcessFile()

	fmt.Printf("pdfProcess.Base64Output: %v\n", pdfProcess.Base64Output)

	assert.NoError(t, err)
	assert.NotEmpty(t, pdfProcess.Base64Output)
}
