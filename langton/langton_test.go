package langton

import (
	"reflect"
	"strings"
	"testing"
)

func Test0(t *testing.T) {
	Iterations(t, 0, `
-|-
―L―
-|-`)
}
func Test1(t *testing.T) {
	Iterations(t, 1, `
--|-
―LR―
--|-`)
}
func Test2(t *testing.T) {
	Iterations(t, 2, `
--|-
―RR―
-L|-
--|-`)
}

func Test5(t *testing.T) {
	Iterations(t, 5, `
--|--
―RLL―
-RR--
--|--`)
}

func Iterations(t *testing.T, n int, expected string) {
	steps := []Step{
		{
			Action: ActionTurnLeft,
		},
		{
			Action: ActionTurnRight,
		},
	}
	ant := NewAnt(
		steps...,
	)
	for i := 0; i < n; i++ {
		ant.Next()
	}
	if strings.TrimSpace(expected) != strings.TrimSpace(ant.String()) {
		t.Errorf("expected \n%s, obtained\n%s", expected, ant)
	}

}

func TestDirection_Turn(t *testing.T) {
	type args struct {
		action Action
	}
	tests := []struct {
		name string
		d    Direction
		args args
		want Direction
	}{
		{
			name: "TurnRight",
			d:    DirectionTop,
			args: args{
				action: ActionTurnRight,
			},
			want: DirectionRight,
		},
		{
			name: "TurnLeft",
			d:    DirectionTop,
			args: args{
				action: ActionTurnLeft,
			},
			want: DirectionLeft,
		},
		{
			name: "FromLeft TurnRight",
			d:    DirectionLeft,
			args: args{
				action: ActionTurnRight,
			},
			want: DirectionTop,
		},
		{
			name: "FromLeft TurnLeft",
			d:    DirectionLeft,
			args: args{
				action: ActionTurnLeft,
			},
			want: DirectionDown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Turn(tt.args.action); got != tt.want {
				t.Errorf("Direction.Turn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCell_UpdateNextStep(t *testing.T) {
	type fields struct {
		Point Point
		Step  Step
	}
	type args struct {
		steps []Step
	}
	steps := Steps{
		Step{
			Action: ActionTurnLeft,
		},
		Step{
			Action: ActionTurnRight,
		},
	}
	steps.Numerate()

	tests := []struct {
		name   string
		fields fields
		args   args
		expect int
	}{
		{
			name: "0 to 1",
			args: args{
				steps: steps,
			},
			fields: fields{
				Step: steps[0],
			},
			expect: 1,
		},
		{
			name: "1 to 0",
			args: args{
				steps: steps,
			},
			fields: fields{
				Step: steps[1],
			},
			expect: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := &Cell{
				Point: tt.fields.Point,
				Step:  tt.fields.Step,
			}
			cell.UpdateNextStep(tt.args.steps)
			if steps[tt.expect] != cell.Step {
				t.Errorf("cell.UpdateNextStep = %v, want %v", cell.Step, steps[tt.expect])
			}
		})
	}
}

func TestStepsFromString(t *testing.T) {
	type args struct {
		steps string
	}
	tests := []struct {
		name string
		args args
		want Steps
	}{
		{
			name: "Left Right",
			args: args{
				steps: "LR",
			},
			want: Steps{
				Step{
					Action: ActionTurnLeft,
				},
				Step{
					Action: ActionTurnRight,
				},
			},
		},
		{
			name: "More steps",
			args: args{
				steps: "LRRLLRRR",
			},
			want: Steps{
				Step{
					Action: ActionTurnLeft,
				},
				Step{
					Action: ActionTurnRight,
				},
				Step{
					Action: ActionTurnRight,
				},
				Step{
					Action: ActionTurnLeft,
				},
				Step{
					Action: ActionTurnLeft,
				},
				Step{
					Action: ActionTurnRight,
				},
				Step{
					Action: ActionTurnRight,
				},
				Step{
					Action: ActionTurnRight,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StepsFromString(tt.args.steps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StepsFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
