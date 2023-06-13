package qrcode_test

import (
	"fmt"
	"pdf-generator-with-qr/service/qrcode"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQRCodeWithIcon(t *testing.T) {
	key := "pdf-generator-with-qr"
	text := fmt.Sprintf("https://github.com/fajarcandraaa/%s", key)

	qrCode, err := qrcode.GenerateQRCodeWithIcon(text, "privy_logo.png", "qrcode-with-icon.png")

	assert.NoError(t, err)
	assert.NotEmpty(t, qrCode)
}
