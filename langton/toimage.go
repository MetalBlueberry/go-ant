package langton

import (
	"image"
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
)

// ToImage generates a image.Paletted with the current ant state.
// The cell size is in pixels
// If the cell size is bigger than 5, the ant will be drawn as a black dot
func ToImage(ant *Ant, palette color.Palette, cellSize int) *image.Paletted {

	r := image.Rect(
		0,
		0,
		int(ant.Dimensions.width)*cellSize,
		int(ant.Dimensions.height)*cellSize,
	)
	palette = append(palette, colornames.Black, colornames.Red)
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

	black := len(palette) - 2
	red := len(palette) - 1
	if cellSize > 5 {
		cell := ant.Position
		for sx := 0; sx < cellSize; sx++ {
			for sy := 0; sy < cellSize; sy++ {
				radius := cellSize / 2
				if distance2From(sx, sy, radius, radius) <= (radius-1)*(radius-1) {
					var color int
					switch {
					case ant.Direction == DirectionLeft && sx < radius && sy == radius:
						color = red
					case ant.Direction == DirectionRight && sx > radius && sy == radius:
						color = red
					case ant.Direction == DirectionTop && sx == radius && sy > radius:
						color = red
					case ant.Direction == DirectionDown && sx == radius && sy < radius:
						color = red
					default:
						color = black
					}

					img.SetColorIndex(
						int((cell.X+ant.Dimensions.width/2)*int64(cellSize)+int64(sx)),
						int((cell.Y+ant.Dimensions.height/2)*int64(cellSize)+int64(sy)),
						uint8(color),
					)
				}
			}
		}
	}
	return img
}

func distance2From(ax, ay, bx, by int) int {
	x := bx - ax
	y := by - ay
	return x*x + y*y
}

// ToPalette is an utility function to use colorful.Color arrays as a color.Palette
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
