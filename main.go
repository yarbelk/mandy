package main

import (
	"github.com/yarbelk/mandy/lib"
	"fmt"
	"math"
)

var COLORS [16]string = [...]string{
  "\033[31m",
  "\033[32m",
  "\033[33m",
  "\033[34m",
  "\033[35m",
  "\033[36m",
  "\033[37m",
  "\033[1;30m",
  "\033[1;31m",
  "\033[1;32m",
  "\033[1;33m",
  "\033[1;34m",
  "\033[1;35m",
  "\033[1;36m",
  "\033[1;37m",
  "\033[1;38m",
}


func getCharForVal(value int16) string {
	var i int = int(math.Floor(math.Log2(float64(value))))
	if i > 15 {
		i = 15
	} else if i < 0 {
		i = 0
	}
	return fmt.Sprintf("%s*\033[0m", COLORS[i])
}

func updateWindow(inputChan chan mandy.PixelValue, doneChan chan bool) {
	var lasty int = 0
	for input := range inputChan {
		if lasty != input.Y {
			fmt.Println("")
		}
		fmt.Printf("%s", getCharForVal(input.Value))
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

	go mandy.Mandy(inputChan, outputChan, 32767, 2.0)

	doneChan := make(chan bool)
	go updateWindow(outputChan, doneChan)
	for _ = range doneChan {
	}
}
