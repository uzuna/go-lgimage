package lgimage

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"

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
	dc.LoadFontFace("./assets/Lato-Regular.ttf", 32)
	dc.DrawStringAnchored("Anchor Text", float64(W/2), float64(H/2), 0.5, 0.5)
	dc.DrawString("Fill Text", float64(W/2), float64(H/2))

	dc.SavePNG("./demo/demo.png")
}

func TestGGTextsize(t *testing.T) {
	W := 300
	H := 300
	dc := gg.NewContext(W, H)

	fLat16, _ := gg.LoadFontFace("./assets/Lato-Regular.ttf", 16)
	fLat64, _ := gg.LoadFontFace("./assets/Lato-Regular.ttf", 64)

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

	// Text Box
	tbox := TextBox{
		FontFace: fLat16,
		Text:     []string{"Code1", "nCode2", "code3"},
		Color:    color.NRGBA{245, 23, 22, 220},
	}
	tbox.Draw(dc, 0, 0, float64(dc.Width()), float64(dc.Height()))

	dc.SavePNG("./demo/font.png")
}

func TestGGMatrix(t *testing.T) {
	dc := gg.NewContext(300, 300)

	col := uint64(10)
	row := uint64(10)
	width := 300.0
	height := 300.0

	m := Matrix{
		X: 0, Y: 0,
		W: width, H: height,
		Row: row, Col: col,
	}

	dfn := func(dc *gg.Context, x, y, dx, dy float64, ix, iy uint64) {
		r := 10.0 + (15.0 * float64(ix) / float64(col)) + (15.0 * float64(iy) / float64(row))

		// Must use NRGBA. RGNA is an 8-bit alpha-premultipled color
		c := color.NRGBA{uint8(ix * 255 / col), uint8(iy * 255 / row), 150, 180}
		dc.SetColor(c)

		dc.DrawCircle(x, y, r)
		dc.Fill()
	}

	m.Draw(dc, dfn)

	err := dc.SavePNG("./demo/out.png")
	assert.NoError(t, err)
}

func TestGGWithBoxLayout(t *testing.T) {
	dc := gg.NewContext(300, 300)

	l := Layout{}
	// Header
	fLat16, _ := gg.LoadFontFace("./assets/Lato-Regular.ttf", 16)
	tbox := TextBox{
		FontFace: fLat16,
		Text:     []string{"Code1"},
		Color:    color.NRGBA{245, 23, 22, 220},
	}
	l.Header = tbox

	tboxbtm := TextBox{
		FontFace: fLat16,
		Text:     []string{"Footx"},
		Color:    color.NRGBA{12, 90, 200, 220},
		Align:    "right",
	}
	l.Footer = tboxbtm

	// Scale box
	content := func(dc *gg.Context, x, y, w, h float64) {

		var min, ofx float64
		min = w
		if min > h {
			min = h
			ofx = (w - h) / 2 // centering
		}
		col := uint64(10)
		row := uint64(10)

		m := Matrix{
			X: x + ofx, Y: y,
			W: min, H: min,
			Row: row, Col: col,
		}

		dfn := func(dc *gg.Context, x, y, dx, dy float64, ix, iy uint64) {
			r := 10.0 + (5.0 * float64(ix) / float64(col)) + (5.0 * float64(iy) / float64(row))

			// Must use NRGBA. RGNA is an 8-bit alpha-premultipled color
			c := color.NRGBA{uint8(ix * 255 / col), uint8(iy * 255 / row), 150, 180}
			dc.SetColor(c)

			dc.DrawCircle(x+(dx/2), y+(dy/2), r)
			dc.Fill()
		}

		m.Draw(dc, dfn)
	}
	l.Content = ScaleBoxFunc(content)

	// Left Side
	fLat12, _ := gg.LoadFontFace("./assets/Lato-Regular.ttf", 12)
	lsfn := func(dc *gg.Context, x, y, w, h float64) {
		cmap := make(map[string]color.Color)
		cmap["under"] = color.NRGBA{0, 0, 180, 255}
		cmap["over"] = color.NRGBA{180, 0, 0, 255}

		cfn := ValueMapWithFunc{
			Vmin: 0,
			Vmax: 200,
			ColorFunc: func(vi float64) color.Color {
				c := colorful.Hsv(230-vi*230, 0.8, 0.72)
				return c
			},
			ExceptionList: cmap,
		}
		cs := ColorScale{
			X:    x,
			Y:    y,
			W:    w,
			H:    h,
			Vmin: 0,
			Vmax: 200,
			Cfn:  cfn,
			Font: fLat12,
		}

		cs.DrawVertical(dc)
	}
	lside := NewVerticalBoxMargine(lsfn, 4.0, 36)
	l.LSide = lside

	// Draw
	l.Draw(dc)

	err := dc.SavePNG("./demo/layout_r1.png")
	assert.NoError(t, err)

	dc2 := gg.NewContext(300, 300)
	tbox.Text = []string{"Title: Matrix demo", "Desc: Color Matrix", "X: 10, Y: 10"}
	l.Header = tbox
	l.Draw(dc2)

	err = dc2.SavePNG("./demo/layout_r3.png")
	assert.NoError(t, err)
}

