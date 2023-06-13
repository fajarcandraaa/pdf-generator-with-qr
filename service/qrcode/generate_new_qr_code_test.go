package qrcode_test

import (
	"fmt"
	"pdf-generator-with-qr/service/qrcode"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQRCode(t *testing.T) {
	key := "pdf-generator-with-qr"
	text := fmt.Sprintf("https://github.com/fajarcandraaa/%s", key)

	err := qrcode.QrCodeGen(text, "qrcode-without-icon.png")

	assert.NoError(t, err)
}
