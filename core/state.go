package core

// CullFace represents the faces which are to be culled.
type CullFace uint8

// Supported CullFace values
const (
	CullBack CullFace = 1 << iota
	CullFront
	CullBoth
)

// CullState holds information about which faces should be culled.
type CullState struct {
	Enabled bool
	Mode    CullFace
}

// BlendMode represents the blending mode.
type BlendMode uint8

// Supported BlendMode values
const (
	BlendSrcAlpha BlendMode = 1 << iota
	BlendOneMinusSrcAlpha
	BlendOne
)

// BlendEquation is an equation used to determine how blending should occur.
type BlendEquation uint8

// Supported BlendEquation values
const (
	BlendFuncAdd BlendEquation = 1 << iota
	BlendFuncMax
)

// BlendState holds information about equation, modes and status of a node's blending state.
type BlendState struct {
	Enabled  bool
	SrcMode  BlendMode
	DstMode  BlendMode
	Equation BlendEquation
}

// DepthFunc is an equation used to determine how depth should be tested.
type DepthFunc uint8

// Supported DepthFunc values
const (
	DepthLessEqual DepthFunc = 1 << iota
	DepthLess
	DepthEqual
)

// DepthState holds information about a node's depth function and mask.
type DepthState struct {
	Enabled bool
	Mask    bool
	Func    DepthFunc
}

// ColorState is a node's color state.
type ColorState struct {
	Mask bool
}

// ScissorState is a node's scissor state.
type ScissorState struct {
	Enabled bool
}

// State wraps sub states for a given node, along with material, bounding box drawing flags and program uniforms.
// This will soon become a material.
type State struct {
	uniforms map[string]*Uniform
	textures map[uint32]Texture
	program  Program
	Cull     CullState
	Blend    BlendState
	Depth    DepthState
	Color    ColorState
	Scissor  ScissorState
	AABB     bool
}

// NewAABBState returns a new state suitable to draw AABBs
func NewAABBState() State {
	st := NewState()
	st.Depth.Enabled = true
	st.Depth.Mask = false
	st.Color.Mask = true
	st.Blend.Enabled = false
	st.SetProgram(resourceManager.Program("zpass"))
	return st
}

// NewZPassState returns a new state suitable to draw depth buffer passes
func NewZPassState() State {
	st := NewState()
	st.Depth.Enabled = true
	st.Depth.Mask = true
	st.Depth.Func = DepthLessEqual
	st.Color.Mask = false
	st.Blend.Enabled = false
	st.Scissor.Enabled = false
	st.SetProgram(resourceManager.Program("zpass"))

	return st
}

// NewInstancedZPassState returns a new state suitable to draw depth buffer passes for instanced meshes
func NewInstancedZPassState() State {
	st := NewState()
	st.Depth.Enabled = true
	st.Depth.Mask = true
	st.Depth.Func = DepthLessEqual
	st.Color.Mask = false
	st.Blend.Enabled = false
	st.Scissor.Enabled = false
	st.SetProgram(resourceManager.Program("zpass-instanced"))

	return st
}

// NewState returns a new default state
func NewState() State {
	s := State{
		make(map[string]*Uniform),
		make(map[uint32]Texture),
		nil,
		CullState{true, CullBack},
		BlendState{false, BlendSrcAlpha, BlendOneMinusSrcAlpha, BlendFuncAdd},
		DepthState{true, true, DepthLessEqual},
		ColorState{true},
		ScissorState{false},
		false,
	}
	return s
}

// Copy deep copies the state
func (s *State) Copy() *State {
	ss := *s
	ss.uniforms = make(map[string]*Uniform)
	for k, v := range s.uniforms {
		ss.Uniform(k).Set(v.value)
	}
	for k, v := range s.textures {
		ss.SetTexture(k, v)
	}

	return &ss
}

// SetProgram sets the state's program.
func (s *State) SetProgram(p Program) {
	s.program = p
}

// Program returns the state's program.
func (s *State) Program() Program {
	return s.program
}

// Uniforms returns the state's uniform map.
func (s *State) Uniforms() map[string]*Uniform {
	return s.uniforms
}

// Uniform returns the uniform with the given name.
func (s *State) Uniform(name string) *Uniform {
	_, ok := s.uniforms[name]
	if ok == false {
		s.uniforms[name] = &Uniform{nil, true}
	}
	return s.uniforms[name]
}

// SetTexture sets the state's texture unit `unit` to the texture `t`
func (s *State) SetTexture(unit uint32, t Texture) {
	s.textures[unit] = t
}

// Textures returns the state's texture map
func (s *State) Textures() map[uint32]Texture {
	return s.textures
}
