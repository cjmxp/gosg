package opengl

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg" // registers jpeg handler
	_ "image/png"  // registers png handler
	"runtime"
	"unsafe"

	"github.com/fcvarela/gosg/core"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/glog"
)

// Texture is an OpenGL texture container
type Texture struct {
	ID uint32
}

// Handle implements the core.Texture interface
func (t *Texture) Handle() unsafe.Pointer {
	return unsafe.Pointer(t)
}

// Lt implements the core.Texture interface
func (t *Texture) Lt(other core.Texture) bool {
	if ot, ok := other.(*Texture); ok {
		return t.ID < ot.ID
	}
	return true
}

// Gt implements the core.Texture interface
func (t *Texture) Gt(other core.Texture) bool {
	if ot, ok := other.(*Texture); ok {
		return t.ID > ot.ID
	}
	return false
}

func textureCleanup(t *Texture) {
	glog.Info("Deleting texture: ", t.ID)
}

// NewTexture implements the core.RenderSystem interface
func (rs *RenderSystem) NewTexture(data []byte) core.Texture {
	if data == nil {
		glog.Fatal("Cannot read texture...")
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		glog.Warning("Cannot decode texture image: ", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("Unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	texture := uint32(0)
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	t := &Texture{texture}
	runtime.SetFinalizer(t, textureCleanup)
	return t
}

// NewRawTexture implements the core.RenderSystem interface
func (rs *RenderSystem) NewRawTexture(width, height int, payload []byte) core.Texture {
	texture := uint32(0)
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(width),
		int32(height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(payload))

	t := &Texture{texture}
	runtime.SetFinalizer(t, textureCleanup)
	return t
}
