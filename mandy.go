package mandy

import (
	"math/cmplx"
)

type WindowPoint struct {
	x, y int
	point complex128
}

type PixelValue struct {
	x, y int
	value int32
}

type WindowLimits struct {
	x, y int32
	xMin, xMax float64
	yMin, yMax float64
	xStep, yStep float64
}

type mandelState struct {
	radius, xStep, yStep, xMin, yMin, xMax, yMax float64
	limit int32
}

func NewWindowLimits(x, y int32, xMin, xMax, yMin, yMax float64) WindowLimits {
	var xStep float64 = (xMax - xMin) / float64(x)
	var yStep float64 = (yMax - yMin) / float64(y)
	return WindowLimits{x, y, xMin, xMax, yMin, yMax, xStep, yStep}
}

func ProdWindowPoints(windowLimit *WindowLimits, output chan WindowPoint) {
	var re, im float64
	var x, y int32
	for y=0; y<windowLimit.y; y++ {
		for x=0; x<windowLimit.x; x++ {
			re = windowLimit.xMin + float64(x) * windowLimit.xStep
			im = windowLimit.yMax - float64(y) * windowLimit.yStep
			output <- WindowPoint{x: int(x), y: int(y), point: complex(re, im)}
		}
	}
}

func Converges(mandyVal *complex128, radius float64) bool {
	return (cmplx.Abs(*mandyVal) < radius)
}

func Mandelbrot(z, c complex128) complex128 {
	return cmplx.Pow(z, 2.0+0i) + c
}

func MandelbrotRecursionLimit(input complex128, limit int32, radius float64) int32 {
	var i int32
	var z, c complex128
	c = input
	for i=0; i<limit; i++ {
		z = Mandelbrot(z,c)
		if !Converges(&z, radius) {
			break
		}
	}
	return i
}

func Mandy(inputChan chan WindowPoint, outputChan chan PixelValue, limit int32, radius float64) {
	for input := range inputChan {
		outputChan<- PixelValue{
			x: input.x,
			y: input.y,
			value: MandelbrotRecursionLimit(input.point, limit, radius),
		}
	}
}
