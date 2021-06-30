package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"

	"github.com/yarbelk/mandy/lib/term"
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
	xMin    *float64 = flag.Float64("xmin", -2.5, "Minimum on real axis")
	xMax    *float64 = flag.Float64("xmax", 0.5, "Max on real axis")
	yMin    *float64 = flag.Float64("ymin", -1.0, "Minimum on imaginary axis")
	yMax    *float64 = flag.Float64("ymax", 1.0, "max on imaginary axis")
	depth   *int64   = flag.Int64("depth", 32767, "max recursion depth, int32")
	radius  *float64 = flag.Float64("radius", 2.0, "radius to test for convergence")
	screenX *int64   = flag.Int64("w", 160, "width in characters. (cowsay defaults to 40 fyi, and -n 'fixes' formating)")
	screenY *int64   = flag.Int64("h", 60, "height in characters")
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
func updateWindowCollector(limits term.WindowLimits, inputChan chan term.ConverganceValue, doneChan chan struct{}) {
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

	var inputChan chan term.WindowPoint = make(chan term.WindowPoint)
	var outputChan chan term.ConverganceValue = make(chan term.ConverganceValue)
	doneChan := make(chan struct{})

	// boundry checking for later, because i truncate it to an int32
	if *depth > math.MaxInt32 {
		*depth = math.MaxInt32
	}

	windowLimits := term.NewWindowLimits(
		int32(*screenX), int32(*screenY),
		*xMin, *xMax,
		*yMin, *yMax,
	)

	go term.ProdWindowPoints(&windowLimits, inputChan)

	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			term.Mandy(inputChan, outputChan, int32(*depth), *radius)
		}()
	}
	go updateWindowCollector(windowLimits, outputChan, doneChan)

	<-doneChan
	wg.Wait()
	close(outputChan)
}
