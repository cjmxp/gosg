package demoapp

import "github.com/fcvarela/gosg/core"

type clientApplicationQuitCommand struct{}

func (c *clientApplicationQuitCommand) Run(ac core.ClientApplication) {
	ac.Stop()
}

type clientApplicationShowDemo1Command struct{}

func (c *clientApplicationShowDemo1Command) Run(ac core.ClientApplication) {
	sm := core.GetSceneManager()

	if sm.FrontScene().Name() == "Demo1" {
		return
	}
	sm.PopScene()
	sm.PushScene(makeDemo1Scene())
}

type clientApplicationShowDemo2Command struct{}

func (c *clientApplicationShowDemo2Command) Run(ac core.ClientApplication) {
	sm := core.GetSceneManager()

	if sm.FrontScene().Name() == "Demo2" {
		return
	}
	sm.PopScene()
	sm.PushScene(makeDemo2Scene())
}

type clientApplicationToggleDebugMenuCommand struct{}

func (c *clientApplicationToggleDebugMenuCommand) Run(ac core.ClientApplication) {
	sm := core.GetSceneManager()
	frontScene := sm.FrontScene()

	switch frontScene.Name() {
	case "debugMenu":
		sm.PopScene()
	case "Demo1":
		debugMenu := new(demo1DebugMenuInputComponent)
		debugMenu.shadowTexture = getDemo1SceneShadowTexture(frontScene)
		sm.PushScene(core.NewIMGUIScene("debugMenu", debugMenu))
	}
}
