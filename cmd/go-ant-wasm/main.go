// Copyright [2019] [Mark Farnan]

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// NOTICE:  Much of this demo is a re-write of the 'Moving red Laser' demo by Martin Olsansky https://medium.freecodecamp.org/webassembly-with-golang-is-fun-b243c0e34f02
// It has been re-written to make use of the go-canvas library,  and avoid context calls for drawing.

package main

import (
	"go-ant/langton"
	"image/color"
	"log"
	"strings"
	"syscall/js"
	"time"

	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/markfarnan/go-canvas/canvas"
	"golang.org/x/image/colornames"
)

var done chan struct{}

var cvs *canvas.Canvas2d
var width float64
var height float64

var ant *langton.Ant
var colors color.Palette
var (
	sequence    = "LLLLRRRR"
	cellSize    = 20
	stepsXFrame = 2
)

func main() {

	cvs, _ = canvas.NewCanvas2d(true)

	height = float64(cvs.Height())
	width = float64(cvs.Width())

	cvs.Start(25, nil)
	fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		sequenceInput := getById("sequence")

		sequence = strings.TrimSpace(sequenceInput.Get("value").String())
		log.Print(sequence)
		ant = langton.NewAntFromString(langton.NewBoard(int64(cvs.Height()/(2*cellSize))), sequence)
		p, err := colorful.SoftPalette(len(sequence) + 1)
		if err != nil {
			panic(err)
		}
		colors = langton.ToPalette(p)

		cvs.Gc().SetFillColor(colornames.White)
		cvs.Gc().Clear()

		return nil
	})
	onClick("apply", fn)
	go doEvery(time.Millisecond*200, Render) // Kick off the Render function as go routine as it never returns
	<-done
}

// Helper function which calls the required func (in this case 'render') every time.Duration,  Call as a go-routine to prevent blocking, as this never returns
func doEvery(d time.Duration, f func(time.Time) error) {
	for x := range time.Tick(d) {
		err := f(x)
		if err != nil {
			log.Print("Error:, ", err)
			return
		}
	}
}

// Called from the 'requestAnnimationFrame' function.   It may also be called seperatly from a 'doEvery' function, if the user prefers drawing to be seperate from the annimationFrame callback
func Render(when time.Time) error {
	if ant == nil {
		return nil
	}
	gc := cvs.Gc()
	for i := 0; i < stepsXFrame; i++ {
		cell, err := ant.Next()
		if err != nil {
			return err
		}

		x := (float64(cell.X) + float64(cvs.Width())/float64(2.0*cellSize)) * float64(cellSize)
		y := (float64(cell.Y) + float64(cvs.Height())/float64(2.0*cellSize)) * float64(cellSize)

		gc.BeginPath()
		gc.SetFillColor(colors[cell.Step.Index])
		gc.SetStrokeColor(color.RGBA{0, 0, 0, 0})
		draw2dkit.Rectangle(gc, x, y, x+float64(cellSize), y+float64(cellSize))
		gc.FillStroke()
		gc.Close()
	}

	return nil
}

func onClick(id string, fn js.Func) {
	element := getById("apply")
	element.Set("onclick", fn)
}

func getById(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}
