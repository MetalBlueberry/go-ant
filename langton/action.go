package langton

// Action represents what the ant does on a given Cell
type Action rune

// String returns an Action string representation
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
	// ActionNone does nothing
	ActionNone Action = 0
	// ActionTurnLeft turns left
	ActionTurnLeft = 'L'
	// ActionTurnRight turns right
	ActionTurnRight = 'R'
	// ActionStraight does not change direction
	ActionStraight = 'S'
)
