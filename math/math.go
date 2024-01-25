package math

import "image"

type Window struct {
	W, H int
}

func (win Window) Same_as(ref_win Window) bool {
	return (win.H == ref_win.H) && (win.W == ref_win.W)
}

type Rect[T int | float64] struct {
	X0, Y0, X1, Y1 T
}

func (view Rect[T]) Same_as(ref_view Rect[T]) bool {
	return (view.X0 == ref_view.X0) &&
		(view.Y0 == ref_view.Y0) &&
		(view.X1 == ref_view.X1) &&
		(view.Y1 == ref_view.Y1)
}

func (view Rect[T]) Fix_aspect(win Window) Rect[T] {
	fixed_rect := view

	ref_aspect_ratio := float64(win.W) / float64(win.H)
	aspect_ratio := float64(view.X1-view.X0) / float64(view.Y1-view.Y0)

	if aspect_ratio > ref_aspect_ratio {
		cy := (view.Y0 + view.Y1) / 2.0
		new_dy := float64(view.X1-view.X0) / ref_aspect_ratio
		fixed_rect.Y0 = T(float64(cy) - 0.5*float64(new_dy))
		fixed_rect.Y1 = T(float64(cy) + 0.5*float64(new_dy))
	} else {
		cx := (view.X0 + view.X1) / 2.0
		new_dx := float64(view.Y1-view.Y0) * ref_aspect_ratio
		fixed_rect.X0 = T(float64(cx) - 0.5*float64(new_dx))
		fixed_rect.X1 = T(float64(cx) + 0.5*float64(new_dx))
	}

	return fixed_rect
}

// Use an alias since we need the methods
type MathView = Rect[float64]

// We didn't really need a generic. I wanted to learn generics and
// initially thought Selection should be Rect[int] but as I wrote
// the code, it turned out Rect[float64] was cleaner
type Selection = Rect[float64]

// Interfaces don't allow us to declare variables
type MathematicalObject interface {
	Compute(MathView, Window)
	Image() *image.RGBA
}
