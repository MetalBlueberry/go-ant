package langton

import (
	"errors"
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
