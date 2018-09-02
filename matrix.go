package lgimage

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// Matrix is provide draw the matrix
type Matrix struct {
	X, Y, W, H float64
	Col, Row   uint64
}

type DrawFunc func(dc *gg.Context, x, y, dx, dy float64, ix, iy uint64)

// Draw is tell the point on matrix.
func (m Matrix) Draw(dc *gg.Context, dfn DrawFunc) {
	dx := m.W / float64(m.Col)
	dy := m.H / float64(m.Row)
	for iy := uint64(0); iy < m.Row; iy++ {
		for ix := uint64(0); ix < m.Col; ix++ {
			dfn(dc, m.X+dx*float64(ix), m.Y+dy*float64(iy), dx, dy, ix, iy)
		}
	}
}

// ColorScale is provide draw the color scale
type ColorScale struct {
	X, Y, W, H float64                     // Position
	Vmin, Vmax float64                     // Value Scale
	Cfn        func(v float64) color.Color // Color Function
	Font       font.Face                   // value font
}

// Draw
func (c ColorScale) DrawVertical(dc *gg.Context) {
	s := 20
	dy := c.H / float64(s)
	zY := c.Y + c.H
	dv := (c.Vmax - c.Vmin) / float64(s)

	// 表示は上下と
	dc.SetFontFace(c.Font)
	// fontHeight := dc.FontHeight()

	for i := 0; i < s; i++ {
		v := c.Vmin + (dv * float64(i))
		cc := c.Cfn(v)
		dc.SetColor(cc)
		y := zY - dy*float64(i+1)
		dc.DrawRectangle(c.X, y, c.W, dy)
		// dc.DrawRectangle(20*float64(i), 20, 40, 40)
		dc.Fill()

		// Text
		if i == 0 || i == s-1 || i%3 == 0 {
			dc.SetColor(color.NRGBA{255, 255, 255, 255})
			dc.DrawStringAnchored(fmt.Sprintf("%.1f ", v), c.X+c.W, y, 1, 1)
			dc.Fill()
		}
	}
}
