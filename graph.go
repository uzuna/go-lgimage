package lgimage

import (
	"fmt"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

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
	dv := (c.Vmax - c.Vmin) / float64(splitLen-1)
	step := 4 // 固定値5分割。できればheightの長さに合わせたい。

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
			if v.Value().(float64) > 1000 {
				dc.DrawStringAnchored(fmt.Sprintf("%.2g", v.Value().(float64)), c.X+c.W-1, y, 1, 0.8)
			} else {
				dc.DrawStringAnchored(fmt.Sprintf("%.1f", v.Value().(float64)), c.X+c.W-1, y, 1, 0.8)
			}
			dc.Fill()
		}
	}
}

type Bins struct {
	Vmin, Vmax float64
	Frequency  float64
}

// Histgram is render of histgram
type Histgram struct {
	X, Y, W, H float64 // Position
	Bins       []Bins
}

func (c Histgram) DrawVertical(dc *gg.Context) {
	length := len(c.Bins)
	maxFreq := 0.0
	bHeight := c.H / float64(length)
	for _, v := range c.Bins {
		if v.Frequency > maxFreq {
			maxFreq = v.Frequency
		}
	}
	bWidth := c.W / maxFreq
	zY := c.H + c.Y

	for i := 0; i < length; i++ {
		offY := bHeight * float64(i)
		dc.DrawRectangle(c.X, zY-offY-bHeight, bWidth*c.Bins[i].Frequency, bHeight)

		dc.SetColor(color.NRGBA{255, 255, 255, 255})
		dc.StrokePreserve()
		dc.SetColor(color.NRGBA{128, 200, 255, 180})
		dc.Fill()
	}
}
