package lgimage

import (
	"github.com/fogleman/gg"
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
