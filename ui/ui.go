package ui

import (
	"fmt"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kghose/mandelbrot-go/math"
)

type UI struct {
	view        math.MathView
	zoom_event  ZoomEvent
	window      *glfw.Window
	texture     uint32
	framebuffer uint32
}

func (ui *UI) Init() {
	var err error
	ui.window, err = glfw.CreateWindow(1000, 700, "Mandelbrot Set", nil, nil)
	if err != nil {
		panic(err)
	}

	ui.window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	{
		gl.GenTextures(1, &ui.texture)

		gl.BindTexture(gl.TEXTURE_2D, ui.texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

		gl.BindImageTexture(0, ui.texture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA8)
	}

	{
		gl.GenFramebuffers(1, &ui.framebuffer)
		gl.BindFramebuffer(gl.FRAMEBUFFER, ui.framebuffer)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, ui.texture, 0)

		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, ui.framebuffer)
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	}

	ui.zoom_event = ZoomEvent{}
	ui.window.SetCursorPosCallback(ui.zoom_event.cursor_pos_callback)
	ui.window.SetMouseButtonCallback(ui.zoom_event.mouse_button_callback)
	ui.window.SetScrollCallback(ui.scroll_callback)

}

func (ui *UI) Draw_loop(mandelbrot_set math.MathematicalObject) {

	W, H := ui.window.GetSize()
	win := math.Window{W: W, H: H}
	ui.view = math.MathView{X0: -2.5, Y0: -1.0, X1: 1.0, Y1: 1.0}.Fix_aspect(win)

	for !ui.window.ShouldClose() {

		W, H := ui.window.GetSize()
		win := math.Window{W: W, H: H}
		ui.view = ui.view.Fix_aspect(win)

		if ui.zoom_event.drag_complete {
			ui.view = ui.ZoomedView(win)
			ui.zoom_event.mark_done()
		}
		ui.draw_object(mandelbrot_set)

		if ui.zoom_event.dragging {
			ui.draw_drag()
		}

		ui.window.SetTitle(
			fmt.Sprintf(
				"Mandelbrot Set: (%f, %f) - (%f, %f) %d x %d",
				ui.view.X0, ui.view.Y0,
				ui.view.X1, ui.view.Y1,
				W, H))

		ui.window.SwapBuffers()
		glfw.PollEvents()
	}

}

func (ui *UI) draw_object(mandelbrot_set math.MathematicalObject) {

	W, H := ui.window.GetSize()
	win := math.Window{W: W, H: H}

	mandelbrot_set.Compute(ui.view, win)

	gl.BindTexture(gl.TEXTURE_2D, ui.texture)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA8,
		int32(W), int32(H),
		0, gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(mandelbrot_set.Image().Pix),
	)
	gl.BlitFramebuffer(
		0, 0,
		int32(W), int32(H),
		0, 0,
		int32(W), int32(H),
		gl.COLOR_BUFFER_BIT, gl.LINEAR,
	)

}

// https://stackoverflow.com/a/21451101
func (ui *UI) draw_drag() {

	W, H := ui.window.GetSize()
	win := math.Window{W: W, H: H}
	selection := ui.zoom_event.AspectCorrectedSelection(win)

	gl.Color3f(1.0, 0.2, 0.2)
	gl.LineWidth(2.0)

	gl.PushMatrix()
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(W), 0, float64(H), -1.0, 1.0)
	gl.Viewport(0, 0, int32(W), int32(H))
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex2d(selection.X0, selection.Y0)
	gl.Vertex2d(selection.X1, selection.Y0)
	gl.Vertex2d(selection.X1, selection.Y1)
	gl.Vertex2d(selection.X0, selection.Y1)
	gl.End()
	gl.PopMatrix()
}

func (ui *UI) scroll_callback(w *glfw.Window, xoff float64, yoff float64) {
	W, H := ui.window.GetSize()
	ui.view.X0 -= xoff * float64(W) * 0.01
	ui.view.X1 -= xoff * float64(W) * 0.01
	ui.view.Y0 += yoff * float64(H) * 0.01
	ui.view.Y1 += yoff * float64(H) * 0.01
}

func (ui *UI) ZoomedView(win math.Window) math.MathView {

	selection := ui.zoom_event.AspectCorrectedSelection(win)

	dx := (ui.view.X1 - ui.view.X0) / float64(win.W)
	dy := (ui.view.Y1 - ui.view.Y0) / float64(win.H)

	new_x0 := ui.view.X0 + dx*float64(selection.X0)
	new_x1 := ui.view.X0 + dx*float64(selection.X1)
	new_y0 := ui.view.Y0 + dy*float64(selection.Y0)
	new_y1 := ui.view.Y0 + dy*float64(selection.Y1)

	return math.MathView{
		X0: new_x0,
		Y0: new_y0,
		X1: new_x1,
		Y1: new_y1,
	}

}
