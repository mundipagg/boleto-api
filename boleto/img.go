package boleto

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"

	"github.com/pbnjay/pixfont"
)

func textToImage(text string) string {
	img := image.NewRGBA(image.Rect(0, 0, 530, 20))
	pixfont.DrawString(img, 5, 5, text, color.Black)
	f := bytes.NewBuffer(nil)
	png.Encode(f, img)
	return base64.StdEncoding.EncodeToString(f.Bytes())
}
