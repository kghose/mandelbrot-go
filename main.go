package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func cursor_pos_callback(w *glfw.Window, xpos float64, ypos float64) {
	// w.SetTitle(fmt.Sprintf("Mouse: %f, %f", xpos, ypos))
}

func scroll_callback(w *glfw.Window, xoff float64, yoff float64) {
	// w.SetTitle(fmt.Sprintf("Mouse: %f, %f", xoff, yoff))
}

type UserInput struct {
	dragging      bool
	drag_complete bool
	x,
	y,
	drag_startx,
	drag_starty,
	drag_endx,
	drag_endy int
}

func (ui *UserInput) mouse_button_callback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if !ui.dragging {
			ui.drag_startx = ui.x
			ui.drag_starty = ui.y
		}
	}
	if action == glfw.Release {
		if !ui.dragging {
			ui.drag_startx = ui.x
			ui.drag_starty = ui.y
		}
	}

}

func init() {
	// GLFW: This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	ui := UI{}
	ui.init()
	mandelbrot_set := MandelbrotSet{}
	mandelbrot_set.max_iter = 100
	ui.draw_loop(&mandelbrot_set)

}
