package langton

import (
	"errors"
	"strings"
	"sync"
)

type Ant struct {
	Cells     []Cell
	Position  *Cell
	Direction Direction

	Dimensions Dimensions

	steps      []Step
	totalSteps int64
	sync.Locker
}

type Dimensions struct {
	TopRight   Point
	BottomLeft Point

	width  int64
	height int64
	size   int64
}

type Point struct {
	X int64
	Y int64
}

type Cell struct {
	Point
	Step Step
}

type Steps []Step

type Step struct {
	Index  int
	Action Action

	nextIndex int
}

type Action rune

const (
	ActionNone      Action = 0
	ActionTurnLeft         = 'L'
	ActionTurnRight        = 'R'
	ActionStright          = 'S'
)

type Direction int

const (
	DirectionTop Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
	DirectionInvalid
)

var (
	StepsSimple  Steps = StepsFromString("LR")
	StepsAwesome Steps = StepsFromString("RLLLLRRRLLL")
)

func (d Direction) Turn(action Action) Direction {
	switch action {
	case ActionTurnLeft:
		return (d + DirectionInvalid - 1) % DirectionInvalid
	case ActionTurnRight:
		return (d + DirectionInvalid + 1) % DirectionInvalid
	case ActionStright:
		return d
	default:
		panic("Invalid action provided")
	}
}

func (ant *Ant) TotalSteps() int64 {
	ant.Lock()
	defer ant.Unlock()
	return ant.totalSteps
}

func (ant *Ant) Next() (*Cell, error) {
	ant.Lock()
	defer ant.Unlock()
	ant.totalSteps++

	ant.Direction = ant.Direction.Turn(ant.Position.Step.Action)

	ant.Position.UpdateNextStep(ant.steps)

	nextPoint := ant.Position.Point.Walk(ant.Direction)

	previousPosition := ant.Position
	nextPosition, err := ant.EnsureCellAt(nextPoint)
	if err != nil {
		return previousPosition, err
	}
	ant.Position = nextPosition
	return previousPosition, nil
}

func (cell *Cell) UpdateNextStep(steps []Step) {
	cell.Step = steps[cell.Step.nextIndex]
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

func (ant *Ant) EnsureCellAt(position Point) (*Cell, error) {
	if !ant.Dimensions.isPointInside(position) {
		return nil, errors.New("Next step is out of bounds")
	}
	posIndex := ant.Dimensions.IndexOf(position)
	cell := &ant.Cells[posIndex]
	if cell.Step.Action == ActionNone {
		ant.Cells[posIndex] = Cell{
			Point: position,
			Step:  ant.steps[0],
		}
		cell = &ant.Cells[posIndex]
	}
	return cell, nil
}

func (steps Steps) Numerate() {
	for i := 0; i < len(steps); i++ {
		steps[i].Index = i
		steps[i].nextIndex = i + 1
	}
	steps[len(steps)-1].nextIndex = 0
}

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

func (dim *Dimensions) Init() {
	dim.height = dim.TopRight.X - dim.BottomLeft.X + 1
	dim.width = dim.TopRight.Y - dim.BottomLeft.Y + 1
	dim.size = dim.height * dim.width
}

func (dim *Dimensions) Center() Point {
	return Point{
		X: (dim.BottomLeft.X + dim.TopRight.X) / 2,
		Y: (dim.BottomLeft.Y + dim.TopRight.Y) / 2,
	}
}

func (dim Dimensions) isPointInside(p Point) bool {
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

func StepsFromString(steps string) Steps {
	out := make(Steps, len(steps), len(steps))
	for i, c := range steps {
		switch c {
		case rune(ActionTurnLeft):
			out[i] = Step{
				Action: ActionTurnLeft,
			}
		case rune(ActionTurnRight):
			out[i] = Step{
				Action: ActionTurnRight,
			}
		}
	}
	return out
}

func NewAntFromString(dimensions Dimensions, steps string) *Ant {
	return NewAnt(dimensions, StepsFromString(steps)...)
}

func NewAnt(dimensions Dimensions, steps ...Step) *Ant {

	Steps(steps).Numerate()

	cells := make([]Cell, dimensions.size, dimensions.size)
	cell := &cells[dimensions.IndexOf(dimensions.Center())]
	cell.Point = dimensions.Center()
	cell.Step = steps[0]

	return &Ant{
		Cells:      cells,
		Position:   cell,
		steps:      steps,
		Locker:     &sync.Mutex{},
		Dimensions: dimensions,
	}
}

func (ant *Ant) String() string {
	minX := ant.Dimensions.BottomLeft.X
	minY := ant.Dimensions.BottomLeft.Y
	maxX := ant.Dimensions.TopRight.X
	maxY := ant.Dimensions.TopRight.Y

	builder := strings.Builder{}
	builder.Grow(int((maxX - minX) * (maxY - minY)))
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			cell := ant.Cells[ant.Dimensions.IndexOf(Point{
				X: x,
				Y: y,
			})]
			if cell.Step.Action != ActionNone {
				builder.WriteRune(rune(cell.Step.Action))
			} else {
				switch {
				case y == 0:
					builder.WriteRune('―')
				case x == 0:
					builder.WriteRune('|')
				default:
					builder.WriteRune('-')
				}
			}
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
