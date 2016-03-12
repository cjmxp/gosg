package core

import (
	"log"

	"github.com/fcvarela/gosg/protos"
	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
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

	// Material returns a byte array representing a material
	Material(string) []byte

	// ProgramData returns a byte array representing program data.
	ProgramData(string) []byte
}

// ResourceManager wraps a resourcesystem and contains configuration about the location of each resource type.
type ResourceManager struct {
	system          ResourceSystem
	programs        map[string]Program
	materials       map[string]*protos.Material
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
		materials:       make(map[string]*protos.Material),
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
		r.models[name] = LoadModel(name, resource)
	}
	return r.models[name].Copy()
}

// Program returns a GPU program.
func (r *ResourceManager) Program(name string) Program {
	if r.programs[name] == nil {
		resource := r.system.Program(name)
		r.programs[name] = renderSystem.NewProgram(name, resource)
	}
	return r.programs[name]
}

// Material returns a material.
func (r *ResourceManager) Material(name string) *protos.Material {
	if r.materials[name] == nil {
		resource := r.system.Material(name)
		var material protos.Material
		if err := jsonpb.UnmarshalString(string(resource), &material); err != nil {
			glog.Fatal("Cannot unmarshal material: ", err)
		}
		material.Name = name
		r.materials[name] = &material
	}
	return r.materials[name]
}

// ProgramData returns source file contents for a given program or subprogram
// This is meant to be used by rendersystem implementations to load subresources for a program spec
func (r *ResourceManager) ProgramData(name string) []byte {
	return r.system.ProgramData(name)
}

// Texture returns a Texture.
func (r *ResourceManager) Texture(name string) Texture {
	if r.textures[name] == nil {
		resource := r.system.Texture(name)
		r.textures[name] = renderSystem.NewTexture(resource)
	}
	return r.textures[name]
}
