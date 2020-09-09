package main

import (
	"fmt"
	"go-ant/langoth"
	"log"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	antSpeed := 10 * time.Millisecond

	scrollSpeed := 250 * time.Microsecond

	// steps := langoth.Steps{
	// 	langoth.Step{
	// 		Action: langoth.ActionTurnRight,
	// 	},
	// 	langoth.Step{
	// 		Action: langoth.ActionTurnLeft,
	// 	},
	// 	langoth.Step{
	// 		Action: langoth.ActionTurnLeft,
	// 	},
	// 	langoth.Step{
	// 		Action: langoth.ActionTurnRight,
	// 	},
	// }

	steps := langoth.StepsAwesome2
	palette, err := colorful.SoftPalette(len(steps))
	if err != nil {
		panic(err)
	}
	for i := range palette {
		steps[i].Color = palette[i]
	}

	ant := langoth.NewAnt(steps...)

	var (
		camPos           = pixel.ZV
		camSpeed         = 500.0
		camZoom          = 1.0
		camZoomSpeed     = 1.2
		screenTextMargin = pixel.V(10, -10)
	)

	imd := imdraw.New(nil)

	go func() {
		for {
			<-time.After(antSpeed)
			ant.Lock()
			ant.Next()
			ant.Unlock()
		}
	}()

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if win.Pressed(pixelgl.KeyKPAdd) {
			log.Println(antSpeed)
			antSpeed += scrollSpeed
		}
		if win.Pressed(pixelgl.KeyKPSubtract) {
			log.Println(antSpeed)
			antSpeed -= scrollSpeed
		}

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		win.Clear(colornames.Black)

		cellSize := 5.0

		imd.Clear()
		ant.Lock()
		for _, cell := range ant.Cells {
			imd.Color = cell.Step.Color
			imd.Push(pixel.V(float64(cell.X)*(cellSize), float64(cell.Y)*(cellSize)))
			imd.Push(pixel.V(float64(cell.X)*(cellSize)+cellSize, float64(cell.Y)*(cellSize)+cellSize))
			imd.Rectangle(0)
		}
		ant.Unlock()

		imd.Draw(win)

		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(100, 500), basicAtlas)

		fmt.Fprintf(basicTxt, "Speed: %s\n", antSpeed)
		fmt.Fprintf(basicTxt, "Framerate: %f\n", 1.0/dt)
		win.SetMatrix(pixel.IM)
		basicTxt.Draw(win, pixel.IM.Moved(win.Bounds().Vertices()[1].Sub(basicTxt.Bounds().Vertices()[1]).Add(screenTextMargin)))

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
