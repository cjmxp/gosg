package audiosystem

import (
	"github.com/fcvarela/gosg/core"
	"github.com/golang/glog"
)

// AudioSystem implements the core.AudioSystem interface by wrapping OpenAL
type AudioSystem struct {
	//Device  *al.Device
	//Context *al.Context
}

func init() {
	core.SetAudioSystem(New())
}

// New returns an initialized AudioSystem
func New() *AudioSystem {
	a := AudioSystem{}
	return &a
}

// Start implements the core.AudioSystem interface
func (a *AudioSystem) Start() {
	glog.Info("Starting")
}

// Step implements the core.AudioSystem interface
func (a *AudioSystem) Step() {

}

// Stop implements the core.AudioSystem interface
func (a *AudioSystem) Stop() {
	glog.Info("Stopping")
}
