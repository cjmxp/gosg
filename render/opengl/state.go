package opengl

import (
	"github.com/fcvarela/gosg/core"

	"github.com/fcvarela/gosg/protos"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	currentState        = "clear"
	textureUnitBindings = map[uint32]uint32{}
)

func bindMaterialState(ub core.UniformBuffer, materialName string, force bool) *Program {
	s := core.GetResourceManager().Material(materialName)
	c := core.GetResourceManager().Material(currentState)

	if s.DepthTest != c.DepthTest || force {
		if s.DepthTest {
			gl.Enable(gl.DEPTH_TEST)
		} else {
			gl.Disable(gl.DEPTH_TEST)
		}
	}

	if s.DepthFunc != c.DepthFunc || force {
		switch s.DepthFunc {
		case protos.Material_DEPTH_LESS:
			gl.DepthFunc(gl.LESS)
		case protos.Material_DEPTH_LESS_EQUAL:
			gl.DepthFunc(gl.LEQUAL)
		case protos.Material_DEPTH_EQUAL:
			gl.DepthFunc(gl.EQUAL)
		}
	}

	if s.ScissorTest != c.ScissorTest || force {
		if s.ScissorTest {
			gl.Enable(gl.SCISSOR_TEST)
		} else {
			gl.Disable(gl.SCISSOR_TEST)
		}
	}

	if s.Blending != c.Blending || force {
		if s.Blending {
			gl.Enable(gl.BLEND)
		} else {
			gl.Disable(gl.BLEND)
		}
	}

	if s.BlendSrcMode != c.BlendSrcMode || s.BlendDstMode != c.BlendDstMode || force {
		srcMode := uint32(0)
		dstMode := uint32(0)

		switch s.BlendSrcMode {
		case protos.Material_BLEND_ONE:
			srcMode = gl.ONE
		case protos.Material_BLEND_ONE_MINUS_SRC_ALPHA:
			srcMode = gl.ONE_MINUS_SRC_ALPHA
		case protos.Material_BLEND_SRC_ALPHA:
			srcMode = gl.SRC_ALPHA
		}

		switch s.BlendDstMode {
		case protos.Material_BLEND_ONE:
			dstMode = gl.ONE
		case protos.Material_BLEND_ONE_MINUS_SRC_ALPHA:
			dstMode = gl.ONE_MINUS_SRC_ALPHA
		case protos.Material_BLEND_SRC_ALPHA:
			dstMode = gl.SRC_ALPHA
		}
		gl.BlendFunc(srcMode, dstMode)
	}

	if s.BlendEquation != c.BlendEquation || force {
		switch s.BlendEquation {
		case protos.Material_BLEND_FUNC_ADD:
			gl.BlendEquation(gl.FUNC_ADD)
		case protos.Material_BLEND_FUNC_MAX:
			gl.BlendEquation(gl.MAX)
		}
	}

	if s.DepthWrite != c.DepthWrite || force {
		gl.DepthMask(s.DepthWrite)
	}

	if s.ColorWrite != c.ColorWrite || force {
		gl.ColorMask(s.ColorWrite, s.ColorWrite, s.ColorWrite, s.ColorWrite)
	}

	if s.Culling != c.Culling || force {
		if s.Culling {
			gl.Enable(gl.CULL_FACE)
		} else {
			gl.Disable(gl.CULL_FACE)
		}
	}

	if s.CullFace != c.CullFace || force {
		switch s.CullFace {
		case protos.Material_CULL_BACK:
			gl.CullFace(gl.BACK)
		case protos.Material_CULL_FRONT:
			gl.CullFace(gl.FRONT)
		case protos.Material_CULL_BOTH:
			gl.CullFace(gl.FRONT_AND_BACK)
		}
	}

	glProgram := core.GetResourceManager().Program(s.ProgramName).(*Program)

	if s.ProgramName != c.ProgramName || force {
		glProgram.bind()
	}

	// global uniform buffer
	if ub != nil {
		glProgram.setUniformBufferByName("cameraConstants", ub.(*UniformBuffer))
	}

	currentState = s.Name

	return glProgram
}

func breaksBatch(p *Program, md *core.MaterialData) bool {
	for name, texture := range md.Textures() {
		textureUnit := p.samplerBindings[name]
		if textureUnitBindings[textureUnit] != texture.(*Texture).ID {
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

func bindUniforms(p *Program, md *core.MaterialData) {
	for name, uniform := range md.Uniforms() {
		p.setUniform(name, uniform.(*Uniform))
	}
}
