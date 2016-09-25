package opengl

import (
	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/golang/glog"
)

// Framebuffer implements the core.Framebuffer interface
type Framebuffer struct {
	width         uint32
	height        uint32
	fbo           uint32
	depthTexture  core.Texture
	colorTextures []core.Texture
}

// NewFramebuffer implements the core.RenderSystem interface
// fixme: this should be explicitly attached later (depth, color attachments)
func (r *RenderSystem) NewFramebuffer(width uint32, height uint32, depth bool, layers uint8) core.Framebuffer {
	rt := &Framebuffer{width, height, 0, nil, nil}

	// create the FB & bind it
	gl.GenFramebuffers(1, &rt.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, rt.fbo)

	// create texture handle
	if depth {
		depthTextureDescriptor := core.TextureDescriptor{
			Width:         rt.width,
			Height:        rt.height,
			Mipmaps:       false,
			Target:        core.TextureTarget2D,
			Format:        core.TextureFormatDEPTH,
			SizedFormat:   core.TextureSizedFormatDEPTH32F,
			ComponentType: core.TextureComponentTypeFLOAT,
			Filter:        core.TextureFilterNearest,
			WrapMode:      core.TextureWrapModeClampEdge,
		}
		rt.depthTexture = r.NewTexture(depthTextureDescriptor, nil)

		// fixme: do we need this bind?
		gl.BindTexture(gl.TEXTURE_2D, rt.depthTexture.(*Texture).id)
		gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, gl.DEPTH_ATTACHMENT, rt.depthTexture.(*Texture).id, 0)
	}

	if layers > 0 {
		// prealloc drawbuffer holder
		rt.colorTextures = make([]core.Texture, layers)
		drawBuffers := make([]uint32, layers)

		for i := range rt.colorTextures {
			textureDescriptor := core.TextureDescriptor{
				Width:         rt.width,
				Height:        rt.height,
				Mipmaps:       false,
				Target:        core.TextureTarget2D,
				Format:        core.TextureFormatRGB,
				SizedFormat:   core.TextureSizedFormatRGBA32F,
				ComponentType: core.TextureComponentTypeFLOAT,
				Filter:        core.TextureFilterLinear,
				WrapMode:      core.TextureWrapModeClampEdge,
			}
			rt.colorTextures[i] = r.NewTexture(textureDescriptor, nil)
			drawBuffers[i] = uint32(gl.COLOR_ATTACHMENT0 + i)

			gl.BindTexture(gl.TEXTURE_2D, rt.colorTextures[i].(*Texture).id)
			gl.FramebufferTexture(gl.DRAW_FRAMEBUFFER, uint32(gl.COLOR_ATTACHMENT0+i), rt.colorTextures[i].(*Texture).id, 0)
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

// ColorTextureCount implements the core.RenderTarger interface
func (rt *Framebuffer) ColorTextureCount() uint8 {
	return uint8(len(rt.colorTextures))
}

// DepthTexture implements the core.RenderTarger interface
func (rt *Framebuffer) DepthTexture() core.Texture {
	return rt.depthTexture
}

// ColorTexture implements the core.RenderTarger interface
func (rt *Framebuffer) ColorTexture(idx uint8) core.Texture {
	if int(idx) < len(rt.colorTextures) {
		return rt.colorTextures[idx]
	}
	return nil
}

func (rt *Framebuffer) bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, rt.fbo)
}
