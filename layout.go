package lgimage

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

type Layout struct {
	Header  Box
	Content Box
	Footer  Box
}

type Scale struct {
	W float64
	H float64
}

type Box interface {
	Draw(dc *gg.Context, x, y float64)
	Size() (w, h float64)
}

type TextBox struct {
	FontFace font.Face
	Text     []string
}

func (t TextBox) Size() (w, h float64) {
	fontHeight := float64(t.FontFace.Metrics().Height) / 64
	h = fontHeight * float64(len(t.Text))
	return
}

func (t TextBox) Draw(dc *gg.Context, x, y float64) {
	dc.SetFontFace(t.FontFace)
	fontHeight := float64(t.FontFace.Metrics().Height) / 64
	for i, v := range t.Text {
		dc.DrawString(v, x, y+fontHeight*float64(i+1))
	}
}
