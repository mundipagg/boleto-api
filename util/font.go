package util

import (
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"github.com/mundipagg/boleto-api/log"

	"github.com/golang/freetype"
)

type font struct {
	FtFont *truetype.Font
}

var fnt font

func GetFont() font {

	if (font{}) == fnt {
		fontBytes, err := ioutil.ReadFile("./boleto/Arial.ttf")
		if err != nil {
			l := log.CreateLog()
			l.Fatal(err.Error(), " An error has occurred load font")
		}

		f, err := freetype.ParseFont(fontBytes)
		if err != nil {
			l := log.CreateLog()
			l.Fatal(err.Error(), " An error has occurred load font")
		}

		fnt = font{
			FtFont: f,
		}
	}

	return fnt
}