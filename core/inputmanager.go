package core

// declares the datatypes, input implementations will map to
// their internal spec

// joysticks
var (
	Joystick1    int
	Joystick2    int
	Joystick3    int
	Joystick4    int
	Joystick5    int
	Joystick6    int
	Joystick7    int
	Joystick8    int
	Joystick9    int
	Joystick10   int
	Joystick11   int
	Joystick12   int
	Joystick13   int
	Joystick14   int
	Joystick15   int
	Joystick16   int
	JoystickLast int
)

// keys
var (
	KeyUnknown      int
	KeySpace        int
	KeyApostrophe   int
	KeyComma        int
	KeyMinus        int
	KeyPeriod       int
	KeySlash        int
	Key0            int
	Key1            int
	Key2            int
	Key3            int
	Key4            int
	Key5            int
	Key6            int
	Key7            int
	Key8            int
	Key9            int
	KeySemicolon    int
	KeyEqual        int
	KeyA            int
	KeyB            int
	KeyC            int
	KeyD            int
	KeyE            int
	KeyF            int
	KeyG            int
	KeyH            int
	KeyI            int
	KeyJ            int
	KeyK            int
	KeyL            int
	KeyM            int
	KeyN            int
	KeyO            int
	KeyP            int
	KeyQ            int
	KeyR            int
	KeyS            int
	KeyT            int
	KeyU            int
	KeyV            int
	KeyW            int
	KeyX            int
	KeyY            int
	KeyZ            int
	KeyLeftBracket  int
	KeyBackslash    int
	KeyRightBracket int
	KeyGraveAccent  int
	KeyWorld1       int
	KeyWorld2       int
	KeyEscape       int
	KeyEnter        int
	KeyTab          int
	KeyBackspace    int
	KeyInsert       int
	KeyDelete       int
	KeyRight        int
	KeyLeft         int
	KeyDown         int
	KeyUp           int
	KeyPageUp       int
	KeyPageDown     int
	KeyHome         int
	KeyEnd          int
	KeyCapsLock     int
	KeyScrollLock   int
	KeyNumLock      int
	KeyPrintScreen  int
	KeyPause        int
	KeyF1           int
	KeyF2           int
	KeyF3           int
	KeyF4           int
	KeyF5           int
	KeyF6           int
	KeyF7           int
	KeyF8           int
	KeyF9           int
	KeyF10          int
	KeyF11          int
	KeyF12          int
	KeyF13          int
	KeyF14          int
	KeyF15          int
	KeyF16          int
	KeyF17          int
	KeyF18          int
	KeyF19          int
	KeyF20          int
	KeyF21          int
	KeyF22          int
	KeyF23          int
	KeyF24          int
	KeyF25          int
	KeyKP0          int
	KeyKP1          int
	KeyKP2          int
	KeyKP3          int
	KeyKP4          int
	KeyKP5          int
	KeyKP6          int
	KeyKP7          int
	KeyKP8          int
	KeyKP9          int
	KeyKPDecimal    int
	KeyKPDivide     int
	KeyKPMultiply   int
	KeyKPSubtract   int
	KeyKPAdd        int
	KeyKPEnter      int
	KeyKPEqual      int
	KeyLeftShift    int
	KeyLeftControl  int
	KeyLeftAlt      int
	KeyLeftSuper    int
	KeyRightShift   int
	KeyRightControl int
	KeyRightAlt     int
	KeyRightSuper   int
	KeyMenu         int
	KeyLast         int
)

// actions
var (
	ActionPress   int
	ActionRelease int
	ActionRepeat  int
)

// mouse buttons
var (
	MouseButton1 int
	MouseButton2 int
	MouseButton3 int
	MouseButton4 int
	MouseButton5 int
)

// MousePositionState holds the mouse position information.
type MousePositionState struct {
	Valid bool
	X     float64
	Y     float64
	DistX float64
	DistY float64
}

// MouseButtonState holds the mouse button state.
type MouseButtonState struct {
	Valid  bool
	Active map[int]bool
	Action int
}

