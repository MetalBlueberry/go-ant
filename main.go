package main

import (
	"flag"
	"fmt"
	"go-ant/langton"
	"math"
	"sync/atomic"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var steps string

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

	antSpeed := time.Second

	// scrollSpeed := time.Duration(1)

	palette, err := colorful.SoftPalette(len(steps))
	if err != nil {
		panic(err)
	}

	ant := langton.NewAntFromString(steps)

	var (
		camPos                  = pixel.ZV
		camSpeed                = 500.0
		camZoom                 = 1.0
		camZoomSpeed            = 1.2
		screenTextMargin        = pixel.V(10, -10)
		antStepCount     uint64 = 0
		antRealSpeed     uint64 = 0
	)

	imd := imdraw.New(nil)
	go func() {
		ticker := time.NewTicker(time.Duration(time.Second))
		for {
			<-ticker.C
			atomic.StoreUint64(&antRealSpeed, atomic.LoadUint64(&antStepCount))
			atomic.StoreUint64(&antStepCount, 0)
		}
	}()
	go func() {
		for {
			if antSpeed > 0 {
				<-time.After(antSpeed)
			}
			ant.Next()
			atomic.AddUint64(&antStepCount, 1)
		}
	}()

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if win.Pressed(pixelgl.KeyKPAdd) {
			antSpeed = antSpeed / 2
			if antSpeed < 0 {
				antSpeed = 0
			}
		}
		if win.Pressed(pixelgl.KeyKPSubtract) {
			if antSpeed == 0 {
				antSpeed++
			}
			antSpeed = antSpeed * 2
			if antSpeed > time.Second*5 {
				antSpeed = time.Second * 5
			}
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
			imd.Color = palette[cell.Step.Index]
			imd.Push(pixel.V(float64(cell.X)*(cellSize), float64(cell.Y)*(cellSize)))
			imd.Push(pixel.V(float64(cell.X)*(cellSize)+cellSize, float64(cell.Y)*(cellSize)+cellSize))
			imd.Rectangle(0)
		}
		ant.Unlock()

		imd.Draw(win)

		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(100, 500), basicAtlas)

		p := message.NewPrinter(language.Spanish)

		p.Fprintf(basicTxt, "Delay between steps: %s\n", antSpeed)
		p.Fprintf(basicTxt, "Real Steps Per Seccond: %d\n", atomic.LoadUint64(&antRealSpeed))
		fmt.Fprintf(basicTxt, "Framerate: %f\n", 1.0/dt)
		win.SetMatrix(pixel.IM)
		basicTxt.Draw(win, pixel.IM.Moved(win.Bounds().Vertices()[1].Sub(basicTxt.Bounds().Vertices()[1]).Add(screenTextMargin)))

		win.Update()
	}
}

func main() {
	flag.StringVar(&steps, "steps", "LR", "Provide the sequence as L for left and R for right")
	flag.Parse()
	pixelgl.Run(run)
}
