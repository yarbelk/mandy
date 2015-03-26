package main

import (
	"github.com/yarbelk/mandy/lib"
	"code.google.com/p/goncurses"
	"fmt"
)

func getCharForVal(value int16) goncurses.Char {
	if value >= 8 {
		return '*'
	} else {
		return '-'
	}
}

func updateWindow(inputChan chan mandy.PixelValue, doneChan chan bool) {
	var lasty int = 0
	for input := range inputChan {
		if lasty != input.Y {
			fmt.Println("")
		}
		fmt.Printf("%c", getCharForVal(input.Value))
		lasty = input.Y
	}
	fmt.Println("")

	doneChan <- true
	close(doneChan)
}

func main() {
	var inputChan chan mandy.WindowPoint = make(chan mandy.WindowPoint)
	var outputChan chan mandy.PixelValue = make(chan mandy.PixelValue)

	var windowLimits mandy.WindowLimits
	screenY, screenX := 60, 160

	windowLimits = mandy.NewWindowLimits(
		int32(screenX), int32(screenY),
		-2.0, 0.5,
		-1.0, 1.0,
	)

	go mandy.ProdWindowPoints(&windowLimits, inputChan)

	go mandy.Mandy(inputChan, outputChan, 16, 2.0)

	doneChan := make(chan bool)
	go updateWindow(outputChan, doneChan)
	for _ = range doneChan {
	}
}
