package ui

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kghose/mandelbrot-go/math"
)

type ZoomEvent struct {
	dragging      bool
	drag_complete bool
	selection     math.Selection
}

func (ze ZoomEvent) normalize() math.Selection {
	w0 := ze.selection.X0
	w1 := ze.selection.X1

	if w0 > w1 {
		w0, w1 = w1, w0
	}

	h0 := ze.selection.Y0
	h1 := ze.selection.Y1

	if h0 > h1 {
		h0, h1 = h1, h0
	}

	return math.Selection{X0: w0, Y0: h0, X1: w1, Y1: h1}
}

func (ze ZoomEvent) AspectCorrectedSelection(win math.Window) math.Selection {
	return ze.normalize().Fix_aspect(win)
}

func (de *ZoomEvent) cursor_pos_callback(w *glfw.Window, xpos float64, ypos float64) {
	_, H := w.GetSize()

	// mouse Y pos reference is top of window

	// The end of the drag should always be where our cursor is
	de.selection.X1 = xpos
	de.selection.Y1 = float64(H) - ypos

	// The start of the drag should be where our cursor is if we are not
	// already in the middle of a drag
	if !de.dragging {
		de.selection.X0 = xpos
		de.selection.Y0 = float64(H) - ypos
	}
}

func (de *ZoomEvent) mouse_button_callback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if !de.dragging {
			de.dragging = true
		}
	}
	if action == glfw.Release {
		if de.dragging {
			de.drag_complete = true
		}
	}
}

func (de *ZoomEvent) mark_done() {
	de.dragging = false
	de.drag_complete = false
}
