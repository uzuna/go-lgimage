package lgimage

import (
	"image/color"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/fogleman/gg"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
	"github.com/uzuna/lgimage/ease"
)

func TestHistgram(t *testing.T) {

	os.Mkdir("./demo", 0555)
	// 再現データ作成

	hst := NewHistgramTotalizer(0, 200, 20, []string{"over", "under"})
	for i := 0; i < 400; i++ {
		v, _ := NewValue(math.Pow(rand.Float64(), 2) * 300)
		hst.Add(v)
	}
	fLat16, _ := gg.LoadFontFace("./assets/Lato-Regular.ttf", 16)
	// log.Println(bins)
	// bin詰
	hist := Histgram{
		40, 0, 30, 300, hst.ExBins(),
		fLat16,
	}
	// draw test
	dc := gg.NewContext(300, 300)

	hist.DrawVertical(dc)

	err := dc.SavePNG("demo/hist.png")
	assert.NoError(t, err)
}

func TestCSHistgram(t *testing.T) {

	os.Mkdir("./demo", 0555)
	cmap := make(map[string]color.Color)
	cmap["under"] = color.NRGBA{0, 0, 180, 255}
	cmap["over"] = color.NRGBA{180, 0, 0, 255}

	// ColorFunction
	fLat16, _ := gg.LoadFontFace("./assets/Lato-Regular.ttf", 16)
	cfn := ValueMapWithFunc{
		Vmin: 0,
		Vmax: 200,
		ColorFunc: func(vi float64) color.Color {
			ve := ease.EaseInQuad(vi)
			c := colorful.Hsv(230-ve*230, 0.8, 0.72)
			return c
		},
		ExceptionList: cmap,
	}

	// histgram
	hst := NewHistgramTotalizer(0, 200, 20, []string{"over", "under"})
	for i := 0; i < 400; i++ {
		v, _ := NewValue(math.Pow(rand.Float64(), 2)*250 - 40)
		hst.Add(v)
	}
	chs := ColorScaleWithHistgram{
		X: 0, Y: 0, Middle: 40, W: 80, H: 300,
		Vmin: 0, Vmax: 200,
		Cfn:    cfn,
		Font:   fLat16,
		ExBins: hst.ExBins(),
	}

	// draw test
	dc := gg.NewContext(300, 300)

	chs.DrawVertical(dc)

	err := dc.SavePNG("demo/colorhist.png")
	assert.NoError(t, err)
}

// 中間点を最高にとる
func createQuadbins(vmin, anchor, vmax, max float64, length int) []Bins {
	// generate v points
	l := float64(length)
	width := (vmax - vmin) / l
	anchorBin := 0
	for i := 0; i < length; i++ {
		iv := (width * float64(i))
		if anchor < iv {
			if anchor-iv >= (width / 2) {
				anchorBin = i
			} else {
				anchorBin = i
			}
			break
		}
	}

	// binsデータ生成
	var bins []Bins
	for i := 0; i < anchorBin; i++ {
		bmin := (width * float64(i))
		bmax := (width * float64(i+1))

		x := (bmin + (width / 2)) / anchor
		x = ease.EaseInQuad(x)
		bins = append(bins, Bins{bmin, bmax, x * max})
	}
	for i := anchorBin; i < length; i++ {
		bmin := (width * float64(i))
		bmax := (width * float64(i+1))

		x := 1 - ((bmin - anchor + (width / 2)) / (vmax - anchor))
		x = ease.EaseInQuad(x)
		bins = append(bins, Bins{bmin, bmax, x * max})
	}

	return bins
}
