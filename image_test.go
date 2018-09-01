package main

import (
	"testing"

	"github.com/fogleman/gg"
)

func TestGG(t *testing.T) {
	dc := gg.NewContext(1000, 1000)
	dc.DrawCircle(500, 500, 400)
	dc.SetRGBA(245, 23, 22, 0.8)
	dc.Fill()
	dc.SavePNG("out.png")
}
