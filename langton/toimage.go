package langton

import (
	"image"
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

func ToPalette(palette []colorful.Color) color.Palette {
	colorPalette := make(color.Palette, len(palette)+1)
	colorPalette[0] = color.Alpha{}
	for i := range palette {

		r, g, b := palette[i].RGB255()

		_, _, _, a := palette[i].RGBA()
		colorPalette[i+1] = color.RGBA{
			R: r,
			G: g,
			B: b,
			A: uint8(math.Sqrt(float64(a))),
		}
	}
	return colorPalette
}

func ToImage(ant *Ant, palette color.Palette, cellSize int) *image.Paletted {

	r := image.Rect(
		0,
		0,
		int(ant.Dimensions.width)*cellSize,
		int(ant.Dimensions.height)*cellSize,
	)
	img := image.NewPaletted(r, palette)
	for i := range ant.Cells {
		if ant.Cells[i].Step.Action == ActionNone {
			continue
		}

		for sx := 0; sx < cellSize; sx++ {
			for sy := 0; sy < cellSize; sy++ {
				img.SetColorIndex(
					int((ant.Cells[i].X+ant.Dimensions.width/2)*int64(cellSize)+int64(sx)),
					int((ant.Cells[i].Y+ant.Dimensions.height/2)*int64(cellSize)+int64(sy)),
					uint8(ant.Cells[i].Step.Index+1),
				)
			}
		}
	}
	return img
}
