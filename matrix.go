package lgimage

import (
	"fmt"
	"image/color"
	"math"

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
	X, Y, W, H float64   // Position
	Vmin, Vmax float64   // Value Scale
	Cfn        ColorMap  // Color Function
	Font       font.Face // value font
}

// Draw
func (c ColorScale) DrawVertical(dc *gg.Context) {

	// fontの高さを取得
	dc.SetFontFace(c.Font)
	fontHeight := dc.FontHeight()

	// exception color list
	// reserve rendering area by font size
	list := c.Cfn.List()

	excH := float64(len(list)) * fontHeight
	numH := c.H - excH

	dy := fontHeight
	zY := c.Y + c.H // 書き込み富順高さを指定

	// 下からexceptionを埋めていく
	for i, v := range list {
		y := zY - dy*float64(i+1)

		// Color Rect
		dc.SetColor(v.Color)
		dc.DrawRectangle(c.X, y, c.W, dy)
		dc.Fill()

		// Change text color
		r, g, b, _ := v.Color.RGBA()
		brightness := float64(((r*299)+(g*587)+(b*114))/1000) / 256.0

		// Text
		if brightness < 128 {
			dc.SetColor(color.NRGBA{255, 255, 255, 255})
		} else {
			dc.SetColor(color.NRGBA{0, 0, 0, 255})
		}
		dc.DrawStringAnchored(v.String, c.X+c.W-1, y, 1, 0.8)
		dc.Fill()
	}

	// Paint scale
	splitLen := 20
	zY = c.Y + numH // 書き込み基準高さを指定
	dy = numH / float64(splitLen)
	dv := (c.Vmax - c.Vmin) / float64(splitLen)
	step := int(math.Ceil(fontHeight * 3 / dv))

	for i := 0; i < splitLen; i++ {
		v, _ := NewValue(c.Vmin + (dv * float64(i)))
		cc := c.Cfn.Color(v)
		y := math.Ceil(zY - dy*float64(i+1))

		// Color Rect
		dc.SetColor(cc)
		dc.DrawRectangle(c.X, y, c.W, math.Ceil(dy))
		dc.Fill()

		// Text
		if i == 0 || i == splitLen-1 || i%step == 0 {
			dc.SetColor(color.NRGBA{255, 255, 255, 255})
			dc.DrawStringAnchored(fmt.Sprintf("%.1f ", v.Value().(float64)), c.X+c.W-1, y, 1, 0.8)
			dc.Fill()
		}
	}

}
