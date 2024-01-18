package main

import (
	"fmt"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type UI struct {
	window      *glfw.Window
	texture     uint32
	framebuffer uint32
}

func (ui *UI) init() {
	var err error
	ui.window, err = glfw.CreateWindow(640, 480, "My Window", nil, nil)
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

}

func (ui *UI) draw_loop(mandelbrot_set MathematicalObject) {

	W, H := ui.window.GetSize()
	existing_viewport := ViewPort{W, H, -2.0, -1.0, 1.0, -1.0 + 3.0*float64(H)/float64(W)}
	new_viewport := ViewPort{}

	for !ui.window.ShouldClose() {

		new_viewport.W, new_viewport.H = ui.window.GetSize()
		new_viewport.x0, new_viewport.y0 = existing_viewport.x0, existing_viewport.y0
		dx := (existing_viewport.x1 - existing_viewport.x0) / float64(existing_viewport.W)

		new_viewport.x1 = new_viewport.x0 + dx*float64(new_viewport.W)
		new_viewport.y1 = new_viewport.y0 + dx*float64(new_viewport.H)

		gl.BindTexture(gl.TEXTURE_2D, ui.texture)
		gl.TexImage2D(
			gl.TEXTURE_2D, 0, gl.RGBA8,
			int32(new_viewport.W), int32(new_viewport.H),
			0, gl.RGBA, gl.UNSIGNED_BYTE,
			gl.Ptr(mandelbrot_set.update(new_viewport).Pix),
		)
		gl.BlitFramebuffer(
			0, 0,
			int32(new_viewport.W), int32(new_viewport.H),
			0, 0,
			int32(new_viewport.W), int32(new_viewport.H),
			gl.COLOR_BUFFER_BIT, gl.LINEAR,
		)

		new_viewport = mandelbrot_set.true_viewport()
		ui.window.SetTitle(fmt.Sprintf("(%f, %f) - (%f, %f)", new_viewport.x0, new_viewport.y0, new_viewport.x1, new_viewport.y1))

		ui.window.SwapBuffers()
		glfw.PollEvents()
	}

}

// We choose to normalize on the x-axis
func (v *ViewPort) normalize() {
	dx := (v.x1 - v.x0) / float64(v.W)
	v.y1 = v.y0 + dx*float64(v.H)
}
