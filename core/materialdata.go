package core

// MaterialData contains material properties for a specific drawable
type MaterialData struct {
	uniforms       map[string]Uniform
	uniformBuffers map[string]UniformBuffer
	textures       map[string]Texture
}

// NewMaterialData returns a new MaterialData
func NewMaterialData() MaterialData {
	s := MaterialData{
		make(map[string]Uniform),
		make(map[string]UniformBuffer),
		make(map[string]Texture),
	}
	return s
}

// Uniforms returns the state's uniform map.
func (s *MaterialData) Uniforms() map[string]Uniform {
	return s.uniforms
}

// Uniform returns the uniform with the given name.
func (s *MaterialData) Uniform(name string) Uniform {
	_, ok := s.uniforms[name]
	if ok == false {
		s.uniforms[name] = renderSystem.NewUniform()
	}
	return s.uniforms[name]
}

// SetTexture sets the material's texture named `name` to the provided texture
func (s *MaterialData) SetTexture(name string, t Texture) {
	s.textures[name] = t
}

// Textures returns the material's textures
func (s *MaterialData) Textures() map[string]Texture {
	return s.textures
}

// UniformBuffer returns the uniform buffer with the given name
func (s *MaterialData) UniformBuffer(name string) UniformBuffer {
	_, ok := s.uniformBuffers[name]
	if !ok {
		s.uniformBuffers[name] = renderSystem.NewUniformBuffer()
	}
	return s.uniformBuffers[name]
}

// UniformBuffers returns the state's uniform buffers
func (s *MaterialData) UniformBuffers() map[string]UniformBuffer {
	return s.uniformBuffers
}
