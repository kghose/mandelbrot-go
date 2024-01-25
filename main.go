package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kghose/mandelbrot-go/math"
	"github.com/kghose/mandelbrot-go/ui"
)

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

	mandel_ui := ui.UI{}
	mandel_ui.Init()
	mandelbrot_set := math.MandelbrotSet{}
	mandelbrot_set.Max_iter = 1000
	mandel_ui.Draw_loop(&mandelbrot_set)

}
