package mandy

import (
	"math"
	"math/cmplx"
)

// Converges returns true if the mandlebrot value is less
// than the radius.
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
