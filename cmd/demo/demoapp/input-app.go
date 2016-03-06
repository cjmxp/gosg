package demoapp

import "github.com/fcvarela/gosg/core"

// ApplicationInputComponent implements InputComponent
type applicationInputComponent struct{}

// Run checks the InputSystem for actionable state and returns commands
func (ic *applicationInputComponent) Run() (commands []core.ClientApplicationCommand) {
	// check for quit key, append to command list
	state := *core.GetInputManager().State()
	if state.Keys.Active[core.KeyEscape] == true {
		commands = append(commands, new(clientApplicationQuitCommand))
	}

	if state.Keys.Active[core.Key1] == true {
		commands = append(commands, new(clientApplicationShowDemo1Command))
	}

	if state.Keys.Active[core.Key2] == true {
		commands = append(commands, new(clientApplicationShowDemo2Command))
	}

	// key-up, after down
	if state.Keys.Released[core.KeyE] == true {
		commands = append(commands, new(clientApplicationToggleDebugMenuCommand))
	}

	return commands
}
