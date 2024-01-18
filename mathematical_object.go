package main

import "image"

type ViewPort struct {
	W, H           int
	x0, y0, x1, y1 float64
}

type MathematicalObject interface {
	update(view ViewPort) *image.RGBA
	true_viewport() ViewPort
}
