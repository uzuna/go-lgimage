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
	var lX, tY, rX, bY float64
	if l.Header != nil {
		_, tY = l.Header.Size()
		l.Header.Draw(dc, 0, 0, float64(dc.Width()), tY)
	}
	if l.Footer != nil {
		_, bY = l.Footer.Size()
		y := float64(dc.Height()) - bY
		l.Footer.Draw(dc, 0, y, float64(dc.Width()), bY)
	}

	if l.LSide != nil {
		h := float64(dc.Height()) - (tY + bY)
		l.LSide.Draw(dc, 0, tY, l.LSide.Width(), h)
		lX += l.LSide.Width()
	}
	if l.Content != nil {
		w := float64(dc.Width()) - (lX + rX)
		h := float64(dc.Height()) - (tY + bY)
		l.Content.Draw(dc, lX, tY, w, h)
	}
}

type Scale struct {
	W float64
	H float64
}

type Box interface {
	Draw(dc *gg.Context, x, y, w, h float64)
	Size() (w, h float64)
}

// VerticalBox is has fixed width and auto tune height
type VerticalBox interface {
	Draw(dc *gg.Context, x, y, w, h float64)
	Width() (w float64)
}

type ScaleBox interface {
	Draw(dc *gg.Context, x, y, w, h float64)
}

type ScaleBoxFunc func(dc *gg.Context, x, y, w, h float64)

func (f ScaleBoxFunc) Draw(dc *gg.Context, x, y, w, h float64) {
	f(dc, x, y, w, h)
}

type TextBox struct {
	FontFace font.Face
	Text     []string
	Color    color.Color
	Align    string
}

// Size is return a part of box size.
func (t TextBox) Size() (w, h float64) {
	fontHeight := float64(t.FontFace.Metrics().Height) / 64
	h = fontHeight * float64(len(t.Text))
	return
}

func (t TextBox) Draw(dc *gg.Context, x, y, w, h float64) {
	dc.SetFontFace(t.FontFace)
	dc.SetColor(t.Color)
	fontHeight := float64(t.FontFace.Metrics().Height) / 64
	switch t.Align {
	case "", "left":
		for i, v := range t.Text {
			dc.DrawString(v, x, y+fontHeight*float64(i+1))
		}
	case "right":
		for i, v := range t.Text {
			dc.DrawStringAnchored(v, x+w, y+fontHeight*float64(i+1), 1, 0)
		}
	case "center":
		for i, v := range t.Text {
			dc.DrawStringAnchored(v, x+w/2, y+fontHeight*float64(i+1), 0.5, 0)
		}
	}

}

type VerticalBoxMargine struct {
	drawfn func(dc *gg.Context, x, y, w, h float64)
	margin float64
	width  float64
}

func NewVerticalBoxMargine(drawfn func(dc *gg.Context, x, y, w, h float64), margin, width float64) VerticalBoxMargine {
	return VerticalBoxMargine{
		drawfn: drawfn,
		margin: margin,
		width:  width,
	}
}

func (v VerticalBoxMargine) Draw(dc *gg.Context, x, y, w, h float64) {
	v.drawfn(dc, x+v.margin, y+v.margin, w-v.margin*2, h-v.margin*2)
}

func (v VerticalBoxMargine) Width() float64 {
	return v.width
}
