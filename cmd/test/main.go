package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"os"
)

// CreateBasicGif creates a GIF image with the given width and height.
// It uses white background and a black pixel in the middle of the image.
func CreateBasicGif(out io.Writer, width, height int) {

	palette := []color.Color{color.White, color.Black}
	rect := image.Rect(0, 0, width, height)
	img := image.NewPaletted(rect, palette)

	img.SetColorIndex(width/2, height/2, 2)

	anim := gif.GIF{Delay: []int{0}, Image: []*image.Paletted{img}}

	gif.EncodeAll(out, &anim)
}

func main() {
	file, err := os.Create("tmp.gif")
	if err != nil {
		panic(err)
	}
	CreateBasicGif(file, 1000, 1000)

}
