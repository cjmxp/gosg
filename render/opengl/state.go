package opengl

import (
	"github.com/fcvarela/gosg/core"

	"github.com/fcvarela/gosg/protos"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	clearMaterial       *protos.Material
	currentMaterial     *protos.Material
	textureUnitBindings = map[uint32]uint32{}
)

func bindMaterialState(ub core.UniformBuffer, material *protos.Material, force bool) *Program {
	if material.DepthTest != currentMaterial.DepthTest || force {
		if material.DepthTest {
			gl.Enable(gl.DEPTH_TEST)
		} else {
			gl.Disable(gl.DEPTH_TEST)
		}
	}

	if material.DepthFunc != currentMaterial.DepthFunc || force {
		switch material.DepthFunc {
		case protos.Material_DEPTH_LESS:
			gl.DepthFunc(gl.LESS)
		case protos.Material_DEPTH_LESS_EQUAL:
			gl.DepthFunc(gl.LEQUAL)
		case protos.Material_DEPTH_EQUAL:
			gl.DepthFunc(gl.EQUAL)
		}
	}

	if material.ScissorTest != currentMaterial.ScissorTest || force {
		if material.ScissorTest {
			gl.Enable(gl.SCISSOR_TEST)
		} else {
			gl.Disable(gl.SCISSOR_TEST)
		}
	}

	if material.Blending != currentMaterial.Blending || force {
		if material.Blending {
			gl.Enable(gl.BLEND)
		} else {
			gl.Disable(gl.BLEND)
		}
	}

	if material.BlendSrcMode != currentMaterial.BlendSrcMode || material.BlendDstMode != currentMaterial.BlendDstMode || force {
		srcMode := uint32(0)
		dstMode := uint32(0)

		switch material.BlendSrcMode {
		case protos.Material_BLEND_ONE:
			srcMode = gl.ONE
		case protos.Material_BLEND_ONE_MINUS_SRC_ALPHA:
			srcMode = gl.ONE_MINUS_SRC_ALPHA
		case protos.Material_BLEND_SRC_ALPHA:
			srcMode = gl.SRC_ALPHA
		}

		switch material.BlendDstMode {
		case protos.Material_BLEND_ONE:
			dstMode = gl.ONE
		case protos.Material_BLEND_ONE_MINUS_SRC_ALPHA:
			dstMode = gl.ONE_MINUS_SRC_ALPHA
		case protos.Material_BLEND_SRC_ALPHA:
			dstMode = gl.SRC_ALPHA
		}
		gl.BlendFunc(srcMode, dstMode)
	}

	if material.BlendEquation != currentMaterial.BlendEquation || force {
		switch material.BlendEquation {
		case protos.Material_BLEND_FUNC_ADD:
			gl.BlendEquation(gl.FUNC_ADD)
		case protos.Material_BLEND_FUNC_MAX:
			gl.BlendEquation(gl.MAX)
		}
	}

	if material.DepthWrite != currentMaterial.DepthWrite || force {
		gl.DepthMask(material.DepthWrite)
	}

	if material.ColorWrite != currentMaterial.ColorWrite || force {
		gl.ColorMask(material.ColorWrite, material.ColorWrite, material.ColorWrite, material.ColorWrite)
	}

	if material.Culling != currentMaterial.Culling || force {
		if material.Culling {
			gl.Enable(gl.CULL_FACE)
		} else {
			gl.Disable(gl.CULL_FACE)
		}
	}

	if material.CullFace != currentMaterial.CullFace || force {
		switch material.CullFace {
		case protos.Material_CULL_BACK:
			gl.CullFace(gl.BACK)
		case protos.Material_CULL_FRONT:
			gl.CullFace(gl.FRONT)
		case protos.Material_CULL_BOTH:
			gl.CullFace(gl.FRONT_AND_BACK)
		}
	}

	glProgram := core.GetResourceManager().Program(material.ProgramName).(*Program)

	if material.ProgramName != currentMaterial.ProgramName || force {
		glProgram.bind()
	}

	// global uniform buffer
	if ub != nil {
		glProgram.setUniformBufferByName("cameraConstants", ub.(*UniformBuffer))
	}

	currentMaterial = material

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
