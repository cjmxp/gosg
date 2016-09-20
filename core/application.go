package core

import "github.com/go-gl/glfw/v3.2/glfw"

// ClientApplication is the client app, provided by the user
type ClientApplication interface {
	InputComponent() ClientApplicationInputComponent
	Stop()
	Done() bool
}

// ClientApplicationCommand is a single runnable command
type ClientApplicationCommand interface {
	Run(ac ClientApplication)
}

// ClientApplicationInputComponent checks the input system state
// and triggers commands
type ClientApplicationInputComponent interface {
	Run() []ClientApplicationCommand
}

// Application is the top-level runnable gosg App
type Application struct {
	client ClientApplication
}

// Start starts the application runloop by calling all systems/managers Start methods,
// and calling the ClientApp constructor. On runloop termination, the Stop methods are
// called in reverse order.
func (app *Application) Start(acConstructor func() ClientApplication) {
	windowManager.Start()

	renderSystem.Start()
	audioSystem.Start()
	physicsSystem.Start()
	imguiSystem.Start()

	resourceManager.start()

	// create the client app
	app.client = acConstructor()

	// step the windowsystem to force swap buffers before starting loop
	windowManager.window.SwapBuffers()

	// start main loop, all systems go
	app.runLoop()

	// done
	resourceManager.stop()

	imguiSystem.Stop()
	physicsSystem.Stop()
	audioSystem.Stop()
	renderSystem.Stop()

	glfw.Terminate()
}

func (app *Application) runLoop() {
	var dt = 1.0 / 60.0
	var start = timerManager.GetTime()
	var end = 0.0

	for !app.client.Done() && !windowManager.ShouldClose() {
		// run subsystem updates if not paused
		app.update(dt)

		// compute time delta
		end = timerManager.GetTime()
		dt = end - start

		// safeguard for extreme deltas (breakpoints, suspends)
		if dt > 1.0/10.0 {
			dt = 1.0 / 60.0
		}

		timerManager.SetDt(dt)

		// rotate time
		start = end
	}
}

func (app *Application) update(dt float64) {
	// update client app
	acCommands := app.client.InputComponent().Run()
	for _, command := range acCommands {
		command.Run(app.client)
	}

	// call game object updates
	sceneManager.update(dt)

	// run the culler
	sceneManager.cull()

	// draw the scenes
	sceneManager.draw()

	// play audio
	audioSystem.Step()

	// swap context buffers and poll for input
	windowManager.window.SwapBuffers()
	inputManager.reset()
	glfw.PollEvents()
	glfw.PollEvents()
	glfw.PollEvents()
}
