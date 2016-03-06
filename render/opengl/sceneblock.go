package opengl

import (
	"unsafe"

	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type innerBlock struct {
	viewMatrix           mgl32.Mat4
	projectionMatrix     mgl32.Mat4
	viewProjectionMatrix mgl32.Mat4
	lightCount           mgl32.Vec4
	lightBlocks          [16]core.LightBlock
}

// SceneBlock implements the core.SceneBlock interface
type SceneBlock struct {
	id     uint32
	dirty  bool
	lights []*core.Light
	data   innerBlock
}

const (
	// 16 lights * 16 floats per light + 4 floats lightcount, all mult4 (sizeof float)
	sceneBlockLen = (3*16 + 16*16 + 4) * 4
)

// NewSceneBlock implements the core.RenderSystem interface
func (r *RenderSystem) NewSceneBlock() core.SceneBlock {
	sb := &SceneBlock{}
	gl.GenBuffers(1, &sb.id)
	return sb
}

// SetMatricesFromCamera implements the core.SceneBlock interface
func (sb *SceneBlock) SetMatricesFromCamera(c *core.Camera) {
	var vMatrix = c.ViewMatrix()
	var pMatrix = c.ProjectionMatrix()
	var vpMatrix = pMatrix.Mul4(vMatrix)

	sb.data.viewMatrix = core.Mat4DoubleToFloat(vMatrix)
	sb.data.projectionMatrix = core.Mat4DoubleToFloat(pMatrix)
	sb.data.viewProjectionMatrix = core.Mat4DoubleToFloat(vpMatrix)

	sb.dirty = true
}

// Lights implements the core.SceneBlock interface
func (sb *SceneBlock) Lights() []*core.Light {
	return sb.lights
}

// SetLights implements the core.SceneBlock interface
func (sb *SceneBlock) SetLights(l []*core.Light) {
	sb.data.lightCount = mgl32.Vec4{0.0, 0.0, 0.0, 0.0}
	sb.lights = l

	// copy lights
	for i := range sb.lights {
		if i > 15 {
			break
		}
		sb.data.lightCount = sb.data.lightCount.Add(mgl32.Vec4{1.0, 1.0, 1.0, 1.0})
		sb.data.lightBlocks[i] = l[i].Block
	}

	sb.dirty = true
}

func bindSceneBlock(sb *SceneBlock) {
	if sb.dirty {
		data := sb.data
		gl.BindBuffer(gl.UNIFORM_BUFFER, sb.id)
		gl.BufferData(gl.UNIFORM_BUFFER, sceneBlockLen, unsafe.Pointer(&data), gl.DYNAMIC_DRAW)
		sb.dirty = false
	}

	gl.BindBufferBase(gl.UNIFORM_BUFFER, 0, sb.id)
}
