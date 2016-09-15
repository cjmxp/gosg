package opengl

import (
	"github.com/fcvarela/gosg/core"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/golang/glog"
)

// RenderTarget implements the core.RenderTarget interface
type RenderTarget struct {
	width         uint32
	height        uint32
	fbo           uint32
	depthLayers   uint8
	depthTexture  uint32
	colorTextures []uint32
}

// NewRenderTarget implements the core.RenderSystem interface
func (r *RenderSystem) NewRenderTarget(width uint32, height uint32, depthLayers uint8, layers uint8) core.RenderTarget {
	rt := &RenderTarget{width, height, 0, depthLayers, 0, make([]uint32, layers)}

	// create the FB & bind it
	gl.GenFramebuffers(1, &rt.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, rt.fbo)

	// create texture handle
	if depthLayers == 1 {
		gl.GenTextures(1, &rt.depthTexture)
		gl.BindTexture(gl.TEXTURE_2D, rt.depthTexture)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32F, int32(rt.width), int32(rt.height), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
		gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, rt.depthTexture, 0)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	} else if depthLayers > 1 {
		gl.GenTextures(1, &rt.depthTexture)
		gl.BindTexture(gl.TEXTURE_2D_ARRAY, rt.depthTexture)

		gl.TexImage3D(gl.TEXTURE_2D_ARRAY, 0, gl.DEPTH_COMPONENT32F, int32(width), int32(height), int32(depthLayers), gl.FALSE, gl.DEPTH_COMPONENT, gl.FLOAT, nil)

		for i := 0; i < int(depthLayers); i++ {
			gl.FramebufferTextureLayer(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, rt.depthTexture, 0, int32(i))
		}

		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_COMPARE_MODE, gl.NONE)
	}

	if layers > 0 {
		// prealloc drawbuffer holder
		drawBuffers := make([]uint32, layers)

		for i := 0; i < int(layers); i++ {
			gl.GenTextures(1, &rt.colorTextures[i])
			gl.BindTexture(gl.TEXTURE_2D, rt.colorTextures[i])
			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB32F, int32(rt.width), int32(rt.height), 0, gl.RGB, gl.FLOAT, nil)
			gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, uint32(gl.COLOR_ATTACHMENT0+i), gl.TEXTURE_2D, rt.colorTextures[i], 0)
			drawBuffers[i] = uint32(gl.COLOR_ATTACHMENT0 + i)

			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		}

		// call drawbuffers
		gl.DrawBuffers(int32(layers), &drawBuffers[0])
	} else {
		gl.DrawBuffer(gl.NONE)
	}

	// check status
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	switch status {
	case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
		glog.Fatal("Incomplete attachment")
	case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
		glog.Fatal("Missing attachment")
	case gl.FRAMEBUFFER_UNSUPPORTED:
		glog.Fatal("Invalid combination of internal formats")
	case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
		glog.Fatal("Incomplete draw buffer")
	case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
		glog.Fatal("Incomplete read buffer")
	default:
		break
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	return rt
}

// SetActiveDepthLayer implements the core.RenderTarger interface
func (rt *RenderTarget) SetActiveDepthLayer(l uint8) {
	gl.BindFramebuffer(gl.FRAMEBUFFER, rt.fbo)
	gl.FramebufferTextureLayer(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, rt.depthTexture, 0, int32(l))
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// DepthLayerCount implements the core.RenderTarger interface
func (rt *RenderTarget) DepthLayerCount() uint8 {
	return rt.depthLayers
}

// ColorTextureCount implements the core.RenderTarger interface
func (rt *RenderTarget) ColorTextureCount() uint8 {
	return uint8(len(rt.colorTextures))
}

// DepthTexture implements the core.RenderTarger interface
func (rt *RenderTarget) DepthTexture() core.Texture {
	return &Texture{rt.depthTexture}
}

// ColorTexture implements the core.RenderTarger interface
func (rt *RenderTarget) ColorTexture(idx uint8) core.Texture {
	if int(idx) < len(rt.colorTextures) {
		return &Texture{rt.colorTextures[idx]}
	}

	return &Texture{0}
}

func (rt *RenderTarget) bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, rt.fbo)
}
