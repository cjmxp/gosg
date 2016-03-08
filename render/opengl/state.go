package opengl

import (
	"github.com/fcvarela/gosg/core"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	currentState = core.NewState()
)

func bindRenderState(ub core.UniformBuffer, s core.State, force bool) {
	// apply depth test function
	if s.Depth.Enabled != currentState.Depth.Enabled || force {
		if s.Depth.Enabled {
			gl.Enable(gl.DEPTH_TEST)
		} else {
			gl.Disable(gl.DEPTH_TEST)
		}
	}

	if s.Depth.Func != currentState.Depth.Func || force {
		switch s.Depth.Func {
		case core.DepthLess:
			gl.DepthFunc(gl.LESS)
		case core.DepthLessEqual:
			gl.DepthFunc(gl.LEQUAL)
		case core.DepthEqual:
			gl.DepthFunc(gl.EQUAL)
		}
	}

	if s.Scissor.Enabled != currentState.Scissor.Enabled || force {
		if s.Scissor.Enabled {
			gl.Enable(gl.SCISSOR_TEST)
		} else {
			gl.Disable(gl.SCISSOR_TEST)
		}
	}

	if s.Blend.Enabled != currentState.Blend.Enabled || force {
		if s.Blend.Enabled {
			gl.Enable(gl.BLEND)
		} else {
			gl.Disable(gl.BLEND)
		}
	}

	if s.Blend.SrcMode != currentState.Blend.SrcMode || s.Blend.DstMode != currentState.Blend.DstMode || force {
		srcMode := uint32(0)
		dstMode := uint32(0)
		switch s.Blend.SrcMode {
		case core.BlendOne:
			srcMode = gl.ONE
		case core.BlendOneMinusSrcAlpha:
			srcMode = gl.ONE_MINUS_SRC_ALPHA
		case core.BlendSrcAlpha:
			srcMode = gl.SRC_ALPHA
		}
		switch s.Blend.DstMode {
		case core.BlendOne:
			dstMode = gl.ONE
		case core.BlendOneMinusSrcAlpha:
			dstMode = gl.ONE_MINUS_SRC_ALPHA
		case core.BlendSrcAlpha:
			dstMode = gl.SRC_ALPHA
		}
		gl.BlendFunc(srcMode, dstMode)
	}

	if s.Blend.Equation != currentState.Blend.Equation || force {
		switch s.Blend.Equation {
		case core.BlendFuncAdd:
			gl.BlendEquation(gl.FUNC_ADD)
		case core.BlendFuncMax:
			gl.BlendEquation(gl.MAX)
		}
	}

	if s.Depth.Mask != currentState.Depth.Mask || force {
		gl.DepthMask(s.Depth.Mask)
	}

	if s.Color.Mask != currentState.Color.Mask || force {
		gl.ColorMask(s.Color.Mask, s.Color.Mask, s.Color.Mask, s.Color.Mask)
	}

	if s.Cull.Enabled != currentState.Cull.Enabled || force {
		if s.Cull.Enabled {
			gl.Enable(gl.CULL_FACE)
		} else {
			gl.Disable(gl.CULL_FACE)
		}
	}

	if s.Cull.Mode != currentState.Cull.Mode || force {
		switch s.Cull.Mode {
		case core.CullBack:
			gl.CullFace(gl.BACK)
		case core.CullFront:
			gl.CullFace(gl.FRONT)
		case core.CullBoth:
			gl.CullFace(gl.FRONT_AND_BACK)
		}
	}

	if s.Program() != nil {
		glProgram := s.Program().(*Program)

		if s.Program() != currentState.Program() || force {
			glProgram.bind()
		}

		// bind this node's uniforms
		for name, uniform := range s.Uniforms() {
			glProgram.setUniform(name, uniform.(*Uniform))
		}

		// global uniform buffer
		if ub != nil {
			glProgram.setUniformBufferByName("cameraConstants", ub.(*UniformBuffer))
		}

		// and node uniform buffers
		for name, uniformBuffer := range s.UniformBuffers() {
			glProgram.setUniformBufferByName(name, uniformBuffer.(*UniformBuffer))
		}

		// activate textures
		for sampler, texture := range s.Textures() {
			if currentState.Textures()[sampler] != texture {
				bindTexture(sampler, texture.(*Texture))
			}
		}
	}

	// let everyone else know this is active
	currentState = s
}
