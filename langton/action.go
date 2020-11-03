package langton

type Action rune

func (action Action) String() string {
	switch action {
	case ActionTurnLeft:
		return "Left"
	case ActionTurnRight:
		return "Right"
	case ActionStraight:
		return "Straight"
	default:
		return "Unknown"
	}
}

const (
	ActionNone      Action = 0
	ActionTurnLeft         = 'L'
	ActionTurnRight        = 'R'
	ActionStraight         = 'S'
)
