package core

import (
	"fmt"

	"github.com/fcvarela/gosg/protos"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// RenderSystem is an interface which wraps all logic related to rendering and memory management of
// GPU buffers.
type RenderSystem interface {
	// Start is called at application startup. This is where implementations will want to initialize their provided
	// rendering libraries.
	Start()

	// Stop is called at application shutdown. Implementations which require cleanup may do so here.
	Stop()

	// PrepareWindow initializes a new window
	MakeWindow(cfg WindowConfig) *glfw.Window

	// NewMesh retuns a new mesh.
	NewMesh() Mesh

	// NewIMGUIMesh returns a new IMGUI mesh.
	NewIMGUIMesh() IMGUIMesh

	// ProgramExtension exposes the resource extension of program definitions for the implementation.
	ProgramExtension() string

	// NewProgram creates a new program from a list of subprogram source files.
	NewProgram(name string, data []byte) Program

	// NewTexture creates a new texture from a byte buffer containing an image file, not raw bitmap.
	// This always generates RGBA, unsigned byte, power of two and will generate mipmaps
	// levels from smallest dimension, ie: 2048x1024 = 10 mipmap levels; log2(1024)
	// It also defaults to ClampEdge and mipmapped filtering.
	NewTextureFromImageData(r []byte, d TextureDescriptor) Texture

	// NewUniform creates a new empty uniform
	NewUniform() Uniform

	// NewUniformBuffer creates a new empty uniform buffer
	NewUniformBuffer() UniformBuffer

	// NewRawTexture creates a new texture and allocates storage for it
	NewTexture(descriptor TextureDescriptor, data []byte) Texture

	// NewFramebuffer returns a newly created framebuffer
	NewFramebuffer() Framebuffer

	// ExecuteRenderPlan issues the actual drawing commands to the 3D api
	ExecuteRenderPlan(p RenderPlan)

	// RenderLog returns a log of the render plan
	RenderLog() string
}

// NodeFilter is a mask used by renderpasses to pick which nodes to render from a list
type NodeFilter int

const (
	NodeFilterOpaque NodeFilter = 1 << iota
	NodeFilterTransparent
	NodeFilterShadowCaster
	NodeFilterBloomCaster
)

// RenderPass represents one operation of rendering a set of nodes with a given rasterstate and a program
type RenderPass struct {
	Nodes []*Node
	Name  string
	State *protos.State
}

// RenderStage represents a group of renderpasses from a single camera and list of nodes
type RenderStage struct {
	Name   string
	Camera *Camera
	Passes []RenderPass
}

type RenderPlan struct {
	Stages []RenderStage
}

var (
	renderSystem RenderSystem

	// fixme, this is a hack, don't know where to keep this
	boundsMesh Mesh
)

// SetRenderSystem is meant to be called from RenderSystem implementations on their init method
// to cause the side-effect of setting the core RenderSystem to their provided one.
func SetRenderSystem(rs RenderSystem) {
	renderSystem = rs
}

// GetRenderSystem returns the renderSystem, thereby exposing it to any package importing core.
func GetRenderSystem() RenderSystem {
	return renderSystem
}

// DefaultRenderTechnique does z pre-pass, diffuse pass, transparency pass
func DefaultRenderTechnique(camera *Camera, materialBuckets map[*protos.State][]*Node) RenderStage {
	var out RenderStage
	out.Name = fmt.Sprintf("%s-DefaultRenderTechnique", camera.name)
	out.Camera = camera

	// create a depth prepass, single state, program and all opaque nodes
	var zPrepass = RenderPass{
		State: resourceManager.State("zpass"),
		Name:  "DepthPrePass",
		Nodes: []*Node{},
	}

	var opaquePasses = make([]RenderPass, 0)
	var transparentPasses = make([]RenderPass, 0)

	// create pass per bucket, opaque is default
	for material, nodeBucket := range materialBuckets {
		if material.Blending == true {
			transparentPasses = append(transparentPasses, RenderPass{
				State: material,
				Name:  "Diffuse-Transparent",
				Nodes: nodeBucket,
			})
			continue
		}

		opaquePasses = append(opaquePasses, RenderPass{
			State: material,
			Name:  "Diffuse",
			Nodes: nodeBucket,
		})

		// append opaque nodes to z prepass
		zPrepass.Nodes = append(zPrepass.Nodes, nodeBucket...)
	}

	out.Passes = append(out.Passes, zPrepass)
	out.Passes = append(out.Passes, opaquePasses...)
	out.Passes = append(out.Passes, transparentPasses...)

	return out
}
