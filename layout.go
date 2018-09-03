package lgimage

import (
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

type Layout struct {
	Header  Box
	Content ScaleBox // AutoResize
	RSide   VerticalBox
	LSide   VerticalBox
	Footer  Box
}

func (l Layout) Draw(dc *gg.Context) {
	var cX, cY float64
	if l.Header != nil {
		l.Header.Draw(dc, 0, 0)
		_, cY = l.Header.Size()
	}

	if l.LSide != nil {
		h := float64(dc.Height()) - cY
		l.LSide.Draw(dc, l.LSide.Width(), h, 0, cY)
		cX += l.LSide.Width()
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

// VerticalBox is has fixed width and auto tune height
type VerticalBox interface {
	Draw(dc *gg.Context, w, h, x, y float64)
	Width() (w float64)
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

type VerticalBoxMargine struct {
	drawfn func(dc *gg.Context, w, h, x, y float64)
	margin float64
	width  float64
}

func NewVerticalBoxMargine(drawfn func(dc *gg.Context, w, h, x, y float64), margin, width float64) VerticalBoxMargine {
	return VerticalBoxMargine{
		drawfn: drawfn,
		margin: margin,
		width:  width,
	}
}

func (v VerticalBoxMargine) Draw(dc *gg.Context, w, h, x, y float64) {
	v.drawfn(dc, w-v.margin*2, h-v.margin*2, x+v.margin, y+v.margin)
}

func (v VerticalBoxMargine) Width() float64 {
	return v.width
}
