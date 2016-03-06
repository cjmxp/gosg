package core

import "unsafe"

// Texture is an interface which wraps a GPU texture (OpenGL, etc). This will be removed soon and textures
// will be accessed by name handle via the resource system or abstracted in opaque material definitions.
type Texture interface {
	Handle() unsafe.Pointer
}
