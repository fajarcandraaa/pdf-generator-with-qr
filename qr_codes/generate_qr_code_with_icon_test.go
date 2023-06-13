package qr_codes_test

import (
	"fmt"
	"pdf-generator-with-qr/qr_codes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQRCodeWithIcon(t *testing.T) {
	key := "simple-pdf-generator-with-qr"
	text := fmt.Sprintf("https://github.com/akhidnukhlis/simple-pdf-generator-with-qr/%s", key)

	qrCode, err := qr_codes.GenerateQRCodeWithIcon(text, "favicon.png", "qrcode-with-icon.png")

	assert.NoError(t, err)
	assert.NotEmpty(t, qrCode)
}
