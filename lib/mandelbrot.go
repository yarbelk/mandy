package mandy

import (
	"math"
	"math/cmplx"
)

type WindowPoint struct {
	X, Y  int
	Point complex128
}

type PixelValue struct {
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

func Converges(mandyVal *complex128, radius float64) bool {
	return (cmplx.Abs(*mandyVal) < radius)
}

// Mandelbrot is the actual mandlebrot calculation
func Mandelbrot(z, c complex128) complex128 {
	return cmplx.Pow(z, 2.0+0i) + c
}

// MandelbrotRecursionLimit returns what value the mandelbrot function
// ceases to be convergent at in a radius.  It normalizes back to an
// int16.  some benchmarks show that using the complex128 is
// not as performant as other approaches, but its easier to read and i like it.
// Performance isn't a huge driver for this.
func MandelbrotRecursionLimit(input complex128, limit int32, radius float64) int16 {
	var i int32
	var z, c complex128
	c = input
	for i = 0; i < limit; i++ {
		z = Mandelbrot(z, c)
		if !Converges(&z, radius) {
			break
		}
	}
	f := float64(i) / float64(limit)
	f *= math.MaxInt16
	return int16(f)
}

// Mandy returns the convergence of a point and radius from the input channel
// and by 'returns', i mean puts it into the output channel.  its designed
// for being run in multiple go routines
func Mandy(inputChan chan WindowPoint, outputChan chan PixelValue, limit int32, radius float64) {
	for input := range inputChan {
		outputChan <- PixelValue{
			X:     input.X,
			Y:     input.Y,
			Value: MandelbrotRecursionLimit(input.Point, limit, radius),
		}
	}
}
