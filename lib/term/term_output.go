package term

import mandy "github.com/yarbelk/mandy/lib"

type WindowPoint struct {
	X, Y  int
	Point complex128
}

// ConverganceValue
type ConverganceValue struct {
	X, Y  int
	Value int16
}

// WindowLimits are the terminal window extents and steps. Fails at
// the 'useful zero value' right now
type WindowLimits struct {
	X, Y         int32
	xMin, xMax   float64
	yMin, yMax   float64
	xStep, yStep float64
}

// NewWindowLimits is an ugly piece of code.  given your terminal size x, y
// and your extent in the complext plane, xMin, xMax, yMin, yMax: what
// are your steps and extents.
// it should be memoized and behaviioural.  its not actually needed out side
// of the 'lets render this' question.
func NewWindowLimits(x, y int32, xMin, xMax, yMin, yMax float64) WindowLimits {
	var xStep float64 = (xMax - xMin) / float64(x)
	var yStep float64 = (yMax - yMin) / float64(y)
	return WindowLimits{x, y, xMin, xMax, yMin, yMax, xStep, yStep}
}

// ProdWindowPoints breaks a window into single character output cells.
// it, and all the 'window' stuff should be their own package.
func ProdWindowPoints(windowLimit *WindowLimits, output chan WindowPoint) {
	var re, im float64
	var x, y int32
	for y = 0; y < windowLimit.Y; y++ {
		for x = 0; x < windowLimit.X; x++ {
			re = windowLimit.xMin + float64(x)*windowLimit.xStep
			im = windowLimit.yMax - float64(y)*windowLimit.yStep
			output <- WindowPoint{X: int(x), Y: int(y), Point: complex(re, im)}
		}
	}
	close(output)
}

// Mandy returns the convergence of a point and radius from the input channel
// and by 'returns', i mean puts it into the output channel.  its designed
// for being run in multiple go routines
func Mandy(inputChan chan WindowPoint, outputChan chan ConverganceValue, limit int32, radius float64) {
	for input := range inputChan {
		outputChan <- ConverganceValue{
			X:     input.X,
			Y:     input.Y,
			Value: mandy.MandelbrotRecursionLimit(input.Point, limit, radius),
		}
	}
}
