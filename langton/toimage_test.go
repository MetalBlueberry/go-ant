package langton

import (
	"image/png"
	"os"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

func TestToImage(t *testing.T) {
	type args struct {
		ant     *Ant
		palette []colorful.Color
	}

	ant := NewAntFromString(NewBoard(100), "LR")
	for {
		_, err := ant.Next()
		if err != nil {
			break
		}
	}
	palette, err := colorful.SoftPalette(len(ant.steps))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Simple",
			args: args{
				ant:     ant,
				palette: palette,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToImage(tt.args.ant, tt.args.palette)
			file, err := os.Create("pic.png")
			if err != nil {
				panic(err)
			}
			png.Encode(file, got)
			file.Close()
		})
	}
}
