package main

import (
	"fmt"
	"go-ant/langoth"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1 * time.Millisecond)

	ant := langoth.NewAnt(
		langoth.Step{
			Color:  colornames.Red,
			Action: langoth.ActionTurnLeft,
		},
		langoth.Step{
			Color:  colornames.Blue,
			Action: langoth.ActionTurnRight,
		},
	)

	for !win.Closed() {
		<-ticker.C
		cam := pixel.IM.Moved(win.Bounds().Center())
		win.SetMatrix(cam)
		ant.Next()
		fmt.Print(ant)

		win.Clear(colornames.Black)

		imd := imdraw.New(nil)

		cellSize := 5.0

		for _, cell := range ant.Cells {
			imd.Color = cell.Step.Color
			imd.Push(pixel.V(float64(cell.X)*(cellSize), float64(cell.Y)*(cellSize)))
			imd.Push(pixel.V(float64(cell.X)*(cellSize)+cellSize, float64(cell.Y)*(cellSize)+cellSize))
			imd.Rectangle(0)
		}

		imd.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
