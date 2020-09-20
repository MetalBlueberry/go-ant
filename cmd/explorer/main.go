package main

import (
	"fmt"
	"go-ant/langton"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	power := 12.0
	works := make(chan string)
	go func(works chan string) {
		for i := 1.0; i <= power; i++ {
			c := Combinations(i)
			for i := range c {
				works <- c[i]
			}
		}
		close(works)
	}(works)

	wg := sync.WaitGroup{}
	workers := 4
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for {
				s, open := <-works
				if !open {
					log.Printf("closed %d", id)
					break
				}
				log.Printf("Worker %d on %s\n", id, s)
				Calculate(s)
			}
		}(i)
	}
	wg.Wait()
	log.Print("Done")
}

func Combinations(power float64) []string {
	var combinations int64
	var i int64

	combinations = int64(math.Pow(2, power))

	out := make([]string, combinations)

	pad := strconv.Itoa(int(power))

	for i = 0; i < combinations; i++ {
		binary := fmt.Sprintf("%"+pad+"v", strconv.FormatInt(i, 2))
		replacer := strings.NewReplacer(" ", "L", "0", "L", "1", "R")
		out[i] = fmt.Sprint(replacer.Replace(binary))
	}
	return out
}

func Calculate(steps string) {

	ant := langton.NewAntFromString(
		langton.NewBoard(1000),
		steps,
	)
	_, err := ant.NextN(10000000)
	if err != nil {
		log.Printf("reached limit! %s\n", steps)
	}
	colorfulPalette, err := colorful.SoftPalette(len(steps))
	img := langton.ToImage(ant, langton.ToPalette(colorfulPalette), 1)
	file, err := os.Create("outs/" + steps + ".png")
	if err != nil {
		log.Print("Cannot create file!")
		panic(err)
	}
	err = png.Encode(file, img)
	if err != nil {
		log.Print("Cannot encode image")
		panic(err)
	}
	file.Close()
}
