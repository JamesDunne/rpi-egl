package bcm

// #cgo linux  CFLAGS: -I/opt/vc/include
// #cgo linux LDFLAGS: -L/opt/vc/lib -lEGL -lGLESv2 -lbcm_host
/*
#include <linux/kd.h>
#include <sys/ioctl.h>
#include <fcntl.h>
#include <stdio.h>

#include <bcm_host.h>
#include <EGL/egl.h>

typedef struct {
	EGLint x, y, width, height;
} EGLRect;

typedef struct {
	EGLRect    win;
	EGLDisplay display;
	EGLSurface surface;
	EGLContext context;

	DISPMANX_ELEMENT_HANDLE_T    dmx_element;
	DISPMANX_DISPLAY_HANDLE_T    dmx_display;
	DISPMANX_UPDATE_HANDLE_T     dmx_update;
	EGL_DISPMANX_WINDOW_T        nativeWindow;
} EGLDisplayState;

EGLint eglOpenDisplay(EGLDisplayState *state, EGLint red, EGLint green, EGLint blue)
{
	EGLint attribList[] = {
		EGL_RED_SIZE,          red,
		EGL_GREEN_SIZE,        green,
		EGL_BLUE_SIZE,         blue,
		EGL_ALPHA_SIZE,        EGL_DONT_CARE,
		EGL_DEPTH_SIZE,        EGL_DONT_CARE,
		EGL_STENCIL_SIZE,      EGL_DONT_CARE,
		EGL_MIN_SWAP_INTERVAL, 0,
		EGL_NONE
	};
	EGLint contextAttribs[] = { EGL_CONTEXT_CLIENT_VERSION, 2, EGL_NONE };
	EGLint    num;
	EGLConfig config;

	VC_RECT_T rectDst, rectSrc;

	int32_t success;

	bcm_host_init();

	state->win.x      = 0;
	state->win.y      = 0;
	state->win.width  = 0;
	state->win.height = 0;

	success = graphics_get_display_size(0, &state->win.width, &state->win.height);
	if (success < 0) {
		return 1;
	}

	rectDst.x            = 0;
	rectDst.y            = 0;
	rectDst.width        = state->win.width;
	rectDst.height       = state->win.height;

	rectSrc.x            = 0;
	rectSrc.y            = 0;
	rectSrc.width        = state->win.width  << 16;
	rectSrc.height       = state->win.height << 16;

	state->dmx_display = vc_dispmanx_display_open(0);
	state->dmx_update  = vc_dispmanx_update_start(0);

	state->nativeWindow.width   = state->win.width;
	state->nativeWindow.height  = state->win.height;
	state->nativeWindow.element = vc_dispmanx_element_add(
		state->dmx_update,
		state->dmx_display,
		0,
		&rectDst,
		0,
		&rectSrc,
		DISPMANX_PROTECTION_NONE,
		0,
		0,
		0
	);
	state->dmx_element = state->nativeWindow.element;

	success = vc_dispmanx_update_submit_sync(state->dmx_update);
	if (success != 0) {
		return 2;
	}

	// Create EGL context:
	state->display = eglGetDisplay(EGL_DEFAULT_DISPLAY);
	if (state->display == EGL_NO_DISPLAY) {
		return 3;
	}
	if (!eglInitialize(state->display, NULL, NULL)) {
		return 4;
	}
	if (!eglGetConfigs(state->display, NULL, 0, &num)) {
		return 5;
	}
	if (!eglChooseConfig(state->display, attribList, &config, 1, &num)) {
		return 6;
	}

	state->surface = eglCreateWindowSurface(state->display, config, (EGLNativeWindowType) &state->nativeWindow, NULL);
	if (state->surface == EGL_NO_SURFACE) {
		return 7;
	}
	state->context = eglCreateContext(state->display, config, EGL_NO_CONTEXT, contextAttribs);
	if (state->context == EGL_NO_CONTEXT) {
		return 8;
	}

	if (!eglMakeCurrent(state->display, state->surface, state->surface, state->context)) {
		return 9;
	}

	return 0;
}

EGLBoolean eglCloseDisplay(EGLDisplayState *state)
{
	if (!eglMakeCurrent(state->display, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT)) {
		return EGL_FALSE;
	}
	if (!eglDestroySurface(state->display, state->surface)) {
		return EGL_FALSE;
	}
	if (!eglDestroyContext(state->display, state->context)) {
		return EGL_FALSE;
	}
	if (!eglTerminate(state->display)) {
		return EGL_FALSE;
	}
	if (vc_dispmanx_element_remove(state->dmx_update, state->dmx_element) != 0) {
		return EGL_FALSE;
	}
	if (vc_dispmanx_update_submit_sync(state->dmx_update) != 0) {
		return EGL_FALSE;
	}
	if (vc_dispmanx_display_close(state->dmx_display) != 0) {
		return EGL_FALSE;
	}
	return EGL_TRUE;
}

EGLBoolean eglUpdateDisplay(EGLDisplayState *state)
{
	return eglSwapBuffers(state->display, state->surface);
}

void ttyGraphics(const char *tty)
{
	// Need sudo or tty group for this:
	int kbfd = open(tty, O_WRONLY);
	if (kbfd >= 0) {
		ioctl(kbfd, KDSETMODE, KD_GRAPHICS);
		close(kbfd);
	}
}

void ttyText(const char *tty)
{
	// Need sudo or tty group for this:
	int kbfd = open(tty, O_WRONLY);
	if (kbfd >= 0) {
		ioctl(kbfd, KDSETMODE, KD_TEXT);
		close(kbfd);
	}
}
*/
import "C"

import "fmt"

type Display struct {
	state *C.EGLDisplayState
}

type EGLError int32

func getLastError() EGLError {
	return EGLError(C.eglGetError())
}

func (e EGLError) Error() string {
	return fmt.Sprintf("egl error=0x%04x", int32(e))
}

func OpenDisplay(red, green, blue int) (*Display, error) {
	state := (*C.EGLDisplayState)(C.malloc(C.sizeof_EGLDisplayState))
	step := C.eglOpenDisplay(state, C.EGLint(red), C.EGLint(green), C.EGLint(blue))
	if step != 0 {
		fmt.Printf("step %d\n", step)
		return nil, getLastError()
	}

	return &Display{state: state}, nil
}

func TTYGraphicsMode(tty string) {
	C.ttyGraphics(C.CString(tty))
}

func TTYTextMode(tty string) {
	C.ttyText(C.CString(tty))
}

func (d *Display) Close() {
	C.eglCloseDisplay(d.state)
}

func (d *Display) Width() int {
	return int(d.state.win.width)
}

func (d *Display) Height() int {
	return int(d.state.win.height)
}

func (d *Display) SwapBuffers() error {
	if C.eglUpdateDisplay(d.state) == 0 {
		return getLastError()
	}
	return nil
}

func (d *Display) SwapInterval(interval int32) error {
	if C.eglSwapInterval(d.state.display, C.EGLint(interval)) == 0 {
		return getLastError()
	}
	return nil
}
