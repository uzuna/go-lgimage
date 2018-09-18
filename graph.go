package lgimage

import (
	"fmt"
	"image/color"
	"math"
	"sort"

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

	var exeptionKeys []string
	exmap := make(map[string]int)
	for i, v := range list {
		exeptionKeys = append(exeptionKeys, v.String)
		exmap[v.String] = i
	}
	sort.Strings(exeptionKeys)

	// 下からexceptionを埋めていく
	for i, key := range exeptionKeys {
		iv := exmap[key]
		v := list[iv]
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
			} else if v.Value().(float64) < 0.1 {
				dc.DrawStringAnchored(fmt.Sprintf("%.3f", v.Value().(float64)), c.X+c.W-1, y, 1, 0.8)
			} else if v.Value().(float64) < 1 {
				dc.DrawStringAnchored(fmt.Sprintf("%.2f", v.Value().(float64)), c.X+c.W-1, y, 1, 0.8)
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
	Bins       BinsWithException
	Font       font.Face
}

func (c Histgram) DrawVertical(dc *gg.Context) {
	length := len(c.Bins.Bins)
	lengthExp := len(c.Bins.ExceptionBin)
	maxFreq := 0.0

	for _, v := range c.Bins.Bins {
		if v.Frequency > maxFreq {
			maxFreq = v.Frequency
		}
	}
	for _, v := range c.Bins.ExceptionBin {
		if v > maxFreq {
			maxFreq = v
		}
	}
	bWidth := c.W / maxFreq
	zY := c.H + c.Y

	// fontの高さを取得
	dc.SetFontFace(c.Font)
	fontHeight := dc.FontHeight()
	excH := float64(lengthExp) * fontHeight
	numH := c.H - excH

	dy := fontHeight

	// exception
	var exeptionKeys []string
	exmap := make(map[string]float64)
	for k, v := range c.Bins.ExceptionBin {
		exeptionKeys = append(exeptionKeys, k)
		exmap[k] = v
	}
	sort.Strings(exeptionKeys)

	for i, k := range exeptionKeys {
		offY := dy * float64(i)
		dc.DrawRectangle(c.X, zY-offY-dy, bWidth*exmap[k], dy)

		dc.SetColor(color.NRGBA{255, 80, 80, 255})
		dc.StrokePreserve()
		dc.SetColor(color.NRGBA{128, 200, 255, 180})
		dc.Fill()
	}

	//数値系
	zY = c.Y + numH
	bHeight := numH / float64(length)
	for i := 0; i < length; i++ {
		offY := bHeight * float64(i)
		dc.DrawRectangle(c.X, zY-offY-bHeight, bWidth*c.Bins.Bins[i].Frequency, bHeight)

		dc.SetColor(color.NRGBA{255, 255, 255, 255})
		dc.StrokePreserve()
		dc.SetColor(color.NRGBA{128, 200, 255, 180})
		dc.Fill()
	}
}

// HistgramTotalizer 集計器
type HistgramTotalizer struct {
	Vmin, Vmax, vWidth float64
	Step               int
	count              int
	bins               []int
	exbins             map[string]float64
}

// NewHistgramTotalizer is construction HistgramTotalizer
func NewHistgramTotalizer(vmin, vmax float64, step int, exceptions []string) *HistgramTotalizer {
	width := (vmax - vmin) / float64(step-1)
	exbin := make(map[string]float64)
	for _, v := range exceptions {
		exbin[v] = 0
	}
	return &HistgramTotalizer{
		Vmin: vmin, Vmax: vmax, Step: step,
		vWidth: width,
		bins:   make([]int, step),
		exbins: exbin,
	}
}

func (h *HistgramTotalizer) Add(v Value) {
	switch x := v.Value().(type) {
	case float64:
		// over range
		if x >= h.Vmax {
			if _, ok := h.exbins["over"]; !ok {
				h.exbins["over"] = 0
			}
			h.exbins["over"]++
		} else if x < h.Vmin {
			if _, ok := h.exbins["under"]; !ok {
				h.exbins["under"] = 0
			}
			h.exbins["under"]++
		} else {
			x := math.Floor(x / h.vWidth)
			h.bins[int(x)]++
		}
	case string:
		if _, ok := h.exbins[x]; ok {
			h.exbins[x] = 0
		}
		h.exbins[x]++
	}
}

func (h *HistgramTotalizer) ExBins() BinsWithException {
	bins := make([]Bins, h.Step)
	for i, v := range h.bins {
		bins[i] = Bins{
			Vmin:      float64(i) * h.vWidth,
			Vmax:      float64(i+1) * h.vWidth,
			Frequency: float64(v),
		}
	}

	return BinsWithException{
		Bins:         bins,
		ExceptionBin: h.exbins,
	}
}

type BinsWithException struct {
	Bins         []Bins             // histgrams
	ExceptionBin map[string]float64 // 数値以外
}

// ColorScale is provide draw the color scale
type ColorScaleWithHistgram struct {
	X, Y, W, H, Middle float64           // Position
	Vmin, Vmax         float64           // Value Scale
	Cfn                ColorMap          // Color Function
	Font               font.Face         // value font
	ExBins             BinsWithException // histgram情報
}

func (c ColorScaleWithHistgram) DrawVertical(dc *gg.Context) {
	// Set ColorScale
	cs := ColorScale{
		X: c.X, Y: c.Y, W: c.Middle, H: c.H,
		Vmin: c.Vmin, Vmax: c.Vmax,
		Cfn:  c.Cfn,
		Font: c.Font,
	}
	cs.DrawVertical(dc)

	// Set Histgram
	hs := Histgram{
		X: c.X + c.Middle, Y: c.Y, W: c.W - c.Middle, H: c.H,
		Bins: c.ExBins,
		Font: c.Font,
	}
	hs.DrawVertical(dc)
}
