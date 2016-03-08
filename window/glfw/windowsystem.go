// Package glfw implements a GLFW based windowsystem for gosg
package glfw

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fcvarela/gosg/core"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/glog"
)

// WindowSystem implements the core.WindowSystem interface
type WindowSystem struct {
	name          string
	window        *glfw.Window
	cfg           core.WindowConfig
	active        bool
	cursorVisible bool
}

// Monitor implements the core.MonitorVideoMode interface
type Monitor struct {
	monitor *glfw.Monitor
}

// MonitorVideoMode wraps a GLFW VidMode and implements the core.MonitorVideoMode interface
type MonitorVideoMode struct {
	videoMode *glfw.VidMode
}

func init() {
	core.SetWindowSystem(New())
}

// Width implements the core.MonitorVideoMode interface
func (mv *MonitorVideoMode) Width() int {
	return mv.videoMode.Width
}

// Height implements the core.MonitorVideoMode interface
func (mv *MonitorVideoMode) Height() int {
	return mv.videoMode.Height
}

// Hz implements the core.MonitorVideoMode interface
func (mv *MonitorVideoMode) Hz() int {
	return mv.videoMode.RefreshRate
}

// String implements the core.MonitorVideoMode interface
func (mv *MonitorVideoMode) String() string {
	return fmt.Sprintf("%dx%d@%d", mv.Width(), mv.Height(), mv.Hz())
}

// VideoModes returns the core.Monitor interface
func (m *Monitor) VideoModes() []core.MonitorVideoMode {
	videoModes := m.monitor.GetVideoModes()

	myVideoModes := make([]core.MonitorVideoMode, len(videoModes))
	for i := range videoModes {
		myVideoModes[i] = &MonitorVideoMode{videoModes[i]}
	}

	return myVideoModes
}

// Name returns the monitor's name
func (m *Monitor) Name() string {
	return m.monitor.GetName()
}

var (
	focusOutCursorX float64
	focusOutCursorY float64
)

