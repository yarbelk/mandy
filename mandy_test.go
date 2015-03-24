package mandy

import (
	"testing"
	"math"
	"math/cmplx"
)

func complexEquality(x, y complex128) bool {
	return (math.Abs(cmplx.Abs(x - y)) < 0.0000001)
}


func TestMandelbrotFunction(t *testing.T) {
	var input complex128 = 1.0 + 1.0i
	var expected complex128 = 1.0 + 3.0i
	var result complex128
	result = Mandelbrot(input, input)

	if !complexEquality(expected, result) {
		t.Fatalf("Expected %v, was %v\n", expected, result)
	}
}

func TestConvergesOpenFunction(t *testing.T) {
	var radius float64 = 2.0
	var divergent complex128 = 2.0 + 0i
	var convergent_complex complex128 = 1.0 + 1.0i
	var convergent_near_border complex128 = 1.99999 + 0i

	if Converges(&divergent, radius) {
		t.Fatalf("Expected convergance to be open set\n")
	}

	if !Converges(&convergent_complex, radius) {
		t.Fatalf("Expected inside radius to converge: %v should be inside %f.", convergent_complex, radius)
	}

	if !Converges(&convergent_near_border, radius) {
		t.Fatalf("Expected inside radius to converge: %v should be inside %f.", convergent_near_border, radius)
	}
}

func TestGetRecursionLimitFunction(t *testing.T) {
	var two complex128 = 1.0 + 1.0i
	var hits_limit complex128 = 0 + 0i
	var result int32

	result = MandelbrotRecursionLimit(two, 10, 2.0)
	if result != 1 {
		t.Fatalf("Expected %v to have a limit of 2, was %d", two, result)
	}

	result = MandelbrotRecursionLimit(hits_limit, 20, 2.0)
	if result != 20 {
		t.Fatalf("Expected %v to have a limit of 2, was %d", two, result)
	}
}

func TestGeneratesStepsFromLimits(t *testing.T) {
	var windowLimits = NewWindowLimits(2, 2, -1.0, 1.0, -1.0, 1.0)

	if math.Abs(windowLimits.xStep - 1.0) > 0.00001 {
		t.Fatalf("Expected xStep to be ~= 1.0, was %f", windowLimits.xStep)
	}

	if math.Abs(windowLimits.yStep - 1.0) > 0.00001 {
		t.Fatalf("Expected yStep to be ~= 1.0, was %f", windowLimits.yStep)
	}
}

/* Expect the chanel to spit out a set of complex numbers and
 * x/y coordinates at the right resolution and with y flipped
 */
func TestGenerateListOfPointPixelPairs(t *testing.T) {
	var windowLimits WindowLimits
	var testPoint, expectedPoint WindowPoint

	testChan := make(chan WindowPoint, 4)
	expected := [...]WindowPoint{
		WindowPoint{x:0, y:0, point: -1.0 + 1.0i},
		WindowPoint{x:1, y:0, point: 0.0 + 1.0i},
		WindowPoint{x:0, y:1, point: -1.0 + 0.0i},
		WindowPoint{x:1, y:1, point: 0.0 + 0.0i},
	}

	windowLimits = NewWindowLimits(2, 2, -1.0, 1.0, -1.0, 1.0)
	ProdWindowPoints(&windowLimits, testChan)
	for _, testPoint = range expected {

		expectedPoint = <-testChan
		if expectedPoint != testPoint {
			t.Fatalf("expected %v to equal %v", expectedPoint, testPoint)
		}
	}
}

func TestMandelbrotOutputFromListOfPointPixelPairs(t *testing.T) {
	var output PixelValue

	inputChan := make(chan WindowPoint, 4)
	outputChan := make(chan PixelValue, 4)

	expected := [...]PixelValue{
		PixelValue{0, 0, 2},
		PixelValue{1, 0, 10},
		PixelValue{0, 1, 10},
		PixelValue{1, 1, 10},
	}

	inputChan <- WindowPoint{x:0, y:0, point: -1.0 + 1.0i}
	inputChan <- WindowPoint{x:1, y:0, point: 0.0 + 1.0i}
	inputChan <- WindowPoint{x:0, y:1, point: -1.0 + 0.0i}
	inputChan <- WindowPoint{x:1, y:1, point: 0.0 + 0.0i}

	close(inputChan)
	go Mandy(inputChan, outputChan, 10, 2.0)

	for _, expectedPoint := range expected {
		output = <-outputChan
		if (expectedPoint != output) {
			t.Fatalf("output %v should be %v", output, expectedPoint)
		}
	}
}
