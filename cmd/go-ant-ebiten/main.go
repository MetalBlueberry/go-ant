package main

import (
	"go-ant/langton"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/math/f64"
	"golang.org/x/text/message"
)

const (
	screenSize = 600
)

type Game struct {
	ant     *langton.Ant
	palette color.Palette
	camera  Camera
	world   *ebiten.Image

	properties *properties
	printer    *message.Printer
}

type properties struct {
	camSpeed       float64
	zoomSpeed      float64
	wheelZoomSpeed float64
	startDrag      image.Point
	sequence       string
	nextSequence   string

	antStepsPerSeccond float64
	antPendingSteps    float64
}

func defaultProperties() *properties {
	return &properties{
		camSpeed:       100.0,
		zoomSpeed:      1,
		wheelZoomSpeed: 5,
		sequence:       "LR",

		antStepsPerSeccond: 1,
	}

}

var previous time.Time

func (g *Game) Update(screen *ebiten.Image) error {
	if previous.IsZero() {
		previous = time.Now()
		return nil
	}
	delta := time.Since(previous).Seconds()
	previous = time.Now()

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camera.Position[0] -= g.properties.camSpeed * delta / g.camera.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camera.Position[0] += g.properties.camSpeed * delta / g.camera.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camera.Position[1] -= g.properties.camSpeed * delta / g.camera.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camera.Position[1] += g.properties.camSpeed * delta / g.camera.ZoomFactor
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.properties.startDrag = image.Pt(ebiten.CursorPosition())

	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		next := image.Pt(ebiten.CursorPosition())
		target := g.properties.startDrag.Sub(next)
		g.camera.Position[0] += float64(target.X) / g.camera.ZoomFactor
		g.camera.Position[1] += float64(target.Y) / g.camera.ZoomFactor
		g.properties.startDrag = next
	}

	_, mouseWheel := ebiten.Wheel()
	if mouseWheel > 0 {
		g.camera.ZoomFactor += g.properties.wheelZoomSpeed * delta * g.camera.ZoomFactor
	}
	if mouseWheel < 0 {
		g.camera.ZoomFactor -= g.properties.wheelZoomSpeed * delta * g.camera.ZoomFactor
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.camera.ZoomFactor -= g.properties.zoomSpeed * delta * g.camera.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.camera.ZoomFactor += g.properties.zoomSpeed * delta * g.camera.ZoomFactor
	}
	if g.camera.ZoomFactor < 0.1 {
		g.camera.ZoomFactor = 0.1
	}
	if g.camera.ZoomFactor > 20 {
		g.camera.ZoomFactor = 20
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.camera.Rotation += 10
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.camera.Rotation -= 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyKPAdd) {
		g.properties.antStepsPerSeccond *= 1.3 * (1 + delta)
	}
	if ebiten.IsKeyPressed(ebiten.KeyKPSubtract) {
		g.properties.antStepsPerSeccond /= 1.3 * (1 + delta)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.properties.nextSequence += "L"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.properties.nextSequence += "R"
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if len(g.properties.nextSequence) != 0 {
			g.properties.sequence = g.properties.nextSequence
			g.properties.nextSequence = ""
		}

		antGridSize := int64(1000)
		g.ant = langton.NewAntFromString(
			langton.NewBoard(antGridSize),
			g.properties.sequence,
		)

		p, err := colorful.HappyPalette(len(g.properties.sequence))
		if err != nil {
			panic(err)
		}
		g.palette = langton.ToPalette(p)
	}

	g.properties.antPendingSteps += g.properties.antStepsPerSeccond * delta
	steps := math.Floor(g.properties.antPendingSteps)
	g.properties.antPendingSteps = g.properties.antPendingSteps - steps
	_, err := g.ant.NextN(int(steps))
	if err != nil {
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	tmp, err := ebiten.NewImage(
		screen.Bounds().Dx(),
		screen.Bounds().Dy(),
		ebiten.FilterDefault,
	)
	if err != nil {
		panic(err)
	}

	wm := g.camera.WorldMatrix()

	g.camera.DrawAnt(g.ant, tmp, g.palette)

	cx, cy := ebiten.CursorPosition()
	mx, my := wm.Apply(float64(cx), float64(cy))
	mx = math.Floor(mx)
	my = math.Floor(my)

	cell, _ := g.ant.CellAt(langton.Point{
		int64(mx),
		int64(my),
	})

	screen.DrawImage(tmp, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrint(screen,
		g.printer.Sprintf(
			`TPS: %0.2f - FPS: %0.2f
Cell %s
Now playing "%s"
Steps x Seccond %0.2f
Total Steps %d,
Type sequence with LR and press Enter to play: "%s"`,
			ebiten.CurrentTPS(),
			ebiten.CurrentFPS(),
			cell,
			g.properties.sequence,
			g.properties.antStepsPerSeccond,
			g.ant.TotalSteps(),
			g.properties.nextSequence,
		))
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

	sequence := "LR"
	antGridSize := int64(1000)
	ant := langton.NewAntFromString(
		langton.NewBoard(antGridSize),
		sequence,
	)

	p, err := colorful.HappyPalette(len(sequence))
	if err != nil {
		panic(err)
	}

	g := &Game{
		camera: Camera{
			ViewPort:   f64.Vec2{screenSize, screenSize},
			Position:   f64.Vec2{-screenSize / 2, -screenSize / 2},
			ZoomFactor: 10,
		},
		ant:        ant,
		palette:    langton.ToPalette(p),
		properties: defaultProperties(),

		printer: message.NewPrinter(message.MatchLanguage("en")),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
