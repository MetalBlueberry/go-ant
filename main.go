package main

import (
	"flag"
	"fmt"
	"go-ant/langton"
	"image/png"
	"log"
	"math"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/browser"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var steps string
var antSpeed int64

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

	// scrollSpeed := time.Duration(1)

	palette, err := colorful.SoftPalette(len(steps))
	if err != nil {
		panic(err)
	}

	ant := langton.NewAntFromString(
		langton.NewBoard(1000),
		steps)

	go runWebServer(ant, palette)

	var (
		camPos                  = pixel.ZV
		camSpeed                = 500.0
		camZoom                 = 1.0
		camZoomSpeed            = 1.2
		screenTextMargin        = pixel.V(10, -10)
		antStepCount     uint64 = 0
		antRealSpeed     uint64 = 0
	)

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
				<-time.After(time.Duration(antSpeed))
			}
			_, err := ant.Next()
			if err != nil {
				return
			}
			atomic.AddUint64(&antStepCount, 1)
		}
	}()

	loadLastPic := LastPic(ant, palette)

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
			if time.Duration(antSpeed) > time.Second*5 {
				antSpeed = int64(time.Second * 5)
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
		if win.JustPressed(pixelgl.KeyS) {
			browser.OpenURL("http://127.0.0.1:8080/pic")
		}

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		win.Clear(colornames.Black)

		loadLastPic().Draw(win, pixel.IM)

		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(100, 500), basicAtlas)

		p := message.NewPrinter(language.Spanish)

		p.Fprintf(basicTxt, "Delay between steps: %s\n", time.Duration(antSpeed))
		p.Fprintf(basicTxt, "Real Steps Per Seccond: %d\n", atomic.LoadUint64(&antRealSpeed))
		p.Fprintf(basicTxt, "Total Steps: %d\n", ant.TotalSteps())
		fmt.Fprintf(basicTxt, "Framerate: %f\n", 1.0/dt)
		fmt.Fprint(basicTxt, "Press S to save the current picture")
		win.SetMatrix(pixel.IM)
		basicTxt.Draw(win, pixel.IM.Moved(win.Bounds().Vertices()[1].Sub(basicTxt.Bounds().Vertices()[1]).Add(screenTextMargin)))

		win.Update()
	}
}

func main() {
	flag.StringVar(&steps, "steps", "LR", "Provide the sequence as L for left and R for right")
	flag.Int64Var(&antSpeed, "speed", time.Second.Nanoseconds(), "the number of nanoseconds to want between interactions. 0 for no wait")
	flag.Parse()

	pixelgl.Run(run)

}

func LastPic(ant *langton.Ant, palette []colorful.Color) func() *pixel.Sprite {
	steps := ant.TotalSteps()
	var (
		sprite *pixel.Sprite
	)
	img := langton.ToImage(ant, palette)
	pic := pixel.PictureDataFromImage(img)
	sprite = pixel.NewSprite(pic, pic.Bounds())

	return func() *pixel.Sprite {
		if ant.TotalSteps() == steps {
			return sprite
		}
		steps = ant.TotalSteps()
		img := langton.ToImage(ant, palette)
		pic := pixel.PictureDataFromImage(img)
		sprite = pixel.NewSprite(pic, pic.Bounds())
		return sprite
	}
}

func runWebServer(ant *langton.Ant, palette []colorful.Color) {
	http.HandleFunc("/pic", func(w http.ResponseWriter, r *http.Request) {
		img := langton.ToImage(ant, palette)
		err := png.Encode(w, img)
		if err != nil {
			log.Printf("error encoding picture: %s", err)
		}
		return
	})

	server := &http.Server{
		Addr: "0.0.0.0:8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
