package mandelbrot

import (
	"testing"

	"github.com/kghose/mandelbrot-go/math"
)

func BenchmarkMandelbrot(b *testing.B) {
	mandelbrot_set := MandelbrotSet{Max_iter: 1000}
	win := math.Window{W: 500, H: 500}
	view := math.MathView{X0: -2.5, Y0: -1.0, X1: 1.0, Y1: 1.0}

	for i := 0; i < b.N; i++ {
		mandelbrot_set.Compute(view, win)
	}
}
