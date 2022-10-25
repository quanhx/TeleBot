package services

import (
	"fmt"
	"github.com/yeqown/go-qrcode"
)

// Define assets accept for payment
const (
	DOGE = "DOGE"
	XRP  = "XRP"
)

const (
	XRPLogo  = "https://cryptologos.cc/logos/thumbs/xrp.png"
	DOGELogo = "https://cryptologos.cc/logos/thumbs/dogecoin.png"
)

var assetsMapping = map[string]string{
	DOGE: DOGELogo,
	XRP:  XRPLogo,
}


func GenQR(text, logo string) (path string) {
	// generating QR Code with source text and output image options.
	var url = fmt.Sprintf("./%s.png", logo)
	fmt.Println(url)
	qrc, err := qrcode.New(text, qrcode.WithLogoImageFilePNG(url)) // Logo: https://cryptologos.cc/logos/thumbs/xrp.png
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
	}

	// save file
	path = fmt.Sprintf("%s.jpeg", logo)
	if err := qrc.Save(path); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
	return path
}
