package core

import (
	"reflect"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/golang/glog"
)

// Uniform is a GPU program uniform of a generic type.
type Uniform struct {
	value interface{}
	dirty bool
}

// Set sets the uniform's value.
func (u *Uniform) Set(value interface{}) {
	u.value = value
	u.dirty = true
}

// Value returns the uniform's value.
func (u *Uniform) Value() interface{} {
	return u.value
}

// Copy returns a copy of the uniform.
func (u *Uniform) Copy() *Uniform {
	switch uval := u.Value().(type) {
	case mgl32.Mat4:
		return &Uniform{uval, false}
	case mgl64.Mat4:
		return &Uniform{uval, false}
	case mgl64.Vec4:
		return &Uniform{uval, false}
	case mgl64.Vec3:
		return &Uniform{uval, false}
	case mgl64.Vec2:
		return &Uniform{uval, false}
	case []float32:
		return &Uniform{uval, false}
	case []mgl32.Vec2:
		return &Uniform{uval, false}
	case int:
		return &Uniform{uval, false}
	case float32:
		return &Uniform{uval, false}
	case float64:
		return &Uniform{uval, false}
	default:
		glog.Fatalf("Unsupported uniform type: %s\n", reflect.TypeOf(u.Value()))
	}

	return nil
}

// SetDirty sets the uniform dirty flag.
func (u *Uniform) SetDirty(d bool) {
	u.dirty = d
}

// Dirty returns whether the uniform is dirty.
func (u *Uniform) Dirty() bool {
	return u.dirty
}
