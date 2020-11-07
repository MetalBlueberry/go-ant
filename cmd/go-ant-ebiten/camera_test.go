package main

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/math/f64"
)

func TestMapVector(t *testing.T) {
	type args struct {
		origin f64.Vec2
		geoM   ebiten.GeoM
	}
	tests := []struct {
		name       string
		args       args
		wantOrigin f64.Vec2
		wantX      f64.Vec2
		wantY      f64.Vec2
	}{
		{
			name: "simple",
			args: args{
				origin: f64.Vec2{0, 0},
				geoM:   ebiten.GeoM{},
			},
			wantOrigin: f64.Vec2{0, 0},
			wantX:      f64.Vec2{1, 0},
			wantY:      f64.Vec2{0, 1},
		},
		{
			name: "translated",
			args: args{
				origin: f64.Vec2{0, 0},
				geoM:   ebiten.TranslateGeo(50, 50),
			},
			wantOrigin: f64.Vec2{50, 50},
			wantX:      f64.Vec2{1, 0},
			wantY:      f64.Vec2{0, 1},
		},
		{
			name: "Scaled up",
			args: args{
				origin: f64.Vec2{0, 0},
				geoM:   ebiten.ScaleGeo(2, 2),
			},
			wantOrigin: f64.Vec2{0, 0},
			wantX:      f64.Vec2{2, 0},
			wantY:      f64.Vec2{0, 2},
		},
		{
			name: "Scaled down",
			args: args{
				origin: f64.Vec2{0, 0},
				geoM:   ebiten.ScaleGeo(0.5, 0.5),
			},
			wantOrigin: f64.Vec2{0, 0},
			wantX:      f64.Vec2{0.5, 0},
			wantY:      f64.Vec2{0, 0.5},
		},
		{
			name: "Rotated",
			args: args{
				origin: f64.Vec2{0, 0},
				geoM:   ebiten.RotateGeo(math.Pi),
			},
			wantOrigin: f64.Vec2{0, 0},
			wantX:      f64.Vec2{-1, 0},
			wantY:      f64.Vec2{0, -1},
		},
		{
			name: "Trans + scale",
			args: args{
				origin: f64.Vec2{0, 0},
				geoM: func() ebiten.GeoM {
					geoM := ebiten.GeoM{}
					geoM.Translate(50, 50)
					geoM.Scale(2, 2)
					return geoM
				}(),
			},
			wantOrigin: f64.Vec2{100, 100},
			wantX:      f64.Vec2{2, 0},
			wantY:      f64.Vec2{0, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOrigin, gotX, gotY := MapVector(tt.args.origin, tt.args.geoM)
			if !vecEquals(gotX, tt.wantX) {
				t.Errorf("MapVector() gotX = %v, want %v", gotX, tt.wantX)
			}
			if !vecEquals(gotY, tt.wantY) {
				t.Errorf("MapVector() gotY = %v, want %v", gotY, tt.wantY)
			}
			if !vecEquals(gotOrigin, tt.wantOrigin) {
				t.Errorf("MapVector() gotOrigin = %v, want %v", gotOrigin, tt.wantOrigin)
			}
		})
	}
}

func vecEquals(a, b f64.Vec2) bool {
	tolerance := 0.001
	x := a[0] - b[0]
	y := a[1] - b[1]
	d := math.Sqrt(x*x + y*y)
	return d < tolerance
}
