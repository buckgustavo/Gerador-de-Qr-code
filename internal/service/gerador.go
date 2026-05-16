package service

import (
	"encoding/base64"

	"gerador-qrcode/internal/domain"

	qrcode "github.com/skip2/go-qrcode"
)

type GeradorQR struct{}

func NovoGeradorQR() *GeradorQR {
	return &GeradorQR{}
}

func (g *GeradorQR) GerarBase64(c *domain.Cracha) (string, error) {
	png, err := qrcode.Encode(c.ConteudoQR(), qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}
