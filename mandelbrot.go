package main

import (
	"image"
	"image/color"
)

/*
Physical viewport to mathematical viewport:

Physical viewport:          L, W
Mathematical viewport:      x, y, x+kL, y+kW
Initial scaling factor k = .01
Initial physical viewport L = 400, W = 400
Initial mathematical viewport x = -2, y = -2

*/

func (existing_view ViewPort) same_as(new_view ViewPort) bool {
	return (existing_view.x0 == new_view.x0) &&
		(existing_view.y0 == new_view.y0) &&
		(existing_view.x1 == new_view.x1) &&
		(existing_view.y1 == new_view.y1) &&
		(existing_view.W == new_view.W) &&
		(existing_view.H == new_view.H)
}

type MandelbrotSet struct {
	current_view ViewPort
	max_iter     int
	img          *image.RGBA
}

func (set *MandelbrotSet) update(new_view ViewPort) *image.RGBA {

	if set.current_view.same_as(new_view) {
		return set.img
	}
	if set.current_view.H != new_view.H || set.current_view.W != new_view.W {
		set.img = image.NewRGBA(image.Rect(0, 0, new_view.W, new_view.H))
	}
	set.current_view = new_view

	x0 := set.current_view.x0
	y0 := set.current_view.y0
	dx := (set.current_view.x1 - x0) / float64(set.current_view.W)
	dy := (set.current_view.y1 - y0) / float64(set.current_view.H)

	for i := 0; i < set.current_view.W; i++ {
		for j := 0; j < set.current_view.H; j++ {
			k := escape_number(x0+dx*float64(i), y0+dy*float64(j), set.max_iter)
			set.img.Set(i, j, color.Gray16{k})
		}
	}
	return set.img
}

func (set *MandelbrotSet) true_viewport() ViewPort {
	return set.current_view
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
