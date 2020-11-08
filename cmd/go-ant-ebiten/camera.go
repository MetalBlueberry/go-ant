package main

import (
	"fmt"
	"go-ant/langton"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	ViewPort        f64.Vec2
	Position        f64.Vec2
	InitialPosition f64.Vec2
	ZoomFactor      float64
	Rotation        float64
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %f, S: %f",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] * 0.5,
		c.ViewPort[1] * 0.5,
	}
}

func (c *Camera) Apply(x, y int) (int, int) {
	wm := c.WorldMatrix()

	mx, my := wm.Apply(float64(x), float64(y))
	return int(math.Floor(mx)), int(math.Floor(my))
}

func (c *Camera) WorldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.InitialPosition[0], -c.InitialPosition[1])
	m.Translate(-c.Position[0], -c.Position[1])
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		c.ZoomFactor,
		c.ZoomFactor,
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) error {
	return screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.WorldMatrix(),
	})
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.WorldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		return math.NaN(), math.NaN()
	}
}

func MapVector(origin f64.Vec2, geoM ebiten.GeoM) (f64.Vec2, f64.Vec2, f64.Vec2) {
	originx, originy := geoM.Apply(origin[0], origin[1])
	xx, xy := geoM.Apply(origin[0]+1, origin[1])
	yx, yy := geoM.Apply(origin[0], origin[1]+1)

	return f64.Vec2{originx, originy},
		f64.Vec2{xx - originx, xy - originy},
		f64.Vec2{yx - originx, yy - originy}
}

func cache(ant *langton.Ant) func(p langton.Point) (*langton.Cell, error) {
	type pair struct {
		cell *langton.Cell
		err  error
	}
	cache := make(map[langton.Point]pair)
	return func(p langton.Point) (*langton.Cell, error) {
		if v, ok := cache[p]; ok {
			return v.cell, v.err
		}
		c, err := ant.CellAt(p)
		cache[p] = pair{c, err}
		return c, err
	}
}

func (c *Camera) DrawAnt(ant *langton.Ant, screen *ebiten.Image, palette color.Palette) {
	bounds := screen.Bounds()
	geo := c.WorldMatrix()
	geo.Invert()

	origin, vectorx, vectory := MapVector(f64.Vec2{
		float64(bounds.Min.X),
		float64(bounds.Min.Y),
	}, geo)

	dx := float64(bounds.Dx())
	dy := float64(bounds.Dy())
	for sx := 0.0; sx < dx; sx++ {
		xx, xy := sx*vectorx[0], sx*vectorx[1]
		for sy := 0.0; sy < dy; sy++ {
			yx, yy := sy*vectory[0], sy*vectory[1]
			x, y := xx+yx+origin[0], xy+yy+origin[1]
			x = math.Floor(x)
			y = math.Floor(y)

			cell, err := ant.CellAt(langton.Point{int64(x), int64(y)})
			if err != nil {
				continue
			}
			screen.Set(int(sx), int(sy), palette[cell.Step.Index+1])
		}
	}

	// wm := c.WorldMatrix()

	// antx, anty := wm.Apply(float64(cell.X), float64(cell.Y))
	// antmaxx, antmaxy := wm.Apply(float64(cell.X+1), float64(cell.Y+1))
	// antCenterX := math.Ceil((antx + antmaxx) / 2)
	// antCenterY := math.Ceil((anty + antmaxy) / 2)

	antSize := 1 / (math.Sqrt(vectorx[0]*vectorx[0] + vectorx[1]*vectorx[1]))

	if antSize >= 4 {
		cell := ant.Position
		x, y := float64(cell.X), float64(cell.Y)
		wm := c.WorldMatrix()
		sx, sy := wm.Apply(x, y)

		geoM := ebiten.GeoM{}
		geoM.Rotate(c.Rotation * 2 * math.Pi / 360)

		geoM.Translate(-float64(AntImage.Bounds().Dx())/2.0, -float64(AntImage.Bounds().Dx())/2.0)
		switch ant.Direction {
		case langton.DirectionRight:
			geoM.Rotate(math.Pi / 2)
		case langton.DirectionTop:
			geoM.Rotate(math.Pi)
		case langton.DirectionLeft:
			geoM.Rotate(3 * math.Pi / 2)
		}

		switch cell.Step.Action {
		case langton.ActionTurnRight:
			geoM.Rotate(-math.Pi / 2)
		case langton.ActionTurnLeft:
			geoM.Rotate(math.Pi / 2)
		}
		geoM.Translate(float64(AntImage.Bounds().Dx())/2.0, float64(AntImage.Bounds().Dx())/2.0)

		geoM.Scale(antSize/float64(AntImage.Bounds().Dx()), antSize/float64(AntImage.Bounds().Dy()))
		geoM.Translate(sx, sy)

		screen.DrawImage(AntImage, &ebiten.DrawImageOptions{
			GeoM: geoM,
		})
	}
}

func distance2From(ax, ay, bx, by float64) float64 {
	x := bx - ax
	y := by - ay
	return x*x + y*y
}