func maxi(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func mini(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// New returns a new WindowSystem
func New() *WindowSystem {
	err := glfw.Init()
	if err != nil {
		glog.Fatal(err)
	}

	return &WindowSystem{}
}

// InitWindow implements the core.WindowSystem interface
func (w *WindowSystem) InitWindow(name string, cfg core.WindowConfig) {
	w.name = name
	w.cfg = cfg
}

// Start implements the core.WindowSystem interface
func (w *WindowSystem) Start() {
	glog.Info("Starting")

	// map core input keys o ours
	w.SetInputMap()

	// create a window
	glfw.WindowHint(glfw.Decorated, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 0)

	monitor := w.cfg.Monitor.(*Monitor)
	var err error

	if w.cfg.Fullscreen && runtime.GOOS != "darwin" {
		w.window, err = glfw.CreateWindow(w.cfg.Width, w.cfg.Height, w.name, monitor.monitor, nil)
	} else if w.cfg.Fullscreen && runtime.GOOS == "darwin" {
		w.window, err = glfw.CreateWindow(w.cfg.Width/2, w.cfg.Height/2, w.name, nil, nil)
		makeFullScreen(w.window)
		time.Sleep(1 * time.Second)
	} else {
		w.window, err = glfw.CreateWindow(w.cfg.Width, w.cfg.Height, w.name, nil, nil)
	}
	if err != nil {
		glog.Fatal(err)
	}

	w.window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	glog.Info("Checking GL Init status")
	glog.Info("OpenGL version: ", gl.GoStr(gl.GetString(gl.VERSION)))

	w.installCallbacks()
	glfw.SwapInterval(1)
}

// Monitors implements the core.WindowSystem interface
func (w *WindowSystem) Monitors() []core.Monitor {
	monitors := glfw.GetMonitors()

	myMonitors := make([]core.Monitor, len(monitors))
	for i := range monitors {
		myMonitors[i] = &Monitor{monitors[i]}
	}

	return myMonitors
}

// WindowSize implements the core.WindowSystem interface
func (w *WindowSystem) WindowSize() mgl32.Vec2 {
	return mgl32.Vec2{float32(w.cfg.Width), float32(w.cfg.Height)}
}

// ShouldClose implements the core.WindowSystem interface
func (w *WindowSystem) ShouldClose() bool {
	return w.window.ShouldClose()
}

// SetCursorVisible implements the core.WindowSystem interface
func (w *WindowSystem) SetCursorVisible(v bool) {
	w.cursorVisible = v
	if v {
		w.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		w.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	}
}

func (w *WindowSystem) installCallbacks() {
	w.window.SetKeyCallback(w.KeyCallback)

	w.window.SetMouseButtonCallback(w.MouseButtonCallback)
	w.window.SetSizeCallback(w.SizeCallback)
	w.window.SetCursorPosCallback(w.MouseMoveCallback)
	w.window.SetScrollCallback(w.MouseScrollCallback)

	w.window.SetIconifyCallback(w.IconifyCallback)
	w.window.SetFocusCallback(w.FocusCallback)
}

// Step implements the core.WindowSystem interface
func (w *WindowSystem) Step() {
	// swap buffers. this should be moved to a separate call and called by the rendersystem at the end of its work
	w.window.SwapBuffers()

	// reset input state before callbacks
	core.GetInputManager().Reset()

	//if core.GetTimerManager().Paused() {
	// this doesn't work properly in Windows
	//	glfw.WaitEvents()
	//} else {
	glfw.PollEvents()
	//}
}

// Stop implements the core.WindowSystem interface
func (w *WindowSystem) Stop() {
	glog.Info("Stopping")
	glfw.Terminate()
}

// KeyCallback passes keyboard events to the core.InputManager
func (w *WindowSystem) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	core.GetInputManager().KeyCallback(int(key), scancode, int(action), int(mods))
}

// SetActive pauses/resumes the timer manager.
func (w *WindowSystem) SetActive(active bool) {
	w.active = active

	if w.active == true {
		w.window.SetCursorPos(focusOutCursorX, focusOutCursorY)
		w.window.SetCursorPosCallback(w.MouseMoveCallback)
		core.GetTimerManager().Start()
	} else {
		focusOutCursorX, focusOutCursorY = w.window.GetCursorPos()
		w.window.SetCursorPosCallback(nil)
		core.GetTimerManager().Pause()
	}
}

// MouseButtonCallback passes mouse button events to the core.InputManager
func (w *WindowSystem) MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	core.GetInputManager().MouseButtonCallback(int(button), int(action), int(mods))
}

// SizeCallback passes window resize events to the core.InputManager
func (w *WindowSystem) SizeCallback(window *glfw.Window, width, height int) {
	w.cfg.Width, w.cfg.Height = width, height
}

// IconifyCallback passes window iconify events to the core.InputManager
func (w *WindowSystem) IconifyCallback(window *glfw.Window, iconified bool) {
	if iconified {
		w.SetActive(false)
	} else {
		w.SetActive(true)
	}
}

// FocusCallback passes focus events to the core.InputManager
func (w *WindowSystem) FocusCallback(window *glfw.Window, focused bool) {
	if focused {
		w.SetActive(true)
	} else {
		w.SetActive(false)
	}

	w.SetCursorVisible(w.cursorVisible)
}

// MouseMoveCallback passes mouse scroll events to the core.InputManager
func (w *WindowSystem) MouseMoveCallback(window *glfw.Window, x, y float64) {
	core.GetInputManager().MouseMoveCallback(x, y)
}

// MouseScrollCallback passes mouse scroll events to the core.InputManager
func (w *WindowSystem) MouseScrollCallback(window *glfw.Window, x, y float64) {
	if core.GetPlatform() == core.PlatformLinux || core.GetPlatform() == core.PlatformWindows {
		// just defaults, user may want to invert axes
		y = -y
	}
	core.GetInputManager().MouseScrollCallback(x, y)
}

