package demoapp

import (
	"fmt"

	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/mathgl/mgl32"
)

type demo1DebugMenuInputComponent struct {
	shadowTexture core.Texture
}

func (u *demo1DebugMenuInputComponent) Run(n *core.Node) []core.NodeCommand {
	timerManager := core.GetTimerManager()
	imguiSystem := core.GetIMGUISystem()

	imguiSystem.StartFrame(timerManager.Dt())

	// Metrics window
	imguiSystem.SetNextWindowPos(mgl32.Vec2{0.0, 0.0})
	imguiSystem.SetNextWindowSize(mgl32.Vec2{320.0, core.GetWindowSystem().WindowSize()[1]})
	if imguiSystem.Begin("Inspector", core.WindowFlagsNoCollapse|core.WindowFlagsNoResize|core.WindowFlagsNoMove) {
		if imguiSystem.CollapsingHeader("Frame Times") {
			frameHistogram := timerManager.Histogram()
			fpsLabel := fmt.Sprintf("%.0f FPS", timerManager.AvgFPS())
			imguiSystem.PlotHistogram(fpsLabel, frameHistogram.Values, 0.0, frameHistogram.Max, mgl32.Vec2{0.0, 60.0})
		}

		if imguiSystem.CollapsingHeader("ShadowTexture") {
			imguiSystem.Image(u.shadowTexture, mgl32.Vec2{256.0, 256.0})
		}
	}

	imguiSystem.End()

	// Done
	imguiSystem.EndFrame()
	return nil
}
