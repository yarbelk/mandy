package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"sync"

	mandy "github.com/yarbelk/mandy/lib"
	"golang.org/x/crypto/ssh/terminal"
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

// This collects the results and prints it in one go; so it can take advantage of coprocessing and
// multiple workers.  You don't get nice 'fill up the terminal line by line' any more though'
func updateWindowCollector(limits mandy.WindowLimits, inputChan chan mandy.PixelValue, doneChan chan struct{}) {
	cellCount := limits.X * limits.Y
	cells := make([][]int16, limits.Y)
	for i := 0; int32(i) < limits.Y; i++ {
		cells[i] = make([]int16, limits.X)
	}
	for cell := range inputChan {
		cells[cell.Y][cell.X] = cell.Value
		cellCount--
		if cellCount == 0 {
			break
		}
	}
	for _, row := range cells {
		for _, cell := range row {
			fmt.Printf("%s", getCharForVal(cell))
		}
		fmt.Printf("\n")
	}
	doneChan <- struct{}{}
	close(doneChan)
}

func main() {
	flag.Parse()

	var inputChan chan mandy.WindowPoint = make(chan mandy.WindowPoint)
	var outputChan chan mandy.PixelValue = make(chan mandy.PixelValue)
	doneChan := make(chan struct{})

	// boundry checking for later, because i truncate it to an int32
	if *depth > math.MaxInt32 {
		*depth = math.MaxInt32
	}

	var windowLimits mandy.WindowLimits
	screenX, screenY, err := terminal.GetSize(0)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// make it nicer in size so your promts and borders don't mess with the output
	screenX -= 10
	screenY -= 10

	windowLimits = mandy.NewWindowLimits(
		int32(screenX), int32(screenY),
		*xMin, *xMax,
		*yMin, *yMax,
	)

	go mandy.ProdWindowPoints(&windowLimits, inputChan)

	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			mandy.Mandy(inputChan, outputChan, int32(*depth), *radius)
		}()
	}
	go updateWindowCollector(windowLimits, outputChan, doneChan)

	<-doneChan
	wg.Wait()
	close(outputChan)
}
