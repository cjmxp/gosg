package core

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
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
	PrepareViewport(c *Camera)

	// NewMesh retuns a new mesh.
	NewMesh() Mesh

	// NewInstancedMesh returns a new instanced mesh from an existing mesh.
	NewInstancedMesh(Mesh) InstancedMesh

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

	// The following are performance measuring. They will be moved to a generalized CPU/GPU timer
	StartTimer()
	EndTimer()
	RenderTime() time.Duration
	CPUTime() time.Duration
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

// MeshBuckets takes a list of *Node containing Meshes and returns a map with
// each mesh as a key and a list of *Node containing it as values. This is useful
// to bucket Nodes containing the same mesh in order to implement dynamic instancing.
func MeshBuckets(nodes []*Node) map[Mesh][]*Node {
	buckets := make(map[Mesh][]*Node)
	for _, n := range nodes {
		if _, ok := buckets[n.Mesh()]; ok != true {
			buckets[n.Mesh()] = make([]*Node, 0)
		}
		buckets[n.Mesh()] = append(buckets[n.Mesh()], n)
	}

	// bucket per materials, less state changes
	matrixBuckets := make(map[Mesh][]float32)

	for mesh, bucketNodes := range buckets {
		mesh, ok := mesh.(InstancedMesh)
		if !ok {
			continue
		}

		// create array of floats from node transforms
		matrixBuckets[mesh] = make([]float32, 0)
		for _, n := range bucketNodes {
			transform64 := n.WorldTransform()
			transform32 := Mat4DoubleToFloat(transform64)
			matrixBuckets[mesh] = append(matrixBuckets[mesh], transform32[0:16]...)
		}
		mesh.SetInstanceCount(len(bucketNodes))
		mesh.SetModelMatrices(matrixBuckets[mesh])
	}

	return buckets
}

// RenderFn is the building block of rendering, it accepts a uniformbuffer, a camera and a list of nodes
// and renders them.
type RenderFn func(UniformBuffer, []*Node)

// RenderFnZPass sorts the nodes near-to-far, sets a non-color-write state with a special depth-only shader
func RenderFnZPass(ub UniformBuffer, nodes []*Node) {
	for _, node := range nodes {
		st := NewZPassState()
		node.mesh.Draw(ub, &st)
	}
}

// RenderFnZPassInstanced sorts the nodes near-to-far, sets a non-color-write state with a special depth-only shader
func RenderFnZPassInstanced(ub UniformBuffer, nodes []*Node) {
	st := NewInstancedZPassState()
	nodes[0].mesh.Draw(ub, &st)
}

// RenderFnOpaqueDiffusePass renders nodes with zfunc equal
func RenderFnOpaqueDiffusePass(ub UniformBuffer, nodes []*Node) {
	for _, node := range nodes {
		if node.state.Blend.Enabled {
			continue
		}
		st := *node.State()
		st.Depth.Mask = false
		st.Depth.Enabled = true
		st.Depth.Func = DepthEqual
		node.mesh.Draw(ub, &st)
	}
}

// RenderFnTransparentDiffusePass renders transparent nodes with back-to-front with blending
func RenderFnTransparentDiffusePass(ub UniformBuffer, nodes []*Node) {
	for _, node := range nodes {
		if !node.state.Blend.Enabled {
			continue
		}

		st := *node.State()
		st.Depth.Enabled = true
		st.Depth.Mask = true
		st.Depth.Func = DepthLessEqual
		st.Blend.SrcMode = BlendSrcAlpha
		st.Blend.DstMode = BlendOneMinusSrcAlpha
		node.mesh.Draw(ub, &st)
	}
}

// DefaultRenderTechnique does z pre-pass, diffuse pass, transparency pass
func DefaultRenderTechnique(ub UniformBuffer, nodes []*Node) {
	// bucket per materials, less state changes
	meshBuckets := MeshBuckets(nodes)

	// zpass
	for mesh, bucketNodes := range meshBuckets {
		if _, ok := mesh.(InstancedMesh); ok {
			RenderFnZPassInstanced(ub, []*Node{bucketNodes[0]})
		} else {
			RenderFnZPass(ub, bucketNodes)
		}
	}

	// diffuse opaque pass
	for mesh, bucketNodes := range meshBuckets {
		if _, ok := mesh.(InstancedMesh); ok {
			//RenderFnOpaqueDiffusePass([]*Node{bucketNodes[0]})
			RenderFnOpaqueDiffusePass(ub, []*Node{bucketNodes[0]})
		} else {
			RenderFnOpaqueDiffusePass(ub, bucketNodes)
		}
	}

	return
	//	diffuse transparent pass
	for mesh, bucketNodes := range meshBuckets {
		if _, ok := mesh.(InstancedMesh); ok {
			RenderFnTransparentDiffusePass(ub, []*Node{bucketNodes[0]})
		} else {
			RenderFnTransparentDiffusePass(ub, bucketNodes)
		}
	}

	var aabbNodes []*Node

	for _, node := range nodes {
		if node.state.AABB == true {
			aabbNodes = append(aabbNodes, node)
		}
	}

	AABBRenderTechnique(ub, aabbNodes)
}

// IMGUIRenderTechnique renders IMGUI UI nodes.
func IMGUIRenderTechnique(ub UniformBuffer, nodes []*Node) {
	for _, node := range nodes {
		node.mesh.Draw(ub, node.State())
	}
}

// AABBRenderTechnique renders AABBs and OBB. OBBs are rendered red and AABBs white.
func AABBRenderTechnique(ub UniformBuffer, nodes []*Node) {
	// create transform list from nodes, like meshbuckets
	if boundsMesh == nil {
		boundsMesh = NewAABBMesh()
	}

	for _, node := range nodes {
		// nodespace bounds: red
		st := NewAABBState()
		center := node.Bounds().Center()
		size := node.mesh.Bounds().Size()
		st.Uniform("flatColor").Set(mgl64.Vec4{1.0, 0.0, 0.0, 1.0})
		st.Uniform("mMatrix").Set(
			node.WorldTransform().Mul4(
				mgl64.Translate3D(center[0], center[1], center[2]).Mul4(
					mgl64.Scale3D(size[0], size[1], size[2]))))
		boundsMesh.Draw(ub, &st)

		// world bounds: white
		st = NewAABBState()
		center = node.worldBounds.Center()
		size = node.worldBounds.Size()
		// world bounds, no need for node transforms
		st.Uniform("flatColor").Set(mgl64.Vec4{1.0, 1.0, 1.0, 1.0})
		st.Uniform("mMatrix").Set(
			mgl64.Translate3D(center[0], center[1], center[2]).Mul4(
				mgl64.Scale3D(size[0], size[1], size[2])))
		boundsMesh.Draw(ub, &st)
	}
}