// SetInputMap implements the core.WindowSystem interface by setting the core values
// of each input type to the internal ones.
func (w *WindowSystem) SetInputMap() {
	core.Joystick1 = int(glfw.Joystick1)
	core.Joystick2 = int(glfw.Joystick2)
	core.Joystick3 = int(glfw.Joystick3)
	core.Joystick4 = int(glfw.Joystick4)
	core.Joystick5 = int(glfw.Joystick5)
	core.Joystick6 = int(glfw.Joystick6)
	core.Joystick7 = int(glfw.Joystick7)
	core.Joystick8 = int(glfw.Joystick8)
	core.Joystick9 = int(glfw.Joystick9)
	core.Joystick10 = int(glfw.Joystick10)
	core.Joystick11 = int(glfw.Joystick11)
	core.Joystick12 = int(glfw.Joystick12)
	core.Joystick13 = int(glfw.Joystick13)
	core.Joystick14 = int(glfw.Joystick14)
	core.Joystick15 = int(glfw.Joystick15)
	core.Joystick16 = int(glfw.Joystick16)
	core.JoystickLast = int(glfw.JoystickLast)

	core.KeyUnknown = int(glfw.KeyUnknown)
	core.KeySpace = int(glfw.KeySpace)
	core.KeyApostrophe = int(glfw.KeyApostrophe)
	core.KeyComma = int(glfw.KeyComma)
	core.KeyMinus = int(glfw.KeyMinus)
	core.KeyPeriod = int(glfw.KeyPeriod)
	core.KeySlash = int(glfw.KeySlash)

	core.Key0 = int(glfw.Key0)
	core.Key1 = int(glfw.Key1)
	core.Key2 = int(glfw.Key2)
	core.Key3 = int(glfw.Key3)
	core.Key4 = int(glfw.Key4)
	core.Key5 = int(glfw.Key5)
	core.Key6 = int(glfw.Key6)
	core.Key7 = int(glfw.Key7)
	core.Key8 = int(glfw.Key8)
	core.Key9 = int(glfw.Key9)

	core.KeySemicolon = int(glfw.KeySemicolon)
	core.KeyEqual = int(glfw.KeyEqual)

	core.KeyA = int(glfw.KeyA)
	core.KeyB = int(glfw.KeyB)
	core.KeyC = int(glfw.KeyC)
	core.KeyD = int(glfw.KeyD)
	core.KeyE = int(glfw.KeyE)
	core.KeyF = int(glfw.KeyF)
	core.KeyG = int(glfw.KeyG)
	core.KeyH = int(glfw.KeyH)
	core.KeyI = int(glfw.KeyI)
	core.KeyJ = int(glfw.KeyJ)
	core.KeyK = int(glfw.KeyK)
	core.KeyL = int(glfw.KeyL)
	core.KeyM = int(glfw.KeyM)
	core.KeyN = int(glfw.KeyN)
	core.KeyO = int(glfw.KeyO)
	core.KeyP = int(glfw.KeyP)
	core.KeyQ = int(glfw.KeyQ)
	core.KeyR = int(glfw.KeyR)
	core.KeyS = int(glfw.KeyS)
	core.KeyT = int(glfw.KeyT)
	core.KeyU = int(glfw.KeyU)
	core.KeyV = int(glfw.KeyV)
	core.KeyW = int(glfw.KeyW)
	core.KeyX = int(glfw.KeyX)
	core.KeyY = int(glfw.KeyY)
	core.KeyZ = int(glfw.KeyZ)

	core.KeyLeftBracket = int(glfw.KeyLeftBracket)
	core.KeyBackslash = int(glfw.KeyBackslash)
	core.KeyRightBracket = int(glfw.KeyRightBracket)
	core.KeyGraveAccent = int(glfw.KeyGraveAccent)
	core.KeyWorld1 = int(glfw.KeyWorld1)
	core.KeyWorld2 = int(glfw.KeyWorld2)
	core.KeyEscape = int(glfw.KeyEscape)
	core.KeyEnter = int(glfw.KeyEnter)
	core.KeyTab = int(glfw.KeyTab)
	core.KeyBackspace = int(glfw.KeyBackspace)
	core.KeyInsert = int(glfw.KeyInsert)
	core.KeyDelete = int(glfw.KeyDelete)
	core.KeyRight = int(glfw.KeyRight)
	core.KeyLeft = int(glfw.KeyLeft)
	core.KeyDown = int(glfw.KeyDown)
	core.KeyUp = int(glfw.KeyUp)
	core.KeyPageUp = int(glfw.KeyPageUp)
	core.KeyPageDown = int(glfw.KeyPageDown)
	core.KeyHome = int(glfw.KeyHome)
	core.KeyEnd = int(glfw.KeyEnd)
	core.KeyCapsLock = int(glfw.KeyCapsLock)
	core.KeyScrollLock = int(glfw.KeyScrollLock)
	core.KeyNumLock = int(glfw.KeyNumLock)
	core.KeyPrintScreen = int(glfw.KeyPrintScreen)
	core.KeyPause = int(glfw.KeyPause)

	core.KeyF1 = int(glfw.KeyF1)
	core.KeyF2 = int(glfw.KeyF2)
	core.KeyF3 = int(glfw.KeyF3)
	core.KeyF4 = int(glfw.KeyF4)
	core.KeyF5 = int(glfw.KeyF5)
	core.KeyF6 = int(glfw.KeyF6)
	core.KeyF7 = int(glfw.KeyF7)
	core.KeyF8 = int(glfw.KeyF8)
	core.KeyF9 = int(glfw.KeyF9)
	core.KeyF10 = int(glfw.KeyF10)
	core.KeyF11 = int(glfw.KeyF11)
	core.KeyF12 = int(glfw.KeyF12)
	core.KeyF13 = int(glfw.KeyF13)
	core.KeyF14 = int(glfw.KeyF14)
	core.KeyF15 = int(glfw.KeyF15)
	core.KeyF16 = int(glfw.KeyF16)
	core.KeyF17 = int(glfw.KeyF17)
	core.KeyF18 = int(glfw.KeyF18)
	core.KeyF19 = int(glfw.KeyF19)
	core.KeyF20 = int(glfw.KeyF20)
	core.KeyF21 = int(glfw.KeyF21)
	core.KeyF22 = int(glfw.KeyF22)
	core.KeyF23 = int(glfw.KeyF23)
	core.KeyF24 = int(glfw.KeyF24)
	core.KeyF25 = int(glfw.KeyF25)
	core.KeyKP0 = int(glfw.KeyKP0)
	core.KeyKP1 = int(glfw.KeyKP1)
	core.KeyKP2 = int(glfw.KeyKP2)
	core.KeyKP3 = int(glfw.KeyKP3)
	core.KeyKP4 = int(glfw.KeyKP4)
	core.KeyKP5 = int(glfw.KeyKP5)
	core.KeyKP6 = int(glfw.KeyKP6)
	core.KeyKP7 = int(glfw.KeyKP7)
	core.KeyKP8 = int(glfw.KeyKP8)
	core.KeyKP9 = int(glfw.KeyKP9)

	core.KeyKPDecimal = int(glfw.KeyKPDecimal)
	core.KeyKPDivide = int(glfw.KeyKPDivide)
	core.KeyKPMultiply = int(glfw.KeyKPMultiply)
	core.KeyKPSubtract = int(glfw.KeyKPSubtract)
	core.KeyKPAdd = int(glfw.KeyKPAdd)
	core.KeyKPEnter = int(glfw.KeyKPEnter)
	core.KeyKPEqual = int(glfw.KeyKPEqual)
	core.KeyLeftShift = int(glfw.KeyLeftShift)
	core.KeyLeftControl = int(glfw.KeyLeftControl)
	core.KeyLeftAlt = int(glfw.KeyLeftAlt)
	core.KeyLeftSuper = int(glfw.KeyLeftSuper)
	core.KeyRightShift = int(glfw.KeyRightShift)
	core.KeyRightControl = int(glfw.KeyRightControl)
	core.KeyRightAlt = int(glfw.KeyRightAlt)
	core.KeyRightSuper = int(glfw.KeyRightSuper)
	core.KeyMenu = int(glfw.KeyMenu)
	core.KeyLast = int(glfw.KeyLast)

	core.MouseButton1 = int(glfw.MouseButton1)
	core.MouseButton2 = int(glfw.MouseButton2)
	core.MouseButton3 = int(glfw.MouseButton3)
	core.MouseButton4 = int(glfw.MouseButton4)
	core.MouseButton5 = int(glfw.MouseButton5)

	core.ActionPress = int(glfw.Press)
	core.ActionRelease = int(glfw.Release)
	core.ActionRepeat = int(glfw.Repeat)
}

// GetTime implements the core.WindowSystem interface
func (w *WindowSystem) GetTime() float64 {
	return glfw.GetTime()
}
