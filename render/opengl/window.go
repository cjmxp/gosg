package opengl

import (
	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang/glog"
)

func (r *RenderSystem) MakeWindow(cfg core.WindowConfig) *glfw.Window {
	// create a window
	glfw.WindowHint(glfw.Decorated, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 0)

	var err error
	var window *glfw.Window
	if cfg.Fullscreen {
		window, err = glfw.CreateWindow(cfg.Width, cfg.Height, cfg.Name, cfg.Monitor, nil)
	} else {
		window, err = glfw.CreateWindow(cfg.Width, cfg.Height, cfg.Name, nil, nil)
	}
	if err != nil {
		glog.Fatal(err)
	}

	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	glog.Info("Checking GL Init status")
	glog.Info("OpenGL version: ", gl.GoStr(gl.GetString(gl.VERSION)))
	glog.Info("OpenGL renderer: ", gl.GoStr(gl.GetString(gl.RENDERER)))

	return window
}