func TestGGColorScale(t *testing.T) {
	dc := gg.NewContext(300, 300)
	fLat16, _ := gg.LoadFontFace("./assets/Lato-Regular.ttf", 14)
	cs := ColorScale{
		X:    0,
		Y:    0,
		W:    40,
		H:    300,
		Vmin: 0,
		Vmax: 200,
		Cfn: ColorMapFunc(func(x Value) color.Color {
			v := x.Value().(float64)
			return color.NRGBA{uint8(v / 200 * 255), 128, 128, 255}
		}),
		Font: fLat16,
	}

	cs.DrawVertical(dc)

	{
		cmap := make(map[string]color.Color)
		cmap["under"] = color.NRGBA{0, 20, 255, 255}
		cmap["over"] = color.NRGBA{255, 0, 0, 255}
		cmap["Ng"] = color.NRGBA{255, 255, 0, 255}
		cmap["nil"] = color.NRGBA{0, 255, 255, 255}

		//
		cfn := ValueMapWithFunc{
			Vmin: 0,
			Vmax: 200,
			ColorFunc: func(vi float64) color.Color {
				return color.NRGBA{uint8(vi * 255), 128, 128, 255}
			},
			ExceptionList: cmap,
		}
		cs := ColorScale{
			X:    60,
			Y:    0,
			W:    40,
			H:    300,
			Vmin: 0,
			Vmax: 200,
			Cfn:  cfn,
			Font: fLat16,
		}

		cs.DrawVertical(dc)
	}

	{
		cmap := make(map[string]color.Color)
		cmap["under"] = color.NRGBA{0, 0, 180, 255}
		cmap["over"] = color.NRGBA{180, 0, 0, 255}

		//
		cfn := ValueMapWithFunc{
			Vmin: 0,
			Vmax: 200,
			ColorFunc: func(vi float64) color.Color {
				c := colorful.Hsv(230-vi*230, 0.8, 0.72)
				return c
			},
			ExceptionList: cmap,
		}
		cs := ColorScale{
			X:    120,
			Y:    0,
			W:    40,
			H:    300,
			Vmin: 0,
			Vmax: 200,
			Cfn:  cfn,
			Font: fLat16,
		}

		cs.DrawVertical(dc)
	}
	{
		cmap := make(map[string]color.Color)
		cmap["under"] = color.NRGBA{0, 0, 180, 255}
		cmap["over"] = color.NRGBA{180, 0, 0, 255}

		//
		cfn := ValueMapWithFunc{
			Vmin: 0,
			Vmax: 1,
			ColorFunc: func(vi float64) color.Color {
				ve := EaseOutQuad(vi)
				c := colorful.Hsv(230-ve*230, 0.8, 0.72)
				return c
			},
			ExceptionList: cmap,
		}
		cs := ColorScale{
			X:    180,
			Y:    0,
			W:    40,
			H:    300,
			Vmin: 0,
			Vmax: 1,
			Cfn:  cfn,
			Font: fLat16,
		}

		cs.DrawVertical(dc)
	}
	{
		cmap := make(map[string]color.Color)
		cmap["under"] = color.NRGBA{0, 0, 180, 255}
		cmap["over"] = color.NRGBA{180, 0, 0, 255}

		//
		cfn := ValueMapWithFunc{
			Vmin: 0,
			Vmax: 2105570,
			ColorFunc: func(vi float64) color.Color {
				ve := EaseInQuad(vi)
				c := colorful.Hsv(230-ve*230, 0.8, 0.72)
				return c
			},
			ExceptionList: cmap,
		}
		cs := ColorScale{
			X:    240,
			Y:    0,
			W:    40,
			H:    300,
			Vmin: 0,
			Vmax: 2105570,
			Cfn:  cfn,
			Font: fLat16,
		}

		cs.DrawVertical(dc)
	}

	err := dc.SavePNG("./demo/colorscale.png")
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

	m := Matrix{
		X: 0, Y: 0,
		W: width, H: height,
		Row: row, Col: col,
	}

	dfn := func(dc *gg.Context, x, y, dx, dy float64, ix, iy uint64) {
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
