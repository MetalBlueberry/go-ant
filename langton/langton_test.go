package langton

import (
	"reflect"
	"strings"
	"testing"
)

var (
	StepsSimple  Steps = StepsFromString("LR")
	StepsAwesome Steps = StepsFromString("RLLLLRRRLLL")
)

func Test0(t *testing.T) {
	Iterations(t, 0, `
--|--
--|--
――L――
--|--
--|--`)
}
func Test1(t *testing.T) {
	Iterations(t, 1, `
--|--
--|--
―LR――
--|--
--|--`)
}
func Test2(t *testing.T) {
	Iterations(t, 2, `
--|--
--|--
―RR――
-L|--
--|--`)
}

func Test5(t *testing.T) {
	Iterations(t, 5, `
--|--
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
		NewDimensions(-2, -2, 2, 2),
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

func TestDirection_Unturn(t *testing.T) {
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
			name: "UnturnRight",
			d:    DirectionTop,
			args: args{
				action: ActionTurnRight,
			},
			want: DirectionLeft,
		},
		{
			name: "TurnLeft",
			d:    DirectionTop,
			args: args{
				action: ActionTurnLeft,
			},
			want: DirectionRight,
		},
		{
			name: "FromLeft TurnRight",
			d:    DirectionLeft,
			args: args{
				action: ActionTurnRight,
			},
			want: DirectionDown,
		},
		{
			name: "FromLeft TurnLeft",
			d:    DirectionLeft,
			args: args{
				action: ActionTurnLeft,
			},
			want: DirectionTop,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Unturn(tt.args.action); got != tt.want {
				t.Errorf("Direction.Unturn() = %v, want %v", got, tt.want)
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

func TestCell_UpdatePreviousStep(t *testing.T) {
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
			Action: ActionStraight,
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
			expect: 2,
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
			cell.UpdatePreviousStep(tt.args.steps)
			if steps[tt.expect] != cell.Step {
				t.Errorf("cell.UpdatePreviousStep = %v, want %v", cell.Step, steps[tt.expect])
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

func BenchmarkNext(b *testing.B) {
	ant := NewAnt(
		NewDimensions(-10000, -10000, 10000, 10000),
		StepsAwesome...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ant.Next()
	}
}

func TestDimensions_IndexOf(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name   string
		fields Dimensions
		args   args
		want   int
	}{
		{
			name: "initial corner",
			args: args{
				p: Point{-1, -1},
			},
			fields: NewDimensions(-1, -1, 1, 1),
			want:   0,
		},
		{
			name: "last corner",
			args: args{
				p: Point{1, 1},
			},
			fields: NewDimensions(-1, -1, 1, 1),
			want:   8,
		},
		{
			name: "step up",
			args: args{
				p: Point{-1, 0},
			},
			fields: NewDimensions(-1, -1, 1, 1),
			want:   3,
		},
		{
			name: "step right",
			args: args{
				p: Point{0, -1},
			},
			fields: NewDimensions(-1, -1, 1, 1),
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dim := &Dimensions{
				TopRight:   tt.fields.TopRight,
				BottomLeft: tt.fields.BottomLeft,
				width:      tt.fields.width,
				height:     tt.fields.height,
				Size:       tt.fields.Size,
			}
			if got := dim.IndexOf(tt.args.p); got != tt.want {
				t.Errorf("Dimensions.IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnt_Grow(t *testing.T) {
	type args struct {
		antSteps    int
		initialSize int64
		newSize     int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "grow grid by 1",
			args: args{
				antSteps:    0,
				initialSize: 2,
				newSize:     3,
			},
			wantErr: false,
		},
		{
			name: "grow grid by 10 with iterations",
			args: args{
				antSteps:    30,
				initialSize: 5,
				newSize:     15,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ant := NewAntFromString(
				NewBoard(tt.args.initialSize),
				"LR",
			)

			ant.NextN(tt.args.antSteps)

			expected := ant.StringMargin(tt.args.newSize - tt.args.initialSize)
			if err := ant.Grow(NewBoard(tt.args.newSize)); (err != nil) != tt.wantErr {
				t.Errorf("Ant.Grow() error = %v, wantErr %v", err, tt.wantErr)
			}
			obtained := ant.String()
			if expected != obtained {
				t.Errorf("Ant.Grow() doesn't keep shape\n%v\nwant\n%v", obtained, expected)
			}
		})
	}
}

func TestAnt_Grow_And_Walk(t *testing.T) {
	growAnt := NewAntFromString(
		NewBoard(1),
		"LR",
	)
	growAnt.NextN(1000)
	growAnt.Grow(NewBoard(2))
	growAnt.NextN(1000)

	staticAnt := NewAntFromString(
		NewBoard(2),
		"LR",
	)
	staticAnt.NextN(1000)

	static := staticAnt.String()
	grow := growAnt.String()
	if static != grow {
		t.Errorf("Ant.Grow() doesn't keep shape\n%v\nwant\n%v", grow, static)
	}
}
