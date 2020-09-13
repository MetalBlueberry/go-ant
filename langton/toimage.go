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

func ToImage(ant *Ant, palette color.Palette) *image.Paletted {

	r := image.Rect(
		0,
		0,
		int(ant.Dimensions.width),
		int(ant.Dimensions.height),
	)
	img := image.NewPaletted(r, palette)
	for i := range ant.Cells {
		if ant.Cells[i].Step.Action == ActionNone {
			continue
		}

		img.SetColorIndex(
			int(ant.Cells[i].X+ant.Dimensions.width/2),
			int(ant.Cells[i].Y+ant.Dimensions.height/2),
			uint8(ant.Cells[i].Step.Index+1),
		)
	}
	return img
}