// MouseScrollState holds the mouse scroll state.
type MouseScrollState struct {
	Valid bool
	X     float64
	Y     float64
}

// MouseState holds mouse input state.
type MouseState struct {
	Valid    bool
	Position MousePositionState
	Scroll   MouseScrollState
	Buttons  MouseButtonState
}

// KeyState holds key input state.
type KeyState struct {
	Valid    bool
	Mods     map[int]bool
	Active   map[int]bool
	Released map[int]bool
}

// InputState wraps mouse and keys input state.
type InputState struct {
	Mouse MouseState
	Keys  KeyState
}

// SetMouseValid sets the mouse state as valid. It will not be processed unless this is set.
func (i *InputState) SetMouseValid(valid bool) {
	i.Mouse.Valid = valid
}

// SetKeysValid sets the key state as valid. It will not be processed unless this is set.
func (i *InputState) SetKeysValid(valid bool) {
	i.Keys.Valid = valid
}

// InputManager wraps global input state. WindowSystem implementations use the manager to expose
// input state to the system.
type InputManager struct {
	state InputState
}

// InputComponent is an interface which returns NodeCommands from nodes. Each node may have its own
// input component which checks the manager for input and determines what commands should be output.
type InputComponent interface {
	// Run returns commands from a given node to itself.
	Run(node *Node) []NodeCommand
}

var (
	inputManager *InputManager
)

func init() {
	inputManager = &InputManager{}
	inputManager.state.Keys.Active = make(map[int]bool)
	inputManager.state.Keys.Released = make(map[int]bool)
	inputManager.state.Mouse.Buttons.Active = make(map[int]bool)
}

// GetInputManager returns the manager.
func GetInputManager() *InputManager {
	return inputManager
}

// State returns the manager's input state.
func (i *InputManager) State() *InputState {
	return &i.state
}

// Reset resets all input state and marks substates as invalid.
func (i *InputManager) Reset() {
	for j := range i.state.Keys.Released {
		i.state.Keys.Released[j] = false
	}

	i.state.Mouse.Valid = false
	i.state.Mouse.Position.Valid = false
	i.state.Mouse.Buttons.Valid = false
	i.state.Mouse.Position.DistX = 0.0
	i.state.Mouse.Position.DistY = 0.0
	i.state.Mouse.Scroll = MouseScrollState{false, 0.0, 0.0}
}

// KeyCallback is called by windowsystems to register key events.
func (i *InputManager) KeyCallback(key int, scancode int, action int, mods int) {
	i.state.Keys.Valid = true

	if action == ActionPress {
		i.state.Keys.Active[key] = true
		i.state.Keys.Released[key] = false
	} else if action == ActionRepeat {
		i.state.Keys.Active[key] = true
	} else if action == ActionRelease {
		i.state.Keys.Active[key] = false
		i.state.Keys.Released[key] = true
	}
}

// MouseButtonCallback is called by windowsystems to register mouse button events.
func (i *InputManager) MouseButtonCallback(button int, action int, mods int) {
	i.state.Mouse.Buttons.Valid = true

	if action == ActionPress {
		i.state.Mouse.Buttons.Active[button] = true
	} else {
		i.state.Mouse.Buttons.Active[button] = false
	}
}

// MouseScrollCallback is called by windowsystems to register mouse scroll events.
func (i *InputManager) MouseScrollCallback(x, y float64) {
	i.state.Mouse.Valid = true
	i.state.Mouse.Scroll.Valid = true

	i.state.Mouse.Scroll.X = x
	i.state.Mouse.Scroll.Y = y
}

// MouseMoveCallback is called by windowsystems to register mouse move events.
func (i *InputManager) MouseMoveCallback(x, y float64) {
	i.state.Mouse.Valid = true
	i.state.Mouse.Position.Valid = true

	i.state.Mouse.Position.DistX = x - i.state.Mouse.Position.X
	i.state.Mouse.Position.DistY = y - i.state.Mouse.Position.Y

	i.state.Mouse.Position.X = x
	i.state.Mouse.Position.Y = y
}
