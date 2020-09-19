package langton

import (
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

func TestToImage(t *testing.T) {
	type args struct {
		ant     *Ant
		palette color.Palette
	}

	ant := NewAntFromString(NewBoard(10), "LR")
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
				palette: ToPalette(palette),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToImage(tt.args.ant, tt.args.palette, 9)
			file, err := os.Create("pic.png")
			if err != nil {
				panic(err)
			}
			png.Encode(file, got)
			file.Close()
		})
	}
}
