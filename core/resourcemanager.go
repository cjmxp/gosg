package core

import (
	"encoding/json"
	"log"

	"github.com/golang/glog"
)

// ResourceSystem is an interface which wraps all resource management logic.
type ResourceSystem interface {
	// Start is called at application startup time. Implementations requiring init may do so here.
	Start()

	// Stop is called at application shutdown time. Implementations requiring cleanup/saving may do so here.
	Stop()

	// Model returns a byte array representing a model.
	Model(string) []byte

	// Texture returns a byte array representing a texture.
	Texture(string) []byte

	// Program returns a byte array representing a program.
	Program(string) []byte

	// ProgramData returns a byte array representing program data.
	ProgramData(string) []byte
}

// ResourceManager wraps a resourcesystem and contains configuration about the location of each resource type.
type ResourceManager struct {
	system          ResourceSystem
	programs        map[string]Program
	models          map[string]*Node
	instancedModels map[string]*Node
	textures        map[string]Texture
}

var (
	resourceManager *ResourceManager
)

func init() {
	resourceManager = &ResourceManager{
		programs:        make(map[string]Program),
		models:          make(map[string]*Node),
		instancedModels: make(map[string]*Node),
		textures:        make(map[string]Texture),
	}
}

// GetResourceManager returns the resource manager. Used by client applications to load assets.
func GetResourceManager() *ResourceManager {
	return resourceManager
}

func (r *ResourceManager) start() {
	glog.Info("Starting")
	r.system.Start()
	glog.Info("Started")
}

func (r *ResourceManager) stop() {
	glog.Info("Stopping...")
	r.system.Stop()
	glog.Info("Stopped")
}

// SetSystem sets the resource manager's resource system
func (r *ResourceManager) SetSystem(s ResourceSystem) {
	if r.system != nil {
		log.Fatal("Can't replace previously registered resource system. Please make sure you're not importing twice")
	}
	r.system = s
}

// Model returns a scenegraph node with a subtree of nodes containing meshes which represent a complex model.
func (r *ResourceManager) Model(name string) *Node {
	if r.models[name] == nil {
		resource := r.system.Model(name)
		r.models[name] = LoadModel(name, resource, false)
	}

	return r.models[name].Copy()
}

// InstancedModel returns a Model from Model(), configured with instanced meshes.
func (r *ResourceManager) InstancedModel(name string) *Node {
	if r.models[name] == nil {
		resource := r.system.Model(name)
		r.models[name] = LoadModel(name, resource, true)
	}

	return r.models[name].Copy()
}

// Program returns a GPU program.
func (r *ResourceManager) Program(name string) Program {
	if r.programs[name] == nil {
		resource := r.system.Program(name)

		// we know programs are json files with key=>subresourcename entries
		var subresourcenames map[string]string
		if err := json.Unmarshal(resource, &subresourcenames); err != nil {
			glog.Fatal("Cannot parse program: ", err)
		}

		var subresources = make(map[string][]byte)
		for k, v := range subresourcenames {
			subresources[k] = r.system.ProgramData(v)
		}

		r.programs[name] = renderSystem.NewProgram(name, subresources)
	}
	return r.programs[name]
}

// Texture returns a Texture.
func (r *ResourceManager) Texture(name string) Texture {
	if r.textures[name] == nil {
		resource := r.system.Texture(name)
		r.textures[name] = renderSystem.NewTexture(resource)
	}
	return r.textures[name]
}
