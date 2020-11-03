package langton

import "fmt"

type Steps []Step

type Step struct {
	Index  int
	Action Action

	nextIndex     int
	previousIndex int
}

func (step Step) String() string {
	return fmt.Sprintf(
		"%d: %s",
		step.Index,
		step.Action,
	)
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

func (steps Steps) Numerate() {
	for i := 0; i < len(steps); i++ {
		steps[i].previousIndex = i - 1
		steps[i].Index = i
		steps[i].nextIndex = i + 1
	}
	steps[len(steps)-1].nextIndex = 0
	steps[0].previousIndex = len(steps) - 1
}
