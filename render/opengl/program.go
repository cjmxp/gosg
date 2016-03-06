package opengl

import (
	"runtime"
	"strings"

	"github.com/fcvarela/gosg/core"
	"github.com/golang/glog"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Program is an OpenGL program
type Program struct {
	name                string
	id                  uint32
	shaders             map[string]*shader
	compileLog          string
	uniformLog          string
	uniformLocations    map[string]int32
	uniformBlockIndexes map[string]uint32
}

func programCleanup(p *Program) {
	glog.Infof("Finalizer called for program: %v\n", p)
}

var (
	programTypeMap = map[string]uint32{
		"compute":            gl.COMPUTE_SHADER,
		"tesselationControl": gl.TESS_CONTROL_SHADER,
		"tesselationEval":    gl.TESS_EVALUATION_SHADER,
		"vertex":             gl.VERTEX_SHADER,
		"geometry":           gl.GEOMETRY_SHADER,
		"fragment":           gl.FRAGMENT_SHADER,
	}
)

// ProgramExtension implements the core.RenderSystem interface.
func (r *RenderSystem) ProgramExtension() string {
	return "gl.json"
}

// NewProgram implements the core.RenderSystem interface.
func (r *RenderSystem) NewProgram(name string, shaders map[string][]byte) core.Program {
	// create program
	prog := Program{
		name,
		0,
		make(map[string]*shader),
		"",
		"",
		make(map[string]int32),
		make(map[string]uint32),
	}

	// set shaders
	for k, v := range shaders {
		prog.shaders[k] = newShader(k, programTypeMap[k], v)
	}

	// set finalizer (hooks into gl to cleanup)
	runtime.SetFinalizer(&prog, programCleanup)

	// compile it
	prog.id = gl.CreateProgram()

	for _, s := range prog.shaders {
		if s != nil {
			s.compile()
			gl.AttachShader(prog.id, s.id)
		}
	}

	gl.LinkProgram(prog.id)
	status := int32(0)
	gl.GetProgramiv(prog.id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		logLength := int32(0)
		gl.GetProgramiv(prog.id, gl.INFO_LOG_LENGTH, &logLength)

		progLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(prog.id, logLength, nil, gl.Str(progLog))
		prog.compileLog = progLog

		glog.Fatalf("failed to link program: %v", progLog)
	}

	for _, s := range prog.shaders {
		if s != nil {
			gl.DeleteShader(s.id)
		}
	}

	// hard-bind scene block and node block(this is shared across all programs)
	sceneBlockIndex := gl.GetUniformBlockIndex(prog.id, gl.Str("sceneBlock"+"\x00"))
	gl.UniformBlockBinding(prog.id, sceneBlockIndex, 0)

	nodeBlockIndex := gl.GetUniformBlockIndex(prog.id, gl.Str("nodeBlock"+"\x00"))
	gl.UniformBlockBinding(prog.id, nodeBlockIndex, 1)

	// debug available uniforms
	var uniformcount int32
	gl.GetProgramiv(prog.id, gl.ACTIVE_UNIFORMS, &uniformcount)
	for i := uint32(0); i < uint32(uniformcount); i++ {
		// prepare content
		uniformName := strings.Repeat("\x00", int(128)+1)
		var uniformLen int32
		var uniformSize int32
		var uniformType uint32

		// extract locations and name
		gl.GetActiveUniform(prog.id, i, 128, &uniformLen, &uniformSize, &uniformType, gl.Str(uniformName))

		// extract location
		location := gl.GetUniformLocation(prog.id, gl.Str(uniformName))

		// save into location cache
		goUniformname := gl.GoStr(gl.Str(uniformName))
		prog.uniformLocations[goUniformname] = location
		//glog.Infof("Uniform: %s Size: %d Type: %d Slot: %d\n", goUniformname, uniformSize, uniformType, location)
	}
	//glog.Info(prog.uniformLocations)

	return &prog
}

// Name implements the core.Program interface
func (p *Program) Name() string {
	return p.name
}

func bindProgram(p *Program) {
	gl.UseProgram(p.id)
}

func unbindProgram() {
	gl.UseProgram(0)
}
