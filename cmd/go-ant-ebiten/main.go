package main

import (
	"fmt"
	"go-ant/langton"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/math/f64"
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
}

type properties struct {
	camSpeed  float64
	zoomSpeed float64
}

func defaultProperties() *properties {
	return &properties{
		camSpeed:  100.0,
		zoomSpeed: 1,
	}

}

var previous time.Time

func (g *Game) Update(screen *ebiten.Image) error {
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

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.camera.ZoomFactor -= g.properties.zoomSpeed * delta * g.camera.ZoomFactor
		if g.camera.ZoomFactor < 0.1 {
			g.camera.ZoomFactor = 0.1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.camera.ZoomFactor += g.properties.zoomSpeed * delta * g.camera.ZoomFactor
		if g.camera.ZoomFactor > 10 {
			g.camera.ZoomFactor = 10
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.camera.Rotation += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.camera.Rotation -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.camera.Reset()
	}

	_, err := g.ant.NextN(1000)
	if err != nil {
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	wm := g.camera.WorldMatrix()

	DrawImage(g.ant, screen, g.palette, wm)

	cx, cy := ebiten.CursorPosition()
	mx, my := wm.Apply(float64(cx), float64(cy))
	mx = math.Floor(mx)
	my = math.Floor(my)

	cell, _ := g.ant.CellAt(langton.Point{
		int64(mx),
		int64(my),
	})

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"TPS: %0.2f\nFPS: %0.2f\nmouse %s\ncell %s",
			ebiten.CurrentTPS(),
			ebiten.CurrentFPS(),
			image.Pt(int(mx), int(my)),
			cell,
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
	// ebiten.SetMaxTPS(25)

	sequence := "LLLRRRRRRLLL"
	antGridSize := int64(500)
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
			ZoomFactor: 5,
		},
		ant:        ant,
		palette:    langton.ToPalette(p),
		properties: defaultProperties(),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func DrawImage(ant *langton.Ant, screen *ebiten.Image, palette color.Palette, geo ebiten.GeoM) {
	bounds := screen.Bounds()
	for sx := 0; sx < bounds.Dx(); sx++ {
		for sy := 0; sy < bounds.Dy(); sy++ {
			x, y := geo.Apply(float64(sx), float64(sy))
			x = math.Floor(x)
			y = math.Floor(y)
			cell, err := ant.CellAt(langton.Point{int64(x), int64(y)})
			if err != nil {
				continue
			}
			screen.Set(sx, sy, palette[cell.Step.Index+1])
		}
	}
	// for i := range ant.Cells {
	// 	if ant.Cells[i].Step.Action == langton.ActionNone {
	// 		continue
	// 	}

	// 	x, y := geo.Apply(
	// 		// float64(sx),
	// 		// float64(sy),
	// 		float64(ant.Cells[i].X),
	// 		float64(ant.Cells[i].Y),
	// 	)
	// 	p := image.Pt(int(x), int(y))
	// 	if p.In(screen.Bounds()) {
	// 		screen.Set(
	// 			// int(ant.Cells[i].X),
	// 			// int(ant.Cells[i].Y),
	// 			int(x),
	// 			int(y),
	// 			// int((ant.Cells[i].X+ant.Dimensions.Width()/2)*int64(cellSize)+int64(sx)),
	// 			// int((ant.Cells[i].Y+ant.Dimensions.Height()/2)*int64(cellSize)+int64(sy)),
	// 			palette[ant.Cells[i].Step.Index+1],
	// 		)
	// 	}
	// }

	// if cellSize > 5 {
	// 	cell := ant.Position
	// 	for sx := 0; sx < cellSize; sx++ {
	// 		for sy := 0; sy < cellSize; sy++ {
	// 			radius := cellSize / 2
	// 			if distance2From(sx, sy, radius, radius) <= (radius-1)*(radius-1) {
	// 				var color color.Color
	// 				switch {
	// 				case ant.Direction == langton.DirectionLeft && sx < radius && sy == radius:
	// 					color = colornames.Red
	// 				case ant.Direction == langton.DirectionRight && sx > radius && sy == radius:
	// 					color = colornames.Red
	// 				case ant.Direction == langton.DirectionTop && sx == radius && sy > radius:
	// 					color = colornames.Red
	// 				case ant.Direction == langton.DirectionDown && sx == radius && sy < radius:
	// 					color = colornames.Red
	// 				default:
	// 					color = colornames.Black
	// 				}

	// 				screen.Set(
	// 					int((cell.X+ant.Dimensions.Width()/2)*int64(cellSize)+int64(sx)),
	// 					int((cell.Y+ant.Dimensions.Height()/2)*int64(cellSize)+int64(sy)),
	// 					color,
	// 				)
	// 			}
	// 		}
	// 	}
	// }
}

func distance2From(ax, ay, bx, by int) int {
	x := bx - ax
	y := by - ay
	return x*x + y*y
}
