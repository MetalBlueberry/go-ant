package langton

import (
	"fmt"
)

// NewBoard creates a squared board with side equals to size/2+1
// Just because it is handy
func NewBoard(size int64) Dimensions {
	return NewDimensions(-size, -size, size, size)
}

// Creates a board with the given coordinates
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

// Init must be always called after creation, it precalculates some internal values
func (dim *Dimensions) Init() {
	dim.height = dim.TopRight.X - dim.BottomLeft.X + 1
	dim.width = dim.TopRight.Y - dim.BottomLeft.Y + 1
	dim.Size = dim.height * dim.width
}

// Dimensions represent the area that the ant can explore
type Dimensions struct {
	TopRight   Point
	BottomLeft Point

	width  int64
	height int64
	Size   int64
}

// Point is a x y pair
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

// Direction is an enum used to track the ant direction
type Direction int

const (
	// DirectionTop moves up
	DirectionTop Direction = iota
	// DirectionRight moves right
	DirectionRight
	// DirectionDown moves down
	DirectionDown
	// DirectionLeft moves left
	DirectionLeft
	// DirectionInvalid is an invalid direction
	DirectionInvalid
)

// Turn changes the Direction based on the provided Action
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

// Unturn performs the opposite operation to Turn
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

// Walk moves the ant in the given Direction, returns the final position
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

// Center returns the center of the Dimensions
func (dim *Dimensions) Center() Point {
	return Point{
		X: (dim.BottomLeft.X + dim.TopRight.X) / 2,
		Y: (dim.BottomLeft.Y + dim.TopRight.Y) / 2,
	}
}

// Width returns the width of the Dimensions
func (dim *Dimensions) Width() int64 {
	return dim.width
}

// Height returns the Height of the Dimensions
func (dim *Dimensions) Height() int64 {
	return dim.height
}

// isPointInside returns true if the point is inside the Dimensions
func (dim *Dimensions) isPointInside(p Point) bool {
	return p.X >= dim.BottomLeft.X &&
		p.X <= dim.TopRight.X &&
		p.Y >= dim.BottomLeft.Y &&
		p.Y <= dim.TopRight.Y
}

// indexOf returns the index of a given point. Is to map a 2 dimensions into 1 slice
func (dim *Dimensions) indexOf(p Point) int {
	x := p.X - dim.BottomLeft.X
	y := p.Y - dim.BottomLeft.Y
	return int((x) + (y)*dim.height)
}

// String returns a string representation
func (dim *Dimensions) String() string {
	return fmt.Sprintf("%dx%d", dim.width, dim.height)
}
