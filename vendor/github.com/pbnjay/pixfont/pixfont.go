// Package pixfont is a simple pixel font library that allows for text drawing in the
// standard image and image/draw packages. It is very useful when heavyweight TrueType
// fonts and the freetype-go package are not available or desired. Since fonts are
// simple bitmap images, scaling is not possible and aliasing will occur.
//
// See the included fontgen tool if you wish to convert or include your own pixel font
// in your project.
package pixfont

import (
	"bytes"
	"image/color"
)

// DefaultFont is used by the convienence method DrawString, and is initialized
// to the Public Domain 8x8 fixed font with some unicode characters.
var DefaultFont = Font8x8

// Spacing is the pixel spacing to use between letters (1 px by default)
var Spacing = 1

// Drawable is an interface which supports setting an x,y coordinate to a color.
type Drawable interface {
	Set(x, y int, c color.Color)
}

// PixFont represents a simple bitmap or pixel-based font that can be drawn using
// simple opaque-pixel operations (supported by image.Image and easily included
// in other packages).
type PixFont struct {
	charWidth    uint8
	charHeight   uint8
	charmap      map[rune]uint16
	data         []uint32
	varCharWidth uint8
}

// NewPixFont creates a new PixFont with the provided character width/height and
// character map of offsets into a packed uint32 array of bits.
func NewPixFont(w, h uint8, cm map[rune]uint16, d []uint32) *PixFont {
	return &PixFont{w, h, cm, d, w}
}

// SetVariableWidth toggles the PixFont between drawing using variable width
// per character or the default fixed-width representation.
func (p *PixFont) SetVariableWidth(isVar bool) {
	if !isVar {
		p.varCharWidth = p.charWidth
	} else {
		// spaces will be approx 1/3 em (but at least 3px)
		p.varCharWidth = p.charWidth / 3
		if p.varCharWidth < 3 {
			p.varCharWidth = 3
		}
	}
}

// DrawRune uses this PixFont to display a single rune in the provided color and
// position in Drawable. The x,y position represents the top-left corner of the rune.
// Drawable.Set is called for each opaque pixel in the font, leaving all other pixels
// in the Drawable as-is. If the rune has no representation in the PixFont, then
// DrawRune returns false and no drawing is done. DrawRune always returns the number
// of pixels to advance before drawing another character.
func (p *PixFont) DrawRune(dr Drawable, x, y int, c rune, clr color.Color) (bool, int) {
	poff, haveChar := p.charmap[c]
	if !haveChar {
		return false, int(p.varCharWidth)
	}
	w := int(p.charWidth)
	if p.varCharWidth != p.charWidth {
		w = 0
	}
	pindex := int(poff >> 2)
	psub := (poff & 0x03) * 8
	d := p.data[pindex : pindex+int(p.charHeight)]
	for yy := 0; yy < int(p.charHeight); yy++ {
		bitMask := uint32(1) << psub
		for xx := 0; xx < int(p.charWidth); xx++ {
			if (d[yy] & bitMask) != 0 {
				dr.Set(x+xx, y+yy, clr)
				if xx >= w {
					w = xx + Spacing
				}
			}
			bitMask <<= 1
		}
	}
	return true, w
}

// DrawString uses this PixFont to display text in the provided color and the specified
// start position in Drawable. The x,y position represents the top-left corner of the
// first letter of s. Text is drawn by repeated calls to DrawRune for each character.
// DrawString returns the total pixel advance used by the string.
func (p *PixFont) DrawString(dr Drawable, x, y int, s string, clr color.Color) int {
	for _, c := range s {
		_, w := p.DrawRune(dr, x, y, c, clr)
		x += w + Spacing
	}
	return x
}

// MeasureRune measures the advance of a rune drawn using this PixFont.
func (p *PixFont) MeasureRune(c rune) (bool, int) {
	poff, haveChar := p.charmap[c]
	if p.varCharWidth == p.charWidth {
		return haveChar, int(p.charWidth)
	}
	if !haveChar {
		return haveChar, int(p.varCharWidth)
	}
	w := 0
	pindex := int(poff >> 2)
	psub := (poff & 0x03) * 8
	d := p.data[pindex : pindex+int(p.charHeight)]
	for yy := 0; yy < int(p.charHeight); yy++ {
		bitMask := uint32(1) << psub
		for xx := w; xx < int(p.charWidth); xx++ {
			if (d[yy]&bitMask) != 0 && xx >= w {
				w = xx + Spacing
			}
			bitMask <<= 1
		}
	}
	return true, w
}

// MeasureString measures the pixel advance of a string drawn using this PixFont.
func (p *PixFont) MeasureString(s string) int {
	x := 0
	for _, c := range s {
		_, w := p.MeasureRune(c)
		x += w + Spacing
	}
	return x
}

// DrawString is a convienence method that calls DrawString using the DefaultFont
func DrawString(dr Drawable, x, y int, s string, clr color.Color) int {
	return DefaultFont.DrawString(dr, x, y, s, clr)
}

// MeasureString is a convienence method that calls MeasureString using the DefaultFont
func MeasureString(s string) int {
	return DefaultFont.MeasureString(s)
}

///////

// StringDrawable implements Drawable so you can do FIGlet-inspired pixel fonts in
// text. Obviously it's much simpler though.
type StringDrawable struct {
	lines [][]byte
}

func (s *StringDrawable) Set(x, y int, c color.Color) {
	for len(s.lines) <= y {
		s.lines = append(s.lines, make([]byte, x))
	}

	if len(s.lines[y]) <= x {
		nb := make([]byte, 1+(x-len(s.lines[y])))
		s.lines[y] = append(s.lines[y], nb...)
	}

	s.lines[y][x] = byte('X')
}

// String returns the current string representation of this Drawable.
func (s *StringDrawable) String() string {
	return s.PrefixString("")
}

// PrefixString returns the current string representation of this Drawable with a
// user-provided prefix before each line. Useful for adding output in code comments.
func (s *StringDrawable) PrefixString(p string) string {
	r := ""
	for _, line := range s.lines {
		r += p + string(bytes.Replace(line, []byte{0}, []byte(" "), -1)) + "\n"
	}
	return r
}
