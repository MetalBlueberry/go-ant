package langton

import "fmt"

type Cell struct {
	Point
	Step Step
}

func (cell *Cell) UpdateNextStep(steps []Step) {
	cell.Step = steps[cell.Step.nextIndex]
}

func (cell *Cell) UpdatePreviousStep(steps []Step) {
	cell.Step = steps[cell.Step.previousIndex]
}

func (cell *Cell) String() string {
	return fmt.Sprintf(
		"%s, %s",
		cell.Point,
		cell.Step,
	)
}
