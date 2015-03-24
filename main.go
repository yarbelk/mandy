package main

import (
	"github.com/yarbelk/mandy/mandy"
	"code.google.com/p/goncurses"
	"log"
)

func getCharForVal(value int16) goncurses.Char {
	if val >= 500 {
		return '*' || goncurses.C_RED
	} else {
		return '-' || goncurses.C_BLUE
	}
}

func updateWindow(inputChan chan mandy.PixelValue, window *goncurses.Window, doneChan chan bool) {
	var ach goncurses.Char
	var pair int16
	for input := range inputChan {
		stdscr.MoveAddChar(input.y, input.x, getCharForVal(input.value))
	}

	doneChan <- true
	close(doneChan)
}

func main() {
	var inputChan chan WindowPoint = make(chan WindowPoint)
	var outpuChan chan PixelValue = make(chan PixelValue)

	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal("init:", err)
	}
	defer goncurses.End()

	var windowLimits mandy.WindowLimits
	screenY, screenX := stdscr.MaxYX()

	windowLimits = mandy.NewWindowLimits(
		screenX, screenY,
		-2.0, 0.5,
		-1.0, 1.0,
	)

	go ProdWindowPoints(&windowLimits, inputChan)
	close(inputChan)

	go Mandy(inputChan, outpuChan, 500, 2.0)
	close(inputChan)

	doneChan := make(chan bool)
	go updateWindow(outpuChan, stdscr, doneChan)
	for _ := range doneChan {
	}
	stdscr.Refresh()
}
