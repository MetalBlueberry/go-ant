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
	ViewPort   f64.Vec2
	Position   f64.Vec2
	ZoomFactor float64
	Rotation   float64
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
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
	m.Translate(-c.Position[0], -c.Position[1])
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		c.ZoomFactor,
		c.ZoomFactor,
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	m.Invert()
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

func (c *Camera) DrawAnt(ant *langton.Ant, screen *ebiten.Image, palette color.Palette) {
	bounds := screen.Bounds()
	for sx := 0; sx < bounds.Dx(); sx++ {
		for sy := 0; sy < bounds.Dy(); sy++ {
			x, y := c.Apply(sx, sy)
			cell, err := ant.CellAt(langton.Point{int64(x), int64(y)})
			if err != nil {
				continue
			}
			screen.Set(sx, sy, palette[cell.Step.Index+1])
		}
	}
}
