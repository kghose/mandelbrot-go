package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
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

	ui := UI{}
	ui.init()
	mandelbrot_set := MandelbrotSet{}
	mandelbrot_set.max_iter = 1000
	ui.draw_loop(&mandelbrot_set)

}
