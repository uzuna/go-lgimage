package lgimage

import (
	"os"
	"testing"

	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
	"github.com/uzuna/lgimage/ease"
)

func TestHistgram(t *testing.T) {

	os.Mkdir("./demo", 0666)
	// 再現データ作成
	bins := createQuadbins(0, 40, 200, 40, 40)
	// log.Println(bins)
	// bin詰
	hist := Histgram{
		0, 0, 30, 300, bins,
	}
	// draw test
	dc := gg.NewContext(300, 300)

	hist.DrawVertical(dc)

	err := dc.SavePNG("demo/hist.png")
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
