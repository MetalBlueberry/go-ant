package main

import (
	"flag"
	"go-ant/langton"
	"image"
	"image/gif"
	"io"
	"log"
	"os"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/browser"

	"github.com/schollz/progressbar/v3"
)

func main() {

	var (
		framesXSeccond   int64
		iterations       int
		steps            string
		outFile          string
		pixelSize        int
		area             int64
		durationSecconds int
	)

	flag.StringVar(&steps, "steps", "LR", "Ant step sequence")
	flag.StringVar(&outFile, "out", "out.gif", "output file")
	flag.IntVar(&iterations, "iterations", 271433, "Total number of ant iterations")
	flag.Int64Var(&framesXSeccond, "frames-x-seccond", 15, "frames per seccond of the final gif")
	flag.IntVar(&pixelSize, "pixel-size", 3, "determines the final image size by multiplying this value by the area")
	flag.IntVar(&durationSecconds, "duration", 20, "gif duration in seconds")
	flag.Int64Var(&area, "area", 130, "size in cells for the ant to walk")
	flag.Parse()

	var (
		duration        = time.Duration(durationSecconds) * time.Second
		frames          = int(framesXSeccond * duration.Milliseconds() / 1000)
		updatesPerFrame = iterations / frames
	)

	ant := langton.NewAntFromString(
		langton.NewBoard(area/2),
		steps,
	)
	colorfulPalette, err := colorful.SoftPalette(len(steps))
	if err != nil {
		panic(err)
	}
	palette := langton.ToPalette(colorfulPalette)

	images := make([]*image.Paletted, 0, frames) // The successive images.
	delay := make([]int, 0, frames)
	disposal := make([]byte, 0, frames)

	delayValue := int(duration.Seconds()) * 100 / frames
	if delayValue == 0 {
		log.Print("WARNING: fps is too high")
		delayValue = 1
	}

	bar := progressbar.Default(int64(frames), "Calculating")
	optimizer := GifFrameOptimizer()
	for frame := 0; frame < frames; frame++ {
		bar.Add(1)
		err := Calculate(ant, updatesPerFrame)
		if err != nil {
			log.Printf("Bound reached at step %d, use that value as updates next time or increase image size", ant.TotalSteps())
			break
		}

		img := langton.ToImage(ant, palette, pixelSize)
		optimizer(img)

		images = append(images, img)
		delay = append(delay, delayValue)
		disposal = append(disposal, gif.DisposalNone)

	}

	// Last frame stays for a seccond
	delay[len(delay)-1] = 100

	out := &gif.GIF{
		Delay:    delay,
		Image:    images,
		Disposal: disposal,
	}
	file, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encodeBar := progressbar.DefaultBytes(-1, "Saving..")
	err = gif.EncodeAll(io.MultiWriter(file, encodeBar), out)
	if err != nil {
		panic(err)
	}

	browser.OpenFile(outFile)
}

func Calculate(ant *langton.Ant, steps int) error {
	for i := 0; i < steps; i++ {
		_, err := ant.Next()
		if err != nil {
			return err
		}
	}
	return nil
}

// GifFrameOptimizer turns repeated pixels to transparent to the final gif size is minimal.
func GifFrameOptimizer() func(img *image.Paletted) {
	var currentImage *image.Paletted

	return func(img *image.Paletted) {
		if currentImage == nil {
			currentImage = &image.Paletted{}
			currentImage.Palette = img.Palette
			currentImage.Rect = img.Rect
			currentImage.Stride = img.Stride
			currentImage.Pix = make([]uint8, len(img.Pix))
			copy(currentImage.Pix, img.Pix)
			return
		}

		for i := range img.Pix {
			if img.Pix[i] == currentImage.Pix[i] {
				img.Pix[i] = 0
			} else {
				currentImage.Pix[i] = img.Pix[i]
			}
		}
	}
}
