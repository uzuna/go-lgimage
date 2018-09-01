package main

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uzuna/lgimage"

	"github.com/fogleman/gg"
)

func TestGG(t *testing.T) {
	W := 300
	H := 300
	var mask *image.Alpha

	dc := gg.NewContext(W, H)
	// mask

	dc.DrawRectangle(10, 80, 280, 140)
	dc.SetRGBA(0, 0, 0, 1)
	dc.Fill()
	mask = dc.AsMask()
	dc.Clear()

	dc.SetMask(mask)

	// image render
	dc.DrawCircle(float64(W/2), float64(H/2), 120)
	dc.SetRGBA(245, 23, 22, 0.8)
	dc.Fill()
	dc.LoadFontFace("../assets/Lato-Regular.ttf", 32)
	dc.DrawStringAnchored("Anchor Text", float64(W/2), float64(H/2), 0.5, 0.5)
	dc.DrawString("Fill Text", float64(W/2), float64(H/2))

	dc.SavePNG("demo.png")
}

func TestGGTextsize(t *testing.T) {
	W := 300
	H := 300
	dc := gg.NewContext(W, H)

	fLat16, _ := gg.LoadFontFace("../assets/Lato-Regular.ttf", 16)
	fLat64, _ := gg.LoadFontFace("../assets/Lato-Regular.ttf", 64)

	dc.DrawRectangle(10, 80, 280, 140)
	dc.SetRGBA(245, 23, 22, 0.8)
	dc.SetFontFace(fLat16)
	var offset float64
	offset = dc.FontHeight()
	dc.DrawString("16 lag\r\n 4tghu", 0, float64(H/2)-offset)
	dc.DrawString("16 Telagxt2", 0, float64(H/2))

	dc.SetFontFace(fLat64)
	offset = dc.FontHeight()
	dc.DrawString("64 lagg", 60, float64(H/2))
	dc.DrawString("64 lag2", 60, float64(H/2)+offset)

	tbox := lgimage.TextBox{
		FontFace: fLat16,
		Text:     []string{"Code1", "nCode2", "code3"},
	}

	tbox.Draw(dc, 0, 0)

	dc.SavePNG("font.png")
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
