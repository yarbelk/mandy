package term_test

import (
	"testing"

	"github.com/yarbelk/mandy/lib/term"
)

// These tests are pretty bad.  I'vve cleaned them up somewhat, but
// its some of the first testing I did in golang, and i wasn't taking it
// too seriously becausee its a toy

/* Expect the chanel to spit out a set of complex numbers and
 * x/y coordinates at the right resolution and with y flipped
 */
func TestGenerateListOfPointPixelPairs(t *testing.T) {
	var windowLimits term.WindowLimits
	var testPoint, expectedPoint term.WindowPoint

	testChan := make(chan term.WindowPoint, 4)
	expected := []term.WindowPoint{
		term.WindowPoint{X: 0, Y: 0, Point: -1.0 + 1.0i},
		term.WindowPoint{X: 1, Y: 0, Point: 0.0 + 1.0i},
		term.WindowPoint{X: 0, Y: 1, Point: -1.0 + 0.0i},
		term.WindowPoint{X: 1, Y: 1, Point: 0.0 + 0.0i},
	}

	windowLimits = term.NewWindowLimits(2, 2, -1.0, 1.0, -1.0, 1.0)
	term.ProdWindowPoints(&windowLimits, testChan)
	for _, testPoint = range expected {

		expectedPoint = <-testChan
		if expectedPoint != testPoint {
			t.Logf("expected %v to equal %v", expectedPoint, testPoint)
			t.Fail()
		}
	}
}

func TestMandelbrotOutputFromListOfPointPixelPairs(t *testing.T) {
	var output term.ConverganceValue

	inputChan := make(chan term.WindowPoint, 4)
	outputChan := make(chan term.ConverganceValue, 4)

	expected := []term.ConverganceValue{
		term.ConverganceValue{0, 0, 6553},
		term.ConverganceValue{1, 0, 32767},
		term.ConverganceValue{0, 1, 32767},
		term.ConverganceValue{1, 1, 32767},
	}

	inputChan <- term.WindowPoint{X: 0, Y: 0, Point: -1.0 + 1.0i}
	inputChan <- term.WindowPoint{X: 1, Y: 0, Point: 0.0 + 1.0i}
	inputChan <- term.WindowPoint{X: 0, Y: 1, Point: -1.0 + 0.0i}
	inputChan <- term.WindowPoint{X: 1, Y: 1, Point: 0.0 + 0.0i}

	close(inputChan)
	go term.Mandy(inputChan, outputChan, 10, 2.0)

	for _, expectedPoint := range expected {
		output = <-outputChan
		if expectedPoint != output {
			t.Logf("output %v should be %v", output, expectedPoint)
			t.Fail()
		}
	}
}
