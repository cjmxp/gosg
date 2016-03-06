package opengl

import (
	"reflect"

	"github.com/fcvarela/gosg/core"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/golang/glog"
)

// binds uniform in currently active state
func bindUniform(name string, p *Program, u *core.Uniform) {
	if u.Dirty() == false || u.Value() == nil {
		return
	}

	uloc, found := p.uniformLocations[name]
	if found == false {
		return
	}

	switch uval := u.Value().(type) {
	case mgl32.Mat4:
		gl.UniformMatrix4fv(uloc, 1, false, &uval[0])
		u.SetDirty(false)
	case mgl64.Mat4:
		newval := core.Mat4DoubleToFloat(uval)
		gl.UniformMatrix4fv(uloc, 1, false, &newval[0])
		u.SetDirty(false)
	case mgl64.Vec4:
		newval := core.Vec4DoubleToFloat(uval)
		gl.Uniform4fv(uloc, 1, &newval[0])
		u.SetDirty(false)
	case mgl64.Vec3:
		newval := core.Vec3DoubleToFloat(uval)
		gl.Uniform3fv(uloc, 1, &newval[0])
		u.SetDirty(false)
	case mgl64.Vec2:
		newval := core.Vec2DoubleToFloat(uval)
		gl.Uniform2fv(uloc, 1, &newval[0])
		u.SetDirty(false)
	case []float32:
		gl.Uniform1fv(uloc, int32(len(uval)), &uval[0])
		u.SetDirty(false)
	case []mgl32.Vec2:
		newval := make([]float32, len(uval)*2)
		for i := 0; i < len(uval); i++ {
			newval[i*2+0] = uval[i].X()
			newval[i*2+1] = uval[i].Y()
		}
		gl.Uniform2fv(uloc, int32(len(uval)), &newval[0])
		u.SetDirty(false)
	case int:
		gl.Uniform1i(uloc, int32(uval))
		u.SetDirty(false)
	case float32:
		gl.Uniform1f(uloc, uval)
		u.SetDirty(false)
	case float64:
		gl.Uniform1f(uloc, float32(uval))
		u.SetDirty(false)
	default:
		glog.Fatalln("UNSUPPORTED -- Uniform: %s Type: %s\n", name, reflect.TypeOf(u.Value()))
	}
}
