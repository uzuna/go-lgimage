package main

import (
	"bytes"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uzuna/lgimage"

	"github.com/fogleman/gg"
)

func TestGG(t *testing.T) {
	dc := gg.NewContext(300, 300)
	dc.DrawCircle(150, 150, 120)
	dc.SetRGBA(245, 23, 22, 0.8)
	dc.Fill()
	dc.SavePNG("out.png")
}

func TestGGMatrix(t *testing.T) {
	dc := gg.NewContext(300, 300)

	col := uint64(10)
	row := uint64(10)
	width := 300.0
	height := 300.0

	m := lgimage.Matrix{
		X: 0, Y: 0,
		W: width, H: height,
		Row: row, Col: col,
	}

	dfn := func(dc *gg.Context, x, y float64, ix, iy uint64) {
		r := 10.0 + (15.0 * float64(ix) / float64(col)) + (15.0 * float64(iy) / float64(row))

		// Must use NRGBA. RGNA is an 8-bit alpha-premultipled color
		c := color.NRGBA{uint8(ix * 255 / col), uint8(iy * 255 / row), 150, 180}
		dc.SetColor(c)

		dc.DrawCircle(x, y, r)
		dc.Fill()
	}

	m.Draw(dc, dfn)

	err := dc.SavePNG("out.png")
	assert.NoError(t, err)
}

func BenchmarkRenderInsntance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dc := gg.NewContext(300, 300)
		dc.DrawCircle(250, 250, 120)
		dc.SetRGBA(245, 23, 22, 0.8)
		dc.Fill()
	}
}

func BenchmarkRenderMatrix100(b *testing.B) {
	col := uint64(10)
	row := uint64(10)
	width := 300.0
	height := 300.0

	m := lgimage.Matrix{
		X: 0, Y: 0,
		W: width, H: height,
		Row: row, Col: col,
	}

	dfn := func(dc *gg.Context, x, y float64, ix, iy uint64) {
		r := 10.0 + (15.0 * float64(ix) / float64(col)) + (15.0 * float64(iy) / float64(row))

		// Must use NRGBA. RGNA is an 8-bit alpha-premultipled color
		c := color.NRGBA{uint8(ix * 255 / col), uint8(iy * 255 / row), 150, 180}
		dc.SetColor(c)

		dc.DrawCircle(x, y, r)
		dc.Fill()
	}

	for i := 0; i < b.N; i++ {
		dc := gg.NewContext(300, 300)
		m.Draw(dc, dfn)
	}
}

func BenchmarkRenderToBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		dc := gg.NewContext(300, 300)
		dc.DrawCircle(250, 250, 120)
		dc.SetRGBA(245, 23, 22, 0.8)
		dc.Fill()
		dc.EncodePNG(buf)
	}
}
