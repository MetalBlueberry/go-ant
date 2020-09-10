package langton

import (
	"image"

	"github.com/lucasb-eyer/go-colorful"
)

func ToImage(ant *Ant, pallete []colorful.Color) image.Image {
	ant.Lock()
	defer ant.Unlock()

	r := image.Rect(
		int(ant.Dimensions.BottomLeft.X),
		int(ant.Dimensions.BottomLeft.Y),
		int(ant.Dimensions.TopRight.X),
		int(ant.Dimensions.TopRight.Y),
	)

	img := image.NewNRGBA(r)
	for i := range ant.Cells {
		if ant.Cells[i].Step.Action == ActionNone {
			continue
		}
		color := pallete[ant.Cells[i].Step.Index]

		img.Set(
			int(ant.Cells[i].X),
			int(ant.Cells[i].Y),
			color,
		)
	}
	return img
}
