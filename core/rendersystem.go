package core

import (
	"fmt"
	"sort"
	"time"

	"github.com/fcvarela/gosg/protos"
	"github.com/golang/glog"
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

	// The following are performance measuring. They will be moved to a generalized CPU/GPU timer
	StartTimer()
	EndTimer()
	RenderTime() time.Duration
	CPUTime() time.Duration
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
	Nodes    []*Node
	Name     string
	Material *protos.Material
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

func debugNodeMaterials(nodes []*Node) {
	for _, n := range nodes {
		glog.Infof("%s: %#v", n.material, n.materialData)
	}
}

// MaterialBuckets takes a list of *Node and returns a map with each materialName as a key and a list of *Node
// using it as values.
func MaterialBuckets(nodes []*Node) map[*protos.Material][]*Node {
	// sort by material, then bucket
	sort.Sort(NodesByMaterial(nodes))

	buckets := make(map[*protos.Material][]*Node)
	for _, n := range nodes {
		if _, ok := buckets[n.material]; ok != true {
			buckets[n.material] = make([]*Node, 0)
		}
		buckets[n.material] = append(buckets[n.material], n)
	}
	return buckets
}

// DefaultRenderTechnique does z pre-pass, diffuse pass, transparency pass
func DefaultRenderTechnique(camera *Camera, nodes []*Node) (out RenderStage) {
	out.Name = fmt.Sprintf("%s-DefaultRenderTechnique", camera.name)
	out.Camera = camera

	// create a depth prepass, single state, program and all nodes
	out.Passes = append(out.Passes, RenderPass{
		Material: resourceManager.Material("zpass"),
		Name:     "DepthPrePass",
		Nodes:    nodes,
	})

	// get per-material buckets
	materialBuckets := MaterialBuckets(nodes)

	// create pass per bucket
	// fixme: make sure transparent materials are always rendered last
	for material, nodeBucket := range materialBuckets {
		out.Passes = append(out.Passes, RenderPass{
			Material: material,
			Name:     "Diffuse",
			Nodes:    nodeBucket,
		})
	}

	return
}

// IMGUIRenderTechnique renders IMGUI UI nodes.
//func IMGUIRenderTechnique(c *Camera, nodes []*Node) (out RenderStage) {
//for _, node := range nodes {
//	node.mesh.Draw(ub, node.State())
//}
//}

// AABBRenderTechnique renders AABBs and OBB. OBBs are rendered red and AABBs white.
//func AABBRenderTechnique(ub UniformBuffer, nodes []*Node) {
//	// create transform list from nodes, like meshbuckets
//	if boundsMesh == nil {
//		boundsMesh = NewAABBMesh()
//	}
//
//	for _, node := range nodes {
//		// nodespace bounds: red
//		st := NewAABBState()
//		center := node.Bounds().Center()
//		size := node.mesh.Bounds().Size()
//		st.Uniform("flatColor").Set(mgl64.Vec4{1.0, 0.0, 0.0, 1.0})
//		st.Uniform("mMatrix").Set(
//			node.WorldTransform().Mul4(
//				mgl64.Translate3D(center[0], center[1], center[2]).Mul4(
//					mgl64.Scale3D(size[0], size[1], size[2]))))
//		boundsMesh.Draw(ub, &st)
//
//		// world bounds: white
//		st = NewAABBState()
//		center = node.worldBounds.Center()
//		size = node.worldBounds.Size()
//		// world bounds, no need for node transforms
//		st.Uniform("flatColor").Set(mgl64.Vec4{1.0, 1.0, 1.0, 1.0})
//		st.Uniform("mMatrix").Set(
//			mgl64.Translate3D(center[0], center[1], center[2]).Mul4(
//				mgl64.Scale3D(size[0], size[1], size[2])))
//		boundsMesh.Draw(ub, &st)
//	}
//}
