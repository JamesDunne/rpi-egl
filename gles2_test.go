package gles2_test

import (
	"testing"
	"time"
)

import "github.com/JamesDunne/rpi-egl/bcm"
import gl "github.com/JamesDunne/rpi-egl/gles2"

func TestOpenGLInit(t *testing.T) {
	display, err := bcm.OpenDisplay()
	if err != nil {
		t.Fatal(err)
	}
	display.Close()
}

func TestClear(t *testing.T) {
	display, err := bcm.OpenDisplay()
	if err != nil {
		t.Fatal(err)
	}
	defer display.Close()

	gl.ClearColor(0.10, 0.33, 0.33, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	display.SwapBuffers()
	time.Sleep(time.Second * 1)
}
