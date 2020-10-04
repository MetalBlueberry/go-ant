package main

import (
	"fmt"
	"go-ant/langton"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	"golang.org/x/image/math/f64"
)

const (
	screenSize = 300
)

type Game struct {
	ant    *langton.Ant
	palete color.Palette
	camera Camera
	world  *ebiten.Image
}

var (
	camSpeed  float64 = 100.0
	zoomSpeed float64 = 1
)

var previous time.Time

func (g *Game) Update(screen *ebiten.Image) error {
	delta := time.Since(previous).Seconds()
	previous = time.Now()

	_, err := g.ant.NextN(10000)
	if err != nil {

	}

	// screen.Fill(color.RGBA{0xff, 0, 0, 0xff})
	// bounds := screen.Bounds()
	// ebitenutil.DrawRect(screen, float64(bounds.Min.X), float64(bounds.Min.Y), float64(bounds.Dx()), float64(bounds.Dy()), colornames.Blue)
	// ebitenutil.DebugPrint(screen, bounds.String())
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camera.Position[0] -= camSpeed * delta / g.camera.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camera.Position[0] += camSpeed * delta / g.camera.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camera.Position[1] -= camSpeed * delta / g.camera.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camera.Position[1] += camSpeed * delta / g.camera.ZoomFactor
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.camera.ZoomFactor -= zoomSpeed * delta
		if g.camera.ZoomFactor < 0.1 {
			g.camera.ZoomFactor = 0.1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.camera.ZoomFactor += zoomSpeed * delta
		if g.camera.ZoomFactor > 5 {
			g.camera.ZoomFactor = 5
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.camera.Rotation += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.camera.Reset()
	}
	ant := langton.ToImage(g.ant, g.palete, 1)
	img, err := ebiten.NewImageFromImage(ant, ebiten.FilterDefault)
	if err != nil {
		return err
	}
	g.world = img
	// scale := float64(screen.Bounds().Dx()) / float64(img.Bounds().Dx())
	// g.world.DrawImage(img, &ebiten.DrawImageOptions{
	// 	GeoM: ebiten.ScaleGeo(scale, scale),
	// })
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.world == nil {
		return
	}
	drawBounds(g.world, 10, colornames.White)
	g.camera.Render(g.world, screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func drawBounds(dst *ebiten.Image, size float64, clr color.Color) {
	bounds := dst.Bounds()
	x := float64(bounds.Min.X)
	y := float64(bounds.Min.Y)
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	drawHLine(dst, x, y+h-size, w, size, clr)
	drawHLine(dst, x, y, w, size, clr)
	drawVLine(dst, x+size, y, h, size, clr)
	drawVLine(dst, x+w, y, h, size, clr)
}

func drawHLine(dst *ebiten.Image, x, y, len, size float64, clr color.Color) {
	ebitenutil.DrawRect(dst, x, y+size, len, -size, clr)
}
func drawVLine(dst *ebiten.Image, x, y, len, size float64, clr color.Color) {
	ebitenutil.DrawRect(dst, x-size, y, size, len, clr)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenSize, screenSize
}

func main() {
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizable(true)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetMaxTPS(25)

	sequence := "LLLRRRRRRLLL"
	antGridSize := int64(1000)
	ant := langton.NewAntFromString(
		langton.NewBoard(antGridSize),
		sequence,
	)

	p, err := colorful.HappyPalette(len(sequence))
	if err != nil {
		panic(err)
	}

	antBoardSize := (float64(ant.Dimensions.Width()) * 1) / 2

	g := &Game{
		// canvas: image.NewRGBA(image.Rect(0, 0, screenSize, screenSize)),
		camera: Camera{
			ViewPort:   f64.Vec2{screenSize, screenSize},
			Position:   f64.Vec2{antBoardSize - screenSize/2, antBoardSize - screenSize/2},
			ZoomFactor: 1,
		},
		ant:    ant,
		palete: langton.ToPalette(p),
	}
	// world, err := ebiten.NewImage(screenSize, screenHeight, ebiten.FilterDefault)
	// if err != nil {
	// 	panic(err)
	// }
	// g.world = world

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
