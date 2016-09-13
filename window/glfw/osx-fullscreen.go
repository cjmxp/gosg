// +build darwin

package glfw

// #cgo darwin CFLAGS: -D DARWIN -x objective-c
// #cgo darwin LDFLAGS: -framework Cocoa -framework CoreFoundation
/*
#ifdef DARWIN
void toggleFullScreen(void *window);
#include <stdlib.h>
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
void toggleFullScreen(void *window) {
NSWindow *nsWindow = (NSWindow *)window;
    [nsWindow toggleFullScreen:nil];
}
#endif*/
import "C"

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.2/glfw"
)

func makeFullScreen(w *glfw.Window) {
	// this exists because GLFW's full screen is terribly broken in el capitan (at least)
	C.toggleFullScreen(unsafe.Pointer(w.GetCocoaWindow()))
}
