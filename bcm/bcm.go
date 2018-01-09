package bcm

// #cgo linux  CFLAGS: -I/opt/vc/include
// #cgo linux LDFLAGS: -L/opt/vc/lib -lEGL -lGLESv2 -lbcm_host
/*
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
} EGLDisplayState;

EGLNativeWindowType eglNativeWindow(EGLint *winWidth, EGLint *winHeight)
{
	static EGL_DISPMANX_WINDOW_T nativeWindow;
	DISPMANX_ELEMENT_HANDLE_T    element;
	DISPMANX_DISPLAY_HANDLE_T    display;
	DISPMANX_MODEINFO_T          mode;
	DISPMANX_UPDATE_HANDLE_T     update;
	VC_RECT_T                    rectDst, rectSrc;
	int32_t success;

	bcm_host_init();

	success = graphics_get_display_size(0, winWidth, winHeight);
	if (success < 0) {
		return NULL;
	}

	rectDst.x            = 0;
	rectDst.y            = 0;
	rectDst.width        = mode.width;
	rectDst.height       = mode.height;
	rectSrc.x            = 0;
	rectSrc.y            = 0;
	rectSrc.width        = *winWidth  << 16;
	rectSrc.height       = *winHeight << 16;
	nativeWindow.width   = *winWidth;
	nativeWindow.height  = *winHeight;

	display = vc_dispmanx_display_open(0);
	update  = vc_dispmanx_update_start(0);
	nativeWindow.element = vc_dispmanx_element_add(update,
						   display,
						   0,
						   &rectDst,
						   0,
						   &rectSrc,
						   DISPMANX_PROTECTION_NONE,
						   0,
						   0,
						   0);

	vc_dispmanx_update_submit_sync(update);
	//vc_dispmanx_display_close(display);

	return (EGLNativeWindowType) &nativeWindow;
}

EGLDisplayState *eglOpenDisplay()
{
	static EGLDisplayState state;
	EGLint attribList[] = {
		EGL_RED_SIZE,        5,
		EGL_GREEN_SIZE,      6,
		EGL_BLUE_SIZE,       5,
		EGL_ALPHA_SIZE,      EGL_DONT_CARE,
		EGL_DEPTH_SIZE,      EGL_DONT_CARE,
		EGL_STENCIL_SIZE,    EGL_DONT_CARE,
		EGL_SAMPLE_BUFFERS,  0,
		EGL_NONE
	};
	EGLint contextAttribs[] = { EGL_CONTEXT_CLIENT_VERSION, 2, EGL_NONE };
	EGLint    num;
	EGLConfig config;
	EGLNativeWindowType native_window;

	state.win.x      = 0;
	state.win.y      = 0;
	state.win.width  = 0;
	state.win.height = 0;

	native_window = eglNativeWindow(&state.win.width, &state.win.height);
	if (native_window == NULL) {
		return NULL;
	}

	state.display = eglGetDisplay(EGL_DEFAULT_DISPLAY);
	if (state.display == EGL_NO_DISPLAY) {
		return NULL;
	}
	if (!eglInitialize(state.display, NULL, NULL)) {
		return NULL;
	}
	if (!eglGetConfigs(state.display, NULL, 0, &num)) {
		return NULL;
	}
	if (!eglChooseConfig(state.display, attribList, &config, 1, &num)) {
		return NULL;
	}

	state.surface = eglCreateWindowSurface(state.display, config, native_window, NULL);
	if (state.surface == EGL_NO_SURFACE) {
		return NULL;
	}
	state.context = eglCreateContext(state.display, config, EGL_NO_CONTEXT, contextAttribs);
	if (state.context == EGL_NO_CONTEXT) {
		return NULL;
	}

	if (!eglMakeCurrent(state.display, state.surface, state.surface, state.context)) {
		return NULL;
	}

	return &state;
}

void eglCloseDisplay(EGLDisplayState *state)
{
	eglMakeCurrent(state->display, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT);
	eglDestroySurface(state->display, state->surface);
	eglDestroyContext(state->display, state->context);
	eglTerminate(state->display);
}

void eglUpdateDisplay(EGLDisplayState *state)
{
	eglSwapBuffers(state->display, state->surface);
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
	return fmt.Sprintf("egl errno=%d", e)
}

func OpenDisplay() (*Display, error) {
	state := C.eglOpenDisplay()
	if state == nil {
		return nil, getLastError()
	}

	return &Display{state: state}, nil
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

func (d *Display) SwapBuffers() {
	C.eglUpdateDisplay(d.state)
}
