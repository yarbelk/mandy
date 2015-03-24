package main

import (
	tb "github.com/nsf/termbox-go"
	"math/cmplx"
	"math"
	"os"
	"fmt"
	"sync"
)

var outputChars [16]tb.Cell = [16]tb.Cell{
	tb.Cell{'\u2592', 0, tb.ColorBlack},
	tb.Cell{'\u2592', 1*18, tb.ColorBlack},
	tb.Cell{'\u2592', 2*18, tb.ColorBlack},
	tb.Cell{'\u2592', 3*18, tb.ColorBlack},
	tb.Cell{'\u2592', 4*18, tb.ColorBlack},
	tb.Cell{'\u2592', 5*18, tb.ColorBlack},
	tb.Cell{'\u2592', 6*18, tb.ColorBlack},
	tb.Cell{'\u2592', 7*18, tb.ColorBlack},
	tb.Cell{'\u2592', 8*18, tb.ColorBlack},
	tb.Cell{'\u2592', 9*18, tb.ColorBlack},
	tb.Cell{'\u2592', 10*18, tb.ColorBlack},
}

type screenPoint struct {
	x, y int
	point complex128
}

type mandelState struct {
	radius, xStep, yStep, xMin, yMin, xMax, yMax float64
	limit int32
}

func (m* mandelState)convergantPoint(mandyVal *complex128) bool {
	return (cmplx.Abs(*mandyVal) < m.radius)
}

func (m *mandelState) mandelbrot(pointChan <-chan screenPoint, wg sync.WaitGroup) {
	defer wg.Done()
	var distance int32
	if m.limit <= 0 {
		os.Exit(1)
	}
	
	for point := range pointChan {
		var currentVal complex128 = point.point
		for distance=0; distance <= m.limit; distance++ {
			currentVal = cmplx.Pow(currentVal, 2.0+0i) + point.point
			if ! m.convergantPoint(&currentVal) {
				m.printMandyPoint(distance, point.x, point.y)
			}
		}
	}
}

func (m *mandelState) printMandyPoint(distance int32, x, y int) {
	i := int32(math.Min(15, math.Max(0, math.Log2(float64(distance)))))
	cell := outputChars[i]
	tb.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
}


func (m *mandelState) setupMandy(termX, termY int) {
	m.yStep = (m.yMax - m.yMin) / float64(termY)
	m.xStep = (m.xMax - m.xMin) / float64(termX)
}

func (mandy *mandelState) mandyRun(termX, termY int) {
	var yCur, xCur float64
	var point complex128
	var wg sync.WaitGroup

	pointChan := make(chan screenPoint)
	wg.Add(1)
	go mandy.mandelbrot(pointChan, wg)
	for y := 0 ; y < termY; y++ {
		for x := 0 ; x < termX; x++ {
			xCur = mandy.xMin + (mandy.xStep * float64(x))
			yCur = mandy.yMin + (mandy.yStep * float64(y))
			point = complex(xCur, yCur)
			pointChan <- screenPoint{x: x, y: y, point: point}
		}
	}
	close(pointChan)
	wg.Wait()
	tb.SetCursor(0, 0)
	tb.Flush()
	tb.HideCursor()

}

func main() {
	var mandy *mandelState = &mandelState{
		xMin: -2.0,
		yMin: -1.0,
		xMax: 0.5,
		yMax: 1.0,
		limit: 4096,
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
