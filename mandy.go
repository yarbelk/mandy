package main

import (
	tb "github.com/nsf/termbox-go"
	"math/cmplx"
	"math"
	"os"
	"fmt"
	"errors"
)

var outputChars [16]tb.Cell = [16]tb.Cell{
	tb.Cell{'\u2592', 1, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 1*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 2*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 3*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 4*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 5*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 6*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 7*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 8*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 9*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 10*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 11*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 13*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 14*13, tb.ColorBlack},
	tb.Cell{'\u2592', 1 + 15*13, tb.ColorBlack},
}

type mandelState struct {
	radius, xStep, yStep, xMin, yMin, xMax, yMax float64
	limit int32
}

var (
	RADIUS float64 = 2.0
	YSTEP float64 = 0.05
	XSTEP float64 = 0.01
)


func convergantPoint(mandyVal *complex128) bool {
	return (cmplx.Abs(*mandyVal) < RADIUS)
}

func (m *mandelState) mandelbrot(point complex128) (distance int32, err error) {
	if m.limit <= 0 {
		return 0, errors.New("Mandy needs a max depth to recurse to")
	}
	var currentVal complex128 = point;

	for distance=0; distance <= m.limit; distance++ {
		currentVal = cmplx.Pow(currentVal, 2.0+0i) + point
		if ! convergantPoint(&currentVal) {
			return distance, nil
		}
	}
	return distance, nil
}

func (mandy *mandelState) setupMandy(termX, termY int) {
	mandy.yStep = (mandy.yMax - mandy.yMin) / float64(termY)
	mandy.xStep = (mandy.xMax - mandy.xMin) / float64(termX)
}

func (mandy *mandelState) mandyRun(termX, termY int) {
	var yCur, xCur float64
	var point complex128

	for y := 0 ; y < termY; y++ {
		for x := 0 ; x < termX; x++ {
			xCur = mandy.xMin + (mandy.xStep * float64(x))
			yCur = mandy.yMin + (mandy.yStep * float64(y))
			point = complex(xCur, yCur)
			distance, err := mandy.mandelbrot(point)
			if err != nil {
				fmt.Println(err.Error())
			}
			i := int32(math.Min(15, math.Max(0, math.Log2(float64(distance)))))
			cell := outputChars[i]
			tb.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
		}
	}
	tb.SetCursor(0, 0)
	tb.Flush()
	tb.HideCursor()

}

func main() {
	var mandy *mandelState = &mandelState{
		xMin: -1.0,
		yMin: 0.0,
		xMax: 0.25,
		yMax: 0.5,
		limit: 65536,
		radius: 2.3,
	}

	err := tb.Init()
	if err != nil {
		fmt.Fprint(os.Stderr,err.Error())
		os.Exit(1)
	}
	tb.SetOutputMode(tb.Output216)
	defer tb.Close()

	termX, termY := tb.Size()

	mandy.setupMandy(termX, termY)
	mandy.mandyRun(termX, termY)

	for  {
		event := tb.PollEvent()
		if event.Type == tb.EventKey && event.Key == tb.KeyEsc {
			break
		}
	}
}
