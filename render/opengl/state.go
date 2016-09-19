package opengl

import (
	"github.com/fcvarela/gosg/core"

	"github.com/fcvarela/gosg/protos"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	clearState          *protos.State
	currentState        *protos.State
	textureUnitBindings = map[uint32]uint32{}
)

func bindMaterialState(ub core.UniformBuffer, material *protos.State, force bool) *Program {
	if material.DepthTest != currentState.DepthTest || force {
		if material.DepthTest {
			gl.Enable(gl.DEPTH_TEST)
		} else {
			gl.Disable(gl.DEPTH_TEST)
		}
	}

	if material.DepthFunc != currentState.DepthFunc || force {
		switch material.DepthFunc {
		case protos.State_DEPTH_LESS:
			gl.DepthFunc(gl.LESS)
		case protos.State_DEPTH_LESS_EQUAL:
			gl.DepthFunc(gl.LEQUAL)
		case protos.State_DEPTH_EQUAL:
			gl.DepthFunc(gl.EQUAL)
		}
	}

	if material.ScissorTest != currentState.ScissorTest || force {
		if material.ScissorTest {
			gl.Enable(gl.SCISSOR_TEST)
		} else {
			gl.Disable(gl.SCISSOR_TEST)
		}
	}

	if material.Blending != currentState.Blending || force {
		if material.Blending {
			gl.Enable(gl.BLEND)
		} else {
			gl.Disable(gl.BLEND)
		}
	}

	if material.BlendSrcMode != currentState.BlendSrcMode || material.BlendDstMode != currentState.BlendDstMode || force {
		srcMode := uint32(0)
		dstMode := uint32(0)

		switch material.BlendSrcMode {
		case protos.State_BLEND_ONE:
			srcMode = gl.ONE
		case protos.State_BLEND_ONE_MINUS_SRC_ALPHA:
			srcMode = gl.ONE_MINUS_SRC_ALPHA
		case protos.State_BLEND_SRC_ALPHA:
			srcMode = gl.SRC_ALPHA
		}

		switch material.BlendDstMode {
		case protos.State_BLEND_ONE:
			dstMode = gl.ONE
		case protos.State_BLEND_ONE_MINUS_SRC_ALPHA:
			dstMode = gl.ONE_MINUS_SRC_ALPHA
		case protos.State_BLEND_SRC_ALPHA:
			dstMode = gl.SRC_ALPHA
		}
		gl.BlendFunc(srcMode, dstMode)
	}

	if material.BlendEquation != currentState.BlendEquation || force {
		switch material.BlendEquation {
		case protos.State_BLEND_FUNC_ADD:
			gl.BlendEquation(gl.FUNC_ADD)
		case protos.State_BLEND_FUNC_MAX:
			gl.BlendEquation(gl.MAX)
		}
	}

	if material.DepthWrite != currentState.DepthWrite || force {
		gl.DepthMask(material.DepthWrite)
	}

	if material.ColorWrite != currentState.ColorWrite || force {
		gl.ColorMask(material.ColorWrite, material.ColorWrite, material.ColorWrite, material.ColorWrite)
	}

	if material.Culling != currentState.Culling || force {
		if material.Culling {
			gl.Enable(gl.CULL_FACE)
		} else {
			gl.Disable(gl.CULL_FACE)
		}
	}

	if material.CullFace != currentState.CullFace || force {
		switch material.CullFace {
		case protos.State_CULL_BACK:
			gl.CullFace(gl.BACK)
		case protos.State_CULL_FRONT:
			gl.CullFace(gl.FRONT)
		case protos.State_CULL_BOTH:
			gl.CullFace(gl.FRONT_AND_BACK)
		}
	}

	glProgram := core.GetResourceManager().Program(material.ProgramName).(*Program)

	if material.ProgramName != currentState.ProgramName || force {
		glProgram.bind()
	}

	// global uniform buffer
	if ub != nil {
		glProgram.setUniformBufferByName("cameraConstants", ub.(*UniformBuffer))
	}

	currentState = material

	return glProgram
}

func breaksBatch(a *core.MaterialData, b *core.MaterialData) bool {
	for name := range b.Textures() {
		ta, ok := a.Textures()[name]
		if !ok {
			return true
		}

		tb, ok := b.Textures()[name]
		if !ok {
			return true
		}

		if ta.(*Texture).ID != tb.(*Texture).ID {
			return true
		}
	}

	return false
}

func bindTextures(p *Program, md *core.MaterialData) {
	for name, texture := range md.Textures() {
		p.setTexture(name, texture.(*Texture))
	}
}

func bindUniformBuffers(p *Program, md *core.MaterialData) {
	for name, uniformBuffer := range md.UniformBuffers() {
		p.setUniformBufferByName(name, uniformBuffer.(*UniformBuffer))
	}
}

/*
func bindUniforms(p *Program, md *core.MaterialData) {
	for name, uniform := range md.Uniforms() {
		p.setUniform(name, uniform.(*Uniform))
	}
}
*/
