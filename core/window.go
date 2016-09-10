package core

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

// WindowSystem is the interface that wraps all windowing, input and timer logic. These
// are generally provided together by a toolkit such as GLFW or SDL.
type WindowSystem interface {
	// Start starts the WindowSystem. This is where implementations should configure
	// timers and register input callbacks.
	Start()

	// Step is called by the application loop at the end of all other calls. This is
	// where implementations should swap their front/back buffers and poll for HID events.
	Step()

	// Stop is called by the application after its main loop returns. Implementations
	// should call their termination/cleanup logic here.
	Stop()

	// ShouldClose returns whether the application has scheduled its main loop termination.
	// Client applications should call this if they need to detect termination (eg: to save things).
	ShouldClose() bool

	// InitWindow configures the window properties. It allows the client application to specify
	// which monitor should be used to display the window in, the window's size and its title.
	InitWindow(name string, cfg WindowConfig)

	// WindowSize returns the current size of the managed window.
	WindowSize() mgl32.Vec2

	// Returns the current cursor position inside window (will scale and remap)
	CursorPosition() (float64, float64)

	// Monitors returns a list of monitors detected at startup.
	Monitors() []Monitor

	// GetTime returns the time elapsed since application startup and should be implemented
	// using a high resolution timer.
	GetTime() float64
}

// MonitorVideoMode is an interface used to access monitor capabilities.
type MonitorVideoMode interface {
	// Width returns the width in pixels of the MonitorVideoMode.
	Width() int

	// Height returns the height in pixels of the MonitorVideoMode.
	Height() int

	// Hz returns the refresh rate of the MonitorVideoMode.
	Hz() int
}

// Monitor is an interface which wraps information about a physical display device.
type Monitor interface {
	// Name returns the monitor's name.
	Name() string

	// VideoModes returns a list of video modes supported by the monitor.
	VideoModes() []MonitorVideoMode
}

// WindowConfig is used by client applications to request a specific video mode
// from a monitor by calling InitWindow and passing it as an argument.
type WindowConfig struct {
	Monitor           Monitor
	Width, Height, Hz int
	Fullscreen        bool
}

var (
	windowSystem WindowSystem
)

// SetWindowSystem should be called by implementations on their init function. It registers
// the implementation as the active window system.
func SetWindowSystem(w WindowSystem) {
	if windowSystem != nil {
		log.Fatal("Can't replace previously registered window system. Please make sure you're not importing twice")
	}
	windowSystem = w
}

// GetWindowSystem returns the registered window system
func GetWindowSystem() WindowSystem {
	return windowSystem
}
