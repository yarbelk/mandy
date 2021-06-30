package main

import (
	"flag"
	"fmt"
	"math"

	mandy "github.com/yarbelk/mandy/lib"
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

var (
	xMin   *float64 = flag.Float64("xmin", -2.5, "Minimum on real axis")
	xMax   *float64 = flag.Float64("xmax", 0.5, "Max on real axis")
	yMin   *float64 = flag.Float64("ymin", -1.0, "Minimum on imaginary axis")
	yMax   *float64 = flag.Float64("ymax", 1.0, "max on imaginary axis")
	depth  *int64   = flag.Int64("depth", 32767, "max recursion depth, int32")
	radius *float64 = flag.Float64("radius", 2.0, "radius to test for convergence")
)

func getCharForVal(value int16) string {
	var i int = int(math.Floor(math.Log2(float64(value))))
	switch {
	case i < 0:
		i = 0
	case i > 15:
		i = 15
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
	flag.Parse()

	// boundry checking for later
	if *depth > math.MaxInt32 {
		*depth = math.MaxInt32
	}

	var windowLimits mandy.WindowLimits
	screenY, screenX := 60, 160

	windowLimits = mandy.NewWindowLimits(
		int32(screenX), int32(screenY),
		*xMin, *xMax,
		*yMin, *yMax,
	)

	go mandy.ProdWindowPoints(&windowLimits, inputChan)

	go mandy.Mandy(inputChan, outputChan, int32(*depth), *radius)
	doneChan := make(chan bool)
	go updateWindow(outputChan, doneChan)

	<-doneChan
}
