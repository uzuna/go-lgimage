package lgimage

import (
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

type Layout struct {
	Header  Box
	Content ScaleBox // AutoResize
	Footer  Box
}

func (l Layout) Draw(dc *gg.Context) {
	var cX, cY float64
	if l.Header != nil {
		l.Header.Draw(dc, 0, 0)
		_, cY = l.Header.Size()
	}
	if l.Content != nil {
		w := float64(dc.Width()) - cX
		h := float64(dc.Height()) - cY
		l.Content.Draw(dc, w, h, cX, cY)
	}
}

type Scale struct {
	W float64
	H float64
}

type Box interface {
	Draw(dc *gg.Context, x, y float64)
	Size() (w, h float64)
}

type ScaleBox interface {
	Draw(dc *gg.Context, w, h, x, y float64)
}

type ScaleBoxFunc func(dc *gg.Context, w, h, x, y float64)

func (f ScaleBoxFunc) Draw(dc *gg.Context, w, h, x, y float64) {
	f(dc, w, h, x, y)
}

type TextBox struct {
	FontFace font.Face
	Text     []string
	Color    color.Color
}

// Size is return a part of box size.
func (t TextBox) Size() (w, h float64) {
	fontHeight := float64(t.FontFace.Metrics().Height) / 64
	h = fontHeight * float64(len(t.Text))
	return
}

func (t TextBox) Draw(dc *gg.Context, x, y float64) {
	dc.SetFontFace(t.FontFace)
	dc.SetColor(t.Color)
	fontHeight := float64(t.FontFace.Metrics().Height) / 64
	for i, v := range t.Text {
		dc.DrawString(v, x, y+fontHeight*float64(i+1))
	}
}
