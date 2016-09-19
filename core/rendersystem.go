package core

import (
	"fmt"

	"github.com/fcvarela/gosg/protos"
)

// RenderSystem is an interface which wraps all logic related to rendering and memory management of
// GPU buffers.
type RenderSystem interface {
	// Start is called at application startup. This is where implementations will want to initialize their provided
	// rendering libraries.
	Start()

	// Stop is called at application shutdown. Implementations which require cleanup may do so here.
	Stop()

	// PrepareViewport is an implementation specific function which prepares a camera's framebuffers and viewport for
	// rendering according to the camera's policy.
	PrepareRenderTarget(c *Camera)

	// NewMesh retuns a new mesh.
	NewMesh() Mesh

	// NewIMGUIMesh returns a new IMGUI mesh.
	NewIMGUIMesh() IMGUIMesh

	// ProgramExtension exposes the resource extension of program definitions for the implementation.
	ProgramExtension() string

	// NewProgram creates a new program from a list of subprogram source files.
	NewProgram(name string, data []byte) Program

	// NewTexture creates a new texture from a byte buffer of image data
	NewTexture(r []byte) Texture

	// NewUniform creates a new empty uniform
	NewUniform() Uniform

	// NewUniformBuffer creates a new empty uniform buffer
	NewUniformBuffer() UniformBuffer

	// NewRawTexture creates an unverified/unvalidated texture from a width, height and byte buffer.
	NewRawTexture(width, height int, payload []byte) Texture

	// NewRenderTarget returns a render target of width, height, depth layer count and color layer count.
	// These are implementation specific but will generally be framebuffer attachments.
	NewRenderTarget(width uint32, height uint32, depthLayers uint8, layers uint8) RenderTarget

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
