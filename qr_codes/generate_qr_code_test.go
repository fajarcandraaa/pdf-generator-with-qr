package qr_codes_test

import (
	"fmt"
	"pdf-generator-with-qr/qr_codes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQRCode(t *testing.T) {
	key := "simple-pdf-generator-with-qr"
	text := fmt.Sprintf("https://github.com/akhidnukhlis/simple-pdf-generator-with-qr/%s", key)

	err := qr_codes.QrCodeGen(text, "qrcode-without-icon.png")

	assert.NoError(t, err)
}
