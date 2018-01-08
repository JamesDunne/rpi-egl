package gles2_test

import "testing"

import "github.com/JamesDunne/rpi-egl/bcm"
import "github.com/JamesDunne/rpi-egl/gles2"

func TestOpenGLInit(t *testing.T) {
	display := bcm.OpenDisplay()
	if display == nil {
		t.Fatal("display = nil")
	}
	err := gles2.Init()
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
	err := gles2.Init()
	if err != nil {
		t.Fatal(err)
	}
	gles2.ClearColor(0.10, 0.33, 0.33, 1.0)
	gles2.Clear(gles2.COLOR_BUFFER_BIT)
	display.SwapBuffers()
}
