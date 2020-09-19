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
		frames          int
		iterations      int
		steps           string
		outFile         string
		pixelSize       int
		area            int64
		durationSeconds int
		open            bool
		lastFrameDelay  int
	)

	flag.StringVar(&steps, "steps", "LR", "Ant step sequence")
	flag.StringVar(&outFile, "out", "out.gif", "output file")
	flag.IntVar(&iterations, "iterations", 10927, "Total number of ant iterations")
	flag.IntVar(&frames, "frames", 200, "total gif frames")
	flag.IntVar(&pixelSize, "pixel-size", 6, "determines the final image size by multiplying this value by the area")
	flag.IntVar(&durationSeconds, "duration", 10, "gif duration in seconds")
	flag.Int64Var(&area, "area", 70, "size in cells for the ant to walk")
	flag.BoolVar(&open, "open", false, "open the output in a browser")
	flag.IntVar(&lastFrameDelay, "last-frame-delay", 1000, "milliseconds for the last frame")
	flag.Parse()

	var (
		duration           = time.Duration(durationSeconds) * time.Second
		updatesPerFrame    = iterations / frames
		delayBetweenFrames = int(duration.Seconds()) * 100 / frames
	)

	if updatesPerFrame == 0 {
		log.Println("WARNING: less than one iteration per frame, increase iterations or reduce frames")
		log.Println("setting updates per frame to 1")
		updatesPerFrame = 1
	}

	if delayBetweenFrames == 0 {
		log.Print("WARNING: duration is too small for the number of frames, reduce frames or increase duration")
		log.Println("setting delay between frames to 1")
		delayBetweenFrames = 1
	}

	log.Printf("INFO: frame rate %f, updates per frame %d", 100/float64(delayBetweenFrames), updatesPerFrame)

	ant := langton.NewAntFromString(
		langton.NewBoard(area/2),
		steps,
	)

	colorfulPalette, err := colorful.SoftPalette(len(steps))
	if err != nil {
		panic(err)
	}
	palette := langton.ToPalette(colorfulPalette)

	images := make([]*image.Paletted, 0, frames)
	delay := make([]int, 0, frames)
	disposal := make([]byte, 0, frames)

	bar := progressbar.Default(int64(frames), "Calculating")
	optimizer := GifFrameOptimizer()
	for frame := 0; frame < frames; frame++ {
		bar.Add(1)
		err := Calculate(ant, updatesPerFrame)

		img := langton.ToImage(ant, palette, pixelSize)
		optimizer(img)

		images = append(images, img)
		delay = append(delay, delayBetweenFrames)
		disposal = append(disposal, gif.DisposalNone)

		if err != nil {
			log.Printf("Bound reached at step %d, use that value as updates next time or increase image size", ant.TotalSteps())
			break
		}
	}

	// Last frame stays for a seccond
	delay[len(delay)-1] = lastFrameDelay

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

	if open {
		browser.OpenFile(outFile)
	}
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
