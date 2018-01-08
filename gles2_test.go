package gles2_test

import "testing"

import "github.com/JamesDunne/rpi-egl/bcm"
import gl "github.com/JamesDunne/rpi-egl/gles2"

func TestOpenGLInit(t *testing.T) {
	display := bcm.OpenDisplay()
	if display == nil {
		t.Fatal("display = nil")
	}
	err := gl.Init()
	if err != nil {
		t.Fatal(err)
	}
	display.Close()
}

func TestClear(t *testing.T) {
	display := bcm.OpenDisplay()
	if display == nil {
		t.Fatal("display = nil")
	}
	defer display.Close()
	err := gl.Init()
	if err != nil {
		t.Fatal(err)
	}
	gl.ClearColor(0.10, 0.33, 0.33, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	display.SwapBuffers()
}
