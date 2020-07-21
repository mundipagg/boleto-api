package boleto

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"
	
	s "strings"
	"github.com/golang/freetype"
	"github.com/mundipagg/boleto-api/util"
	"golang.org/x/image/font"
)

func textToImage(text string) string {
	
	if s.Contains(text,"  ") {		
		text = s.Replace(text," ","",-1)
		text = s.Replace(text,".","",-1)	
		text = formatDigitableLine(text)
	}
	
	size := float64(13)
	dpi := float64(100)
	rgba := image.NewNRGBA64(image.Rect(0, 0, 530, 20))
	draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(util.GetFont().FtFont)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.Black)
	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(10, 8+int(c.PointToFixed(size)>>7))
	for _, s := range []string{text} {
		c.DrawString(s, pt)
		pt.Y += c.PointToFixed(size)
	}

	data := bytes.NewBuffer(nil)
	png.Encode(data, rgba)
	return base64.StdEncoding.EncodeToString(data.Bytes())
}

func formatDigitableLine(s string) string {
	buf := bytes.Buffer{}
	for idx, c := range s {
		if idx == 5 || idx == 15 || idx == 26 {
			buf.WriteString(".")
		}
		if idx == 10 || idx == 21 || idx == 32 || idx == 33 {
			buf.WriteString(" ")
		}
		buf.WriteByte(byte(c))
	}
	return buf.String()
}
