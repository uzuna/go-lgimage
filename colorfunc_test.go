package lgimage

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorFunc(t *testing.T) {
	data := []struct {
		v      interface{}
		expect color.Color
	}{
		{0.0, color.NRGBA{0, 128, 128, 128}},
		{1200.0, color.NRGBA{32, 32, 32, 255}},
		{"N/A", color.NRGBA{255, 255, 0, 255}},
		{nil, color.NRGBA{0, 255, 255, 255}},
		{"nodata", color.NRGBA{0, 0, 0, 0}},
		{1, color.NRGBA{1, 128, 128, 128}},
		{-3, color.NRGBA{255, 255, 255, 255}},
	}

	var vs []Value

	for _, v := range data {
		c, _ := NewValue(v.v)
		vs = append(vs, c)
	}

	cmap := make(map[string]color.Color)

	cmap["over"] = color.NRGBA{32, 32, 32, 255}
	cmap["N/A"] = color.NRGBA{255, 255, 0, 255}
	cmap["nil"] = color.NRGBA{0, 255, 255, 255}

	cmin := color.NRGBA{0, 128, 128, 128}
	cmax := color.NRGBA{255, 128, 128, 128}
	dr, dg, db, da := cmax.R-cmin.R, cmax.G-cmin.G, cmax.B-cmin.B, cmax.A-cmin.A
	cfn := func(v float64) color.Color {
		return color.NRGBA{
			cmin.R + uint8(float64(dr)*v),
			cmin.G + uint8(float64(dg)*v),
			cmin.B + uint8(float64(db)*v),
			cmin.A + uint8(float64(da)*v),
		}
	}

	cmfn := ValueMapWithFunc{
		Vmin:          0,
		Vmax:          255,
		ColorFunc:     cfn,
		ExceptionList: cmap,
	}

	for i, v := range vs {
		c := cmfn.Color(v)
		// log.Println(c, data[i].expect)
		assert.Equal(t, data[i].expect, c)
	}
}
