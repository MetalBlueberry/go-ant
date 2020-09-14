package main

import (
	"go-ant/langton"
	"image"
	"image/gif"
	"log"
	"os"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	steps := "RLLLLRRRLLL"
	ant := langton.NewAntFromString(
		langton.NewBoard(150),
		steps,
	)
	colorfulPalette, err := colorful.SoftPalette(len(steps))
	if err != nil {
		panic(err)
	}
	palette := langton.ToPalette(colorfulPalette)

	var (
		framesXSeccond  int64 = 40
		duration              = 20 * time.Second
		updates               = 271433
		frames                = int(framesXSeccond * duration.Milliseconds() / 1000)
		updatesPerFrame       = updates / frames
	)

	images := make([]*image.Paletted, 0, frames) // The successive images.
	delay := make([]int, 0, frames)
	disposal := make([]byte, 0, frames)

	delayValue := int(duration.Seconds()) * 100 / frames
	if delayValue == 0 {
		log.Print("WARNING: fps is too high")
		delayValue = 1
	}

	optimizer := GifFrameOptimizer()

	for frame := 0; frame < frames; frame++ {
		log.Printf("frame %d of %d", frame, frames)
		err := Calculate(ant, updatesPerFrame)
		if err != nil {
			log.Printf("Bound reached at step %d, use that value as updates next time or increase image size", ant.TotalSteps())
			break
		}

		img := langton.ToImage(ant, palette, 3)
		optimizer(img)

		images = append(images, img)
		delay = append(delay, delayValue)
		disposal = append(disposal, gif.DisposalNone)

		// file, err := os.Create(fmt.Sprintf("out/frame_%d.png", frame))
		// if err != nil {
		// 	panic(err)
		// }

		// err = png.Encode(file, img)
		// if err != nil {
		// 	panic(err)
		// }
		// err = file.Close()
		// if err != nil {
		// 	panic(err)
		// }
	}

	out := &gif.GIF{
		Delay:    delay,
		Image:    images,
		Disposal: disposal,
	}
	file, err := os.Create("out.gif")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = gif.EncodeAll(file, out)
	if err != nil {
		panic(err)
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
