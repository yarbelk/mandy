package mandy

import (
	"math"
	"math/cmplx"
	"testing"
)

func complexEquality(x, y complex128) bool {
	return (math.Abs(cmplx.Abs(x-y)) < 0.0000001)
}

func TestMandelbrotFunction(t *testing.T) {
	var input complex128 = 1.0 + 1.0i
	var expected complex128 = 1.0 + 3.0i
	var result complex128
	result = Mandelbrot(input, input)

	if !complexEquality(expected, result) {
		t.Logf("Expected %v, was %v\n", expected, result)
	}
}

func TestConvergesOpenFunction(t *testing.T) {
	var radius float64 = 2.0
	var divergent complex128 = 2.0 + 0i
	var convergent_complex complex128 = 1.0 + 1.0i
	var convergent_near_border complex128 = 1.99999 + 0i

	if Converges(&divergent, radius) {
		t.Logf("Expected convergance to be open set\n")
		t.Fail()
	}

	if !Converges(&convergent_complex, radius) {
		t.Logf("Expected inside radius to converge: %v should be inside %f.", convergent_complex, radius)
		t.Fail()
	}

	if !Converges(&convergent_near_border, radius) {
		t.Logf("Expected inside radius to converge: %v should be inside %f.", convergent_near_border, radius)
		t.Fail()
	}
}

func TestGetRecursionLimitFunction(t *testing.T) {
	var two complex128 = 1.0 + 1.0i
	var hitsLimit complex128 = 0 + 0i
	var result int16

	result = MandelbrotRecursionLimit(two, 32767, 2.0)
	if result != 1 {
		t.Logf("Expected %v to have a limit of 1, was %d", two, result)
		t.Fail()
	}

	result = MandelbrotRecursionLimit(hitsLimit, 32767, 2.0)
	if result != 32767 {
		// normalized to 32767
		t.Logf("Expected %v to have a limit of 32767, was %d", hitsLimit, result)
		t.Fail()
	}
}
