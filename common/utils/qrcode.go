package utils

import (
	"github.com/skip2/go-qrcode"
)

func GetQRCodeIO(content string, size ...int) ([]byte, error) {
	s := 256
	if len(size) > 0 {
		s = size[0]
	}
	png, err := qrcode.Encode(content, qrcode.High, s)
	if err != nil {
		return nil, err
	}
	return png, nil
}
