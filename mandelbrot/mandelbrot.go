package mandelbrot

import (
	"image"
	"image/color"
	"sync"

	"github.com/kghose/mandelbrot-go/math"
)

type MandelbrotSet struct {
	view     math.MathView
	win      math.Window
	Max_iter int
	img      *image.RGBA
}

func (mandel *MandelbrotSet) Image() *image.RGBA {
	return mandel.img
}

func (mandel *MandelbrotSet) Compute(new_view math.MathView, new_win math.Window) {

	if mandel.view.Same_as(new_view) && mandel.win.Same_as(new_win) {
		return
	}
	if !mandel.win.Same_as(new_win) {
		mandel.img = image.NewRGBA(image.Rect(0, 0, new_win.W, new_win.H))
	}
	mandel.view = new_view
	mandel.win = new_win

	x0 := mandel.view.X0
	y0 := mandel.view.Y0
	dx := (mandel.view.X1 - x0) / float64(mandel.win.W)
	dy := (mandel.view.Y1 - y0) / float64(mandel.win.H)

	wg := sync.WaitGroup{}
	for i := 0; i < mandel.win.W; i++ {
		wg.Add(1)
		go mandel.compute_column(i, x0, y0, dx, dy, &wg)
	}
	wg.Wait()
}

func (mandel *MandelbrotSet) compute_column(i int, x0 float64, y0 float64, dx float64, dy float64, wg *sync.WaitGroup) {
	for j := 0; j < mandel.win.H; j++ {
		k := escape_number(x0+dx*float64(i), y0+dy*float64(j), mandel.Max_iter)
		mandel.img.Set(i, j, color.Gray16{k})
	}
	wg.Done()
}

func escape_number(x0 float64, y0 float64, max_iter int) uint16 {
	col_scale := float64(0xffff) / float64(max_iter)
	x, x_2 := 0.0, 0.0
	y, y_2 := 0.0, 0.0
	for i := 0; i < max_iter; i++ {
		x_2, y_2 = x*x, y*y
		xtemp := x_2 - y_2 + x0
		y = 2*x*y + y0
		x = xtemp
		if x_2+y_2 >= 4.0 {
			return uint16((max_iter - i) * int(col_scale))
		}
	}
	return 0x0
}
