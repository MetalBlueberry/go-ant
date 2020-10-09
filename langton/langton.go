package langton

import (
	"errors"
	"fmt"
	"strings"
)

type Ant struct {
	Cells     []Cell
	Position  *Cell
	Direction Direction

	Dimensions Dimensions

	steps      []Step
	totalSteps int64
	stuck      bool
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

type Cell struct {
	Point
	Step Step
}

type Steps []Step

type Step struct {
	Index  int
	Action Action

	nextIndex     int
	previousIndex int
}

type Action rune

const (
	ActionNone      Action = 0
	ActionTurnLeft         = 'L'
	ActionTurnRight        = 'R'
	ActionStraight         = 'S'
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

func (ant *Ant) TotalSteps() int64 {
	return ant.totalSteps
}

func (ant *Ant) Stuck() bool {
	return ant.stuck
}

func (ant *Ant) Next() (*Cell, error) {
	if ant.stuck {
		return nil, errors.New("Ant is stuck, grow the grid before calling Next")
	}

	ant.Direction = ant.Direction.Turn(ant.Position.Step.Action)

	ant.Position.UpdateNextStep(ant.steps)

	nextPoint := ant.Position.Point.Walk(ant.Direction)

	nextPosition, err := ant.EnsureCellAt(nextPoint)
	if err != nil {
		ant.stuck = true

		ant.Position.UpdatePreviousStep(ant.steps)
		ant.Direction = ant.Direction.Unturn(ant.Position.Step.Action)

		return ant.Position, err
	}
	ant.Position = nextPosition

	ant.totalSteps++
	return ant.Position, nil
}

func (ant *Ant) NextN(steps int) (cell *Cell, err error) {
	if steps < 0 {
		panic("steps must be >= 0")
	}
	if steps == 0 {
		return ant.Position, nil
	}
	for i := 0; i < steps; i++ {
		cell, err = ant.Next()
		if err != nil {
			return cell, err
		}
	}
	return cell, err
}

func (cell *Cell) UpdateNextStep(steps []Step) {
	cell.Step = steps[cell.Step.nextIndex]
}

func (cell *Cell) UpdatePreviousStep(steps []Step) {
	cell.Step = steps[cell.Step.previousIndex]
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

func (ant *Ant) CellAt(position Point) (*Cell, error) {
	if !ant.Dimensions.isPointInside(position) {
		return nil, errors.New("Next step is out of bounds")
	}
	posIndex := ant.Dimensions.IndexOf(position)
	cell := &ant.Cells[posIndex]
	if cell.Step.Action == ActionNone {
		return nil, errors.New("Cell not initialized")
	}
	return cell, nil
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
		steps[i].previousIndex = i - 1
		steps[i].Index = i
		steps[i].nextIndex = i + 1
	}
	steps[len(steps)-1].nextIndex = 0
	steps[0].previousIndex = len(steps) - 1
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

func (dim Dimensions) String() string {
	return fmt.Sprintf("%dx%d", dim.width, dim.height)
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
		case rune(ActionStraight):
			out[i] = Step{
				Action: ActionStraight,
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

	cells := make([]Cell, dimensions.Size, dimensions.Size)
	cell := &cells[dimensions.IndexOf(dimensions.Center())]
	cell.Point = dimensions.Center()
	cell.Step = steps[0]

	return &Ant{
		Cells:      cells,
		Position:   cell,
		steps:      steps,
		Dimensions: dimensions,
	}
}

func (ant *Ant) Grow(dimensions Dimensions) error {
	if ant.Dimensions.height >= dimensions.height || ant.Dimensions.width >= dimensions.width {
		return errors.New("New dimensions are equal or smaller than the current dimensions")
	}

	newCells := make([]Cell, dimensions.Size, dimensions.Size)
	for i := range ant.Cells {
		old := ant.Cells[i]
		if old.Step.Action == ActionNone {
			continue
		}
		newCells[dimensions.IndexOf(old.Point)] = old
	}
	ant.Cells = newCells
	ant.Dimensions = dimensions
	ant.Position = &newCells[dimensions.IndexOf(ant.Position.Point)]
	ant.stuck = false
	return nil
}

func (ant *Ant) StringMargin(margin int64) string {
	minX := ant.Dimensions.BottomLeft.X - margin
	minY := ant.Dimensions.BottomLeft.Y - margin
	maxX := ant.Dimensions.TopRight.X + margin
	maxY := ant.Dimensions.TopRight.Y + margin

	builder := strings.Builder{}
	builder.Grow(int((maxX - minX) * (maxY - minY)))
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			p := Point{
				X: x,
				Y: y,
			}
			if ant.Dimensions.isPointInside(p) {
				cell := ant.Cells[ant.Dimensions.IndexOf(p)]
				if cell.Step.Action != ActionNone {
					builder.WriteRune(rune(cell.Step.Action))
					continue
				}
			}

			switch {
			case y == 0:
				builder.WriteRune('―')
			case x == 0:
				builder.WriteRune('|')
			default:
				builder.WriteRune('-')
			}
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (ant *Ant) String() string {
	return ant.StringMargin(0)
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
