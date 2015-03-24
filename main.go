package main

import (
	"github.com/yarbelk/mandy/lib"
	"code.google.com/p/goncurses"
	"log"
)

func getCharForVal(value int16) goncurses.Char {
	if value >= 500 {
		return '*'
	} else {
		return '-'
	}
}

func updateWindow(inputChan chan mandy.PixelValue, window *goncurses.Window, doneChan chan bool) {
	for input := range inputChan {
		window.MoveAddChar(input.Y, input.X, getCharForVal(input.Value))
	}

	doneChan <- true
	close(doneChan)
}

func main() {
	var inputChan chan mandy.WindowPoint = make(chan mandy.WindowPoint)
	var outputChan chan mandy.PixelValue = make(chan mandy.PixelValue)

	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal("init:", err)
	}
	defer goncurses.End()

	var windowLimits mandy.WindowLimits
	screenY, screenX := stdscr.MaxYX()

	windowLimits = mandy.NewWindowLimits(
		int32(screenX), int32(screenY),
		-2.0, 0.5,
		-1.0, 1.0,
	)

	go mandy.ProdWindowPoints(&windowLimits, inputChan)

	go mandy.Mandy(inputChan, outputChan, 16, 2.0)

	doneChan := make(chan bool)
	go updateWindow(outputChan, stdscr, doneChan)
	for _ = range doneChan {
	}
	close(outputChan)
	stdscr.Refresh()
}
