package math

import (
	"image"
	"image/color"
)

type MandelbrotSet struct {
	view     MathView
	win      Window
	Max_iter int
	img      *image.RGBA
}

func (mandel *MandelbrotSet) Image() *image.RGBA {
	return mandel.img
}

func (mandel *MandelbrotSet) Compute(new_view MathView, new_win Window) {

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

	for i := 0; i < mandel.win.W; i++ {
		for j := 0; j < mandel.win.H; j++ {
			k := escape_number(x0+dx*float64(i), y0+dy*float64(j), mandel.Max_iter)
			mandel.img.Set(i, j, color.Gray16{k})
		}
	}
}

func escape_number(x0 float64, y0 float64, max_iter int) uint16 {
	col_scale := float64(0xffff) / float64(max_iter)
	x := 0.0
	y := 0.0
	for i := 0; i < max_iter; i++ {
		xtemp := x*x - y*y + x0
		y = 2*x*y + y0
		x = xtemp
		if x*x+y*y >= 4.0 {
			return uint16((max_iter - i) * int(col_scale))
		}
	}
	return 0x0
}