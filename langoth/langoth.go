package langoth

import (
	"image/color"
	"strings"
)

type Ant struct {
	Cells     map[Point]*Cell
	Position  *Cell
	Steps     []Step
	Direction Direction
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
	Color  color.Color
	Action Action

	nextIndex int
}

type Action rune

const (
	ActionTurnLeft  Action = 'L'
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

func (ant *Ant) Next() {

	ant.Direction = ant.Direction.Turn(ant.Position.Step.Action)

	ant.Position.UpdateNextStep(ant.Steps)

	nextPoint := ant.Position.Point.Walk(ant.Direction)

	ant.Position = ant.EnsureCellAt(nextPoint)

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

func (ant *Ant) EnsureCellAt(position Point) *Cell {
	cell, exist := ant.Cells[position]
	if !exist {
		cell = &Cell{
			Point: position,
			Step:  ant.Steps[0],
		}
		ant.Cells[position] = cell
	}
	return cell
}

func (steps Steps) Numerate() {
	for i := 0; i < len(steps); i++ {
		steps[i].nextIndex = i + 1
	}
	steps[len(steps)-1].nextIndex = 0
}

func NewAnt(steps ...Step) *Ant {

	Steps(steps).Numerate()

	initialPoint := Point{
		X: 0,
		Y: 0,
	}

	cells := map[Point]*Cell{
		initialPoint: {
			Step:  steps[0],
			Point: initialPoint,
		},
	}

	return &Ant{
		Cells:    cells,
		Position: cells[initialPoint],
		Steps:    steps,
	}
}

func (ant *Ant) String() string {
	var minX, maxX, minY, maxY int64
	for cell := range ant.Cells {
		minX = min(cell.X, minX)
		minY = min(cell.Y, minY)

		maxX = max(cell.X, maxX)
		maxY = max(cell.Y, maxY)
	}

	minX--
	minY--
	maxX++
	maxY++

	builder := strings.Builder{}
	builder.Grow(int((maxX - minX) * (maxY - minY)))
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			cell, exist := ant.Cells[Point{
				X: x,
				Y: y,
			}]
			if exist {
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
