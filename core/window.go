package core

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/golang/glog"
)

// WindowConfig is used by client applications to request a specific video mode
// from a monitor by calling InitWindow and passing it as an argument.
type WindowConfig struct {
	Name              string
	Monitor           *glfw.Monitor
	Width, Height, Hz int
	Fullscreen        bool
}

// WindowManager exposes windowing to client applications
type WindowManager struct {
	window         *glfw.Window
	cfg            WindowConfig
	active         bool
	cursorPosition mgl64.Vec2
}

var (
	windowManager *WindowManager
)

func init() {
	windowManager = &WindowManager{}

	err := glfw.Init()
	if err != nil {
		glog.Fatal(err)
	}
}

func GetWindowManager() *WindowManager {
	return windowManager
}

func (w *WindowManager) SetWindowConfig(cfg WindowConfig) {
	w.cfg = cfg
}

func (w *WindowManager) Start() {
	w.window = renderSystem.MakeWindow(w.cfg)
	w.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	w.installCallbacks()
	glfw.SwapInterval(1)
}

// WindowSize implements the WindowSystem interface
func (w *WindowManager) WindowSize() mgl32.Vec2 {
	return mgl32.Vec2{float32(w.cfg.Width), float32(w.cfg.Height)}
}

// ShouldClose implements the WindowSystem interface
func (w *WindowManager) ShouldClose() bool {
	return w.window.ShouldClose()
}

func (w *WindowManager) installCallbacks() {
	w.window.SetKeyCallback(inputManager.KeyCallback)
	w.window.SetMouseButtonCallback(inputManager.MouseButtonCallback)
	w.window.SetCursorPosCallback(inputManager.MouseMoveCallback)
	w.window.SetScrollCallback(inputManager.MouseScrollCallback)
	w.window.SetSizeCallback(w.SizeCallback)
}

// SizeCallback is called from glfw when the window is resized
func (w *WindowManager) SizeCallback(window *glfw.Window, width int, height int) {
	w.cfg.Width = width
	w.cfg.Height = height
}

// Stop implements the WindowSystem interface
func (w *WindowManager) Stop() {
	glog.Info("Stopping")
	glfw.Terminate()
}

// SetActive pauses/resumes the timer manager.
func (w *WindowManager) SetActive(active bool) {
	w.active = active

	if w.active == true {
		w.window.SetCursorPosCallback(inputManager.MouseMoveCallback)
		GetTimerManager().Start()
	} else {
		w.window.SetCursorPosCallback(nil)
		GetTimerManager().Pause()
	}
}

// CursorPosition implements the WindowSystem interface
func (w *WindowManager) CursorPosition() (float64, float64) {
	return w.cursorPosition.X(), w.cursorPosition.Y()
}
