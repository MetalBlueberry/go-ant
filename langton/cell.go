package langton

import "fmt"

// Cell represents a Position where the ant can walk and the Action it takes
type Cell struct {
	Point
	Step Step
}

// UpdateNextStep given the sequence of Steps, updates to the next one
func (cell *Cell) UpdateNextStep(steps []Step) {
	cell.Step = steps[cell.Step.nextIndex]
}

// UpdatePreviousStep the opposite of UpdateNextStep
func (cell *Cell) UpdatePreviousStep(steps []Step) {
	cell.Step = steps[cell.Step.previousIndex]
}

// String is the string representation of a Cell
func (cell *Cell) String() string {
	return fmt.Sprintf(
		"%s, %s",
		cell.Point,
		cell.Step,
	)
}
