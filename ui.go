package main

import (
	"fmt"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type UI struct {
	view_port   ViewPort
	drag_event  DragEvent
	window      *glfw.Window
	texture     uint32
	framebuffer uint32
}

func (ui *UI) init() {
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

	ui.drag_event = DragEvent{}
	ui.window.SetCursorPosCallback(ui.drag_event.cursor_pos_callback)
	ui.window.SetMouseButtonCallback(ui.drag_event.mouse_button_callback)
	ui.window.SetScrollCallback(ui.scroll_callback)

}

func (ui *UI) draw_loop(mandelbrot_set MathematicalObject) {

	W, H := ui.window.GetSize()
	DY := float64(H) / float64(W)
	ui.view_port = ViewPort{W, H, -2.5, -1.0 * DY, 1.0, 1.0 * DY}

	for !ui.window.ShouldClose() {

		if ui.drag_event.drag_complete {
			ui.view_port = ui.drag_event.new_viewport(ui.view_port)
			ui.drag_event.mark_done()
		}
		ui.draw_object(mandelbrot_set)

		if ui.drag_event.dragging {
			ui.draw_drag()
		}

		ui.window.SetTitle(
			fmt.Sprintf(
				"Mandelbrot Set: (%f, %f) - (%f, %f) %d x %d",
				ui.view_port.x0, ui.view_port.y0,
				ui.view_port.x1, ui.view_port.y1,
				ui.view_port.W, ui.view_port.H))

		ui.window.SwapBuffers()
		glfw.PollEvents()
	}

}

func (ui *UI) draw_object(mandelbrot_set MathematicalObject) {

	ui.view_port.W, ui.view_port.H = ui.window.GetSize()
	dx := (ui.view_port.x1 - ui.view_port.x0) / float64(ui.view_port.W)
	cy := (ui.view_port.y1 + ui.view_port.y0) / 2.0
	ui.view_port.y0 = cy - 0.5*dx*float64(ui.view_port.H)
	ui.view_port.y1 = cy + 0.5*dx*float64(ui.view_port.H)

	gl.BindTexture(gl.TEXTURE_2D, ui.texture)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA8,
		int32(ui.view_port.W), int32(ui.view_port.H),
		0, gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(mandelbrot_set.update(ui.view_port).Pix),
	)
	gl.BlitFramebuffer(
		0, 0,
		int32(ui.view_port.W), int32(ui.view_port.H),
		0, 0,
		int32(ui.view_port.W), int32(ui.view_port.H),
		gl.COLOR_BUFFER_BIT, gl.LINEAR,
	)

}

// https://stackoverflow.com/a/21451101
func (ui *UI) draw_drag() {
	w0, h0, w1, h1 := ui.drag_event.normalize_drag_rect(ui.view_port)

	gl.Color3f(1.0, 1.0, 0.0)
	gl.LineWidth(5.0)

	gl.PushMatrix()
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(ui.view_port.W), 0, float64(ui.view_port.H), -1.0, 1.0)
	gl.Viewport(0, 0, int32(ui.view_port.W), int32(ui.view_port.H))
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex2d(w0, h0)
	gl.Vertex2d(w1, h0)
	gl.Vertex2d(w1, h1)
	gl.Vertex2d(w0, h1)
	gl.End()
	gl.PopMatrix()
}

func (ui *UI) scroll_callback(w *glfw.Window, xoff float64, yoff float64) {
	ui.view_port.x0 -= xoff * (ui.view_port.x1 - ui.view_port.x0) * 0.01
	ui.view_port.x1 -= xoff * (ui.view_port.x1 - ui.view_port.x0) * 0.01
	ui.view_port.y0 += yoff * (ui.view_port.y1 - ui.view_port.y0) * 0.01
	ui.view_port.y1 += yoff * (ui.view_port.y1 - ui.view_port.y0) * 0.01
}

type DragEvent struct {
	dragging      bool
	drag_complete bool
	win_x0,
	win_y0,
	win_x1,
	win_y1 float64
}

func (de *DragEvent) cursor_pos_callback(w *glfw.Window, xpos float64, ypos float64) {
	// The end of the drag should always be where our cursor is
	de.win_x1 = xpos
	de.win_y1 = ypos

	// The start of the drag should be where our cursor is if we are not
	// already in the middle of a drag
	if !de.dragging {
		de.win_x0 = xpos
		de.win_y0 = ypos
	}
}

func (de *DragEvent) mouse_button_callback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
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

func (de *DragEvent) mark_done() {
	de.dragging = false
	de.drag_complete = false
}

func (de *DragEvent) normalize_drag_rect(vp ViewPort) (w0 float64, h0 float64, w1 float64, h1 float64) {

	w0 = de.win_x0
	w1 = de.win_x1

	if w0 > w1 {
		w0, w1 = w1, w0
	}

	// mouse Y pos reference is top of window
	h0 = float64(vp.H) - de.win_y0
	h1 = float64(vp.H) - de.win_y1

	if h0 > h1 {
		h0, h1 = h1, h0
	}

	dw := w1 - w0
	dh := h1 - h0

	wr := dw / float64(vp.W)
	hr := dh / float64(vp.H)

	if wr < hr {
		cw := (w0 + w1) / 2.0
		w0 = cw - 0.5*hr*float64(vp.W)
		w1 = cw + 0.5*hr*float64(vp.W)
	} else {
		ch := (h0 + h1) / 2.0
		h0 = ch - 0.5*wr*float64(vp.H)
		h1 = ch + 0.5*wr*float64(vp.H)
	}

	return w0, h0, w1, h1
}

func (de *DragEvent) new_viewport(vp ViewPort) ViewPort {

	w0, h0, w1, h1 := de.normalize_drag_rect(vp)
	dx := (vp.x1 - vp.x0) / float64(vp.W)
	dy := (vp.y1 - vp.y0) / float64(vp.H)

	new_x0 := vp.x0 + dx*w0
	new_x1 := vp.x0 + dx*w1
	new_y0 := vp.y0 + dy*h0
	new_y1 := vp.y0 + dy*h1

	return ViewPort{
		vp.W,
		vp.H,
		new_x0,
		new_y0,
		new_x1,
		new_y1,
	}

}
