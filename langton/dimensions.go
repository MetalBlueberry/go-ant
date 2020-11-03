package langton

import (
	"fmt"
)

func NewBoard(size int64) Dimensions {
	return NewDimensions(-size, -size, size, size)
}

func NewDimensions(minX, minY, maxX, maxY int64) Dimensions {
	dim := Dimensions{
		BottomLeft: Point{
			X: minX,
			Y: minY,
		},
		TopRight: Point{
			X: maxX,
			Y: maxY,
		},
	}
	dim.Init()
	return dim
}

type Dimensions struct {
	TopRight   Point
	BottomLeft Point

	width  int64
	height int64
	Size   int64
}

type Point struct {
	X int64
	Y int64
}

func (point Point) String() string {
	return fmt.Sprintf(
		"(X: %d, Y: %d)",
		point.X,
		point.Y,
	)
}

type Direction int

const (
	DirectionTop Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
	DirectionInvalid
)

func (d Direction) Turn(action Action) Direction {
	switch action {
	case ActionTurnLeft:
		return (d + DirectionInvalid - 1) % DirectionInvalid
	case ActionTurnRight:
		return (d + DirectionInvalid + 1) % DirectionInvalid
	case ActionStraight:
		return d
	default:
		panic("Invalid action provided")
	}
}
func (d Direction) Unturn(action Action) Direction {
	switch action {
	case ActionTurnLeft:
		return (d + DirectionInvalid + 1) % DirectionInvalid
	case ActionTurnRight:
		return (d + DirectionInvalid - 1) % DirectionInvalid
	case ActionStraight:
		return d
	default:
		panic("Invalid action provided")
	}
}

func (point Point) Walk(direction Direction) Point {
	switch direction {
	case DirectionTop:
		point.Y++
	case DirectionDown:
		point.Y--
	case DirectionRight:
		point.X++
	case DirectionLeft:
		point.X--
	}

	return point
}

func (dim *Dimensions) Init() {
	dim.height = dim.TopRight.X - dim.BottomLeft.X + 1
	dim.width = dim.TopRight.Y - dim.BottomLeft.Y + 1
	dim.Size = dim.height * dim.width
}

func (dim *Dimensions) Center() Point {
	return Point{
		X: (dim.BottomLeft.X + dim.TopRight.X) / 2,
		Y: (dim.BottomLeft.Y + dim.TopRight.Y) / 2,
	}
}

func (dim *Dimensions) Width() int64 {
	return dim.width
}

func (dim *Dimensions) Height() int64 {
	return dim.height
}

func (dim *Dimensions) isPointInside(p Point) bool {
	return p.X >= dim.BottomLeft.X &&
		p.X <= dim.TopRight.X &&
		p.Y >= dim.BottomLeft.Y &&
		p.Y <= dim.TopRight.Y
}

func (dim *Dimensions) IndexOf(p Point) int {
	x := p.X - dim.BottomLeft.X
	y := p.Y - dim.BottomLeft.Y
	return int((x) + (y)*dim.height)
}

func (dim *Dimensions) String() string {
	return fmt.Sprintf("%dx%d", dim.width, dim.height)
}
