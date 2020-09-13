package main

import (
	"fmt"
	"go-ant/langton"
	"image"
	"image/gif"
	"image/png"
	"os"

	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	ant := langton.NewAntFromString(
		langton.NewBoard(5),
		"LRR",
	)
	colorfulPalette, err := colorful.SoftPalette(3)
	if err != nil {
		panic(err)
	}
	palette := langton.ToPalette(colorfulPalette)

	var (
		frames = 200
	)

	images := make([]*image.Paletted, 0, frames) // The successive images.
	delay := make([]int, 0, frames)

	for frame := 0; frame < frames; frame++ {
		_, err := ant.Next()
		if err != nil {
			break
		}
		img := langton.ToImage(ant, palette)
		images = append(images, img)
		delay = append(delay, 10)

		file, err := os.Create(fmt.Sprintf("out/frame_%d.png", frame))
		if err != nil {
			panic(err)
		}

		err = png.Encode(file, img)
		if err != nil {
			panic(err)
		}
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}

	out := &gif.GIF{
		Delay: delay,
		Image: images,
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
