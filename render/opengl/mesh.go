package opengl

import (
	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl64"
)

type bufferType uint8

const (
	positionBuffer bufferType = iota
	normalBuffer
	tangentBuffer
	bitangentBuffer
	texCoordBuffer
	indexBuffer
	boneIDBuffer
	boneWeightBuffer
	modelMatrixBuffer
)

// Mesh implements the core.Mesh interface
type Mesh struct {
	vao           uint32
	dirty         bool
	indexcount    int32
	name          string
	buffers       map[bufferType]uint32
	bounds        *core.AABB
	primitiveType uint32
	compileList   []func()
}

// InstancedMesh implements the core.InstancedMesh interface
type InstancedMesh struct {
	*Mesh
	instanceCount int32
}

// IMGUIMesh implements the core.IMGUIMesh interface
type IMGUIMesh struct {
	*Mesh
}

// NewMesh implements the core.RenderSystem interface
func (r *RenderSystem) NewMesh() core.Mesh {
	m := Mesh{}

	m.compileList = make([]func(), 0)
	m.bounds = core.NewAABB()
	m.buffers = make(map[bufferType]uint32)

	for i := positionBuffer; i < modelMatrixBuffer; i++ {
		m.buffers[i] = 0
	}

	m.dirty = false

	// add to compile list
	m.compileList = append(m.compileList, func() {
		if m.vao == 0 {
			gl.GenVertexArrays(1, &m.vao)
		}
	})

	m.dirty = true
	return &m
}

// NewIMGUIMesh implements the core.RenderSystem interface
func (r *RenderSystem) NewIMGUIMesh() core.IMGUIMesh {
	return &IMGUIMesh{r.NewMesh().(*Mesh)}
}

// NewInstancedMesh implements the core.RenderSystem interface
func (r *RenderSystem) NewInstancedMesh(m core.Mesh) core.InstancedMesh {
	if m != nil {
		return &InstancedMesh{m.(*Mesh), 0}
	}
	return &InstancedMesh{r.NewMesh().(*Mesh), 0}
}

func (m *Mesh) compile() {
	for _, fn := range m.compileList {
		fn()
	}

	m.compileList = m.compileList[:0]
	m.dirty = false
}

// SetPrimitiveType implements the core.Mesh interface
func (m *Mesh) SetPrimitiveType(t core.PrimitiveType) {
	switch t {
	case core.PrimitiveTypeTriangles:
		m.primitiveType = gl.TRIANGLES
	case core.PrimitiveTypeLines:
		m.primitiveType = gl.LINES
	case core.PrimitiveTypePoints:
		m.primitiveType = gl.POINTS
		m.compileList = append(m.compileList, func() {
			gl.Enable(gl.PROGRAM_POINT_SIZE)
			gl.PointParameteri(gl.POINT_SPRITE_COORD_ORIGIN, gl.LOWER_LEFT)
		})
		m.dirty = true
	}
}

// Bounds implements the core.Mesh interface
func (m *Mesh) Bounds() *core.AABB {
	return m.bounds
}

// SetName implements the core.Mesh interface
func (m *Mesh) SetName(name string) {
	m.name = name
}

// Name implements the core.Mesh interface
func (m *Mesh) Name() string {
	return m.name
}

// SetPositions implements the core.Mesh interface
func (m *Mesh) SetPositions(positions []float32) {
	m.compileList = append(m.compileList, func() {
		gl.BindVertexArray(m.vao)

		buf := uint32(0)

		if m.buffers[positionBuffer] == 0 {
			gl.GenBuffers(1, &buf)
			m.buffers[positionBuffer] = buf
		} else {
			buf = m.buffers[positionBuffer]
		}

		//4 : sizeof float32
		gl.BindBuffer(gl.ARRAY_BUFFER, buf)
		gl.BufferData(gl.ARRAY_BUFFER, len(positions)*4, gl.Ptr(positions), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	})

	// grow our bounds
	for i := 0; i < len(positions); i += 3 {
		p := mgl64.Vec3{
			float64(positions[i+0]),
			float64(positions[i+1]),
			float64(positions[i+2])}
		m.bounds.ExtendWithPoint(p)
	}

	m.dirty = true
}

// SetNormals implements the core.Mesh interface
func (m *Mesh) SetNormals(normals []float32) {
	m.compileList = append(m.compileList, func() {
		gl.BindVertexArray(m.vao)

		buf := uint32(0)

		if m.buffers[normalBuffer] == 0 {
			gl.GenBuffers(1, &buf)
			m.buffers[normalBuffer] = buf
		} else {
			buf = m.buffers[normalBuffer]
		}

		//4 : sizeof float32
		gl.BindBuffer(gl.ARRAY_BUFFER, buf)
		gl.BufferData(gl.ARRAY_BUFFER, len(normals)*4, gl.Ptr(normals), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	})

	m.dirty = true
}

// SetTangents implements the core.Mesh interface
func (m *Mesh) SetTangents(tangents []float32) {
	m.compileList = append(m.compileList, func() {
		gl.BindVertexArray(m.vao)

		buf := uint32(0)

		if m.buffers[tangentBuffer] == 0 {
			gl.GenBuffers(1, &buf)
			m.buffers[tangentBuffer] = buf
		} else {
			buf = m.buffers[tangentBuffer]
		}

		//4 : sizeof float32
		gl.BindBuffer(gl.ARRAY_BUFFER, buf)
		gl.BufferData(gl.ARRAY_BUFFER, len(tangents)*4, gl.Ptr(tangents), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	})

	m.dirty = true
}

// SetBitangents implements the core.Mesh interface
func (m *Mesh) SetBitangents(bitangents []float32) {
	m.compileList = append(m.compileList, func() {
		gl.BindVertexArray(m.vao)

		buf := uint32(0)

		if m.buffers[bitangentBuffer] == 0 {
			gl.GenBuffers(1, &buf)
			m.buffers[bitangentBuffer] = buf
		} else {
			buf = m.buffers[bitangentBuffer]
		}

		//4 : sizeof float32
		gl.BindBuffer(gl.ARRAY_BUFFER, buf)
		gl.BufferData(gl.ARRAY_BUFFER, len(bitangents)*4, gl.Ptr(bitangents), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(3)
		gl.VertexAttribPointer(3, 3, gl.FLOAT, false, 0, nil)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	})

	m.dirty = true
}

// SetTextureCoordinates implements the core.Mesh interface
func (m *Mesh) SetTextureCoordinates(size int32, texcoords []float32) {
	m.compileList = append(m.compileList, func() {
		gl.BindVertexArray(m.vao)

		buf := uint32(0)

		if m.buffers[texCoordBuffer] == 0 {
			gl.GenBuffers(1, &buf)
			m.buffers[texCoordBuffer] = buf
		} else {
			buf = m.buffers[texCoordBuffer]
		}

		//4 : sizeof float32
		gl.BindBuffer(gl.ARRAY_BUFFER, buf)
		gl.BufferData(gl.ARRAY_BUFFER, len(texcoords)*4, gl.Ptr(texcoords), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(4)
		gl.VertexAttribPointer(4, size, gl.FLOAT, false, 0, nil)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	})

	m.dirty = true
}

// SetIndices implements the core.Mesh interface
func (m *Mesh) SetIndices(indices []uint16) {
	m.compileList = append(m.compileList, func() {
		gl.BindVertexArray(m.vao)

		buf := uint32(0)

		if m.buffers[indexBuffer] == 0 {
			gl.GenBuffers(1, &buf)
			m.buffers[indexBuffer] = buf
		} else {
			buf = m.buffers[indexBuffer]
		}

		//2 : sizeof uint16
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*2, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)

		m.indexcount = int32(len(indices))
	})

	m.dirty = true
}

// Draw implements the core.Mesh interface
func (m *Mesh) Draw(ub core.UniformBuffer, st *core.State) {
	if m.dirty {
		m.compile()
	}

	bindRenderState(ub, *st, false)

	gl.BindVertexArray(m.vao)
	gl.DrawElements(m.primitiveType, m.indexcount, gl.UNSIGNED_SHORT, nil)
	gl.BindVertexArray(0)
}

// Draw implements the core.IMGUIMesh interface
func (m *IMGUIMesh) Draw(ub core.UniformBuffer, st *core.State) {
	bindRenderState(ub, *st, false)

	imguiSystem := core.GetIMGUISystem()
	drawData := imguiSystem.GetDrawData()

	m.SetPositions([]float32{0.0, 0.0, 0.0})
	m.SetIndices([]uint16{0})

	m.compile()

	gl.BindVertexArray(m.vao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)

	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.ActiveTexture(gl.TEXTURE0)

	for i := 0; i < drawData.CommandListCount(); i++ {
		cmdlist := drawData.GetCommandList(i)

		gl.BindBuffer(gl.ARRAY_BUFFER, m.buffers[positionBuffer])
		gl.BufferData(gl.ARRAY_BUFFER, cmdlist.VertexBufferSize*5*4, cmdlist.VertexPointer, gl.STREAM_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.buffers[indexBuffer])
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, cmdlist.IndexBufferSize*2, cmdlist.IndexPointer, gl.STREAM_DRAW)

		// position = 0, tcoords = 1, normals/color = 2
		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(8))
		gl.VertexAttribPointer(2, 4, gl.UNSIGNED_BYTE, true, 5*4, gl.PtrOffset(16))

		m.compile()

		var elementIndex int
		for _, cmd := range cmdlist.Commands {
			if tex := (*Texture)(cmd.TextureID); tex != nil {
				gl.BindTexture(gl.TEXTURE_2D, tex.ID)
			}

			gl.Scissor(
				int32(cmd.ClipRect[0]),
				int32(imguiSystem.DisplaySize().Y()-cmd.ClipRect[3]),
				int32(cmd.ClipRect[2]-cmd.ClipRect[0]),
				int32(cmd.ClipRect[3]-cmd.ClipRect[1]))
			gl.DrawElements(m.primitiveType, int32(cmd.ElementCount), gl.UNSIGNED_SHORT, gl.PtrOffset(elementIndex))
			elementIndex += cmd.ElementCount * 2
		}
	}

	gl.BindTexture(gl.TEXTURE_2D, (uint32)(lastTexture))

	gl.BindVertexArray(0)
}

// SetInstanceCount implements the core.InstancedMesh interface
func (i *InstancedMesh) SetInstanceCount(count int) {
	i.instanceCount = int32(count)
}

// SetModelMatrices implements the core.InstancedMesh interface
func (i *InstancedMesh) SetModelMatrices(matrices []float32) {
	i.compileList = append(i.compileList, func() {
		gl.BindVertexArray(i.vao)

		buf := uint32(0)

		if i.buffers[modelMatrixBuffer] == 0 {
			gl.GenBuffers(1, &buf)
			i.buffers[modelMatrixBuffer] = buf
		} else {
			buf = i.buffers[modelMatrixBuffer]
			gl.DeleteBuffers(1, &buf)
			gl.GenBuffers(1, &buf)
			i.buffers[modelMatrixBuffer] = buf
		}

		//4 : sizeof float32
		gl.BindBuffer(gl.ARRAY_BUFFER, buf)
		gl.BufferData(gl.ARRAY_BUFFER, len(matrices)*4, gl.Ptr(matrices), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(5)
		gl.VertexAttribPointer(5, 4, gl.FLOAT, false, 16*4, gl.PtrOffset(0*16))
		gl.VertexAttribDivisor(5, 1)

		gl.EnableVertexAttribArray(6)
		gl.VertexAttribPointer(6, 4, gl.FLOAT, false, 16*4, gl.PtrOffset(1*16))
		gl.VertexAttribDivisor(6, 1)

		gl.EnableVertexAttribArray(7)
		gl.VertexAttribPointer(7, 4, gl.FLOAT, false, 16*4, gl.PtrOffset(2*16))
		gl.VertexAttribDivisor(7, 1)

		gl.EnableVertexAttribArray(8)
		gl.VertexAttribPointer(8, 4, gl.FLOAT, false, 16*4, gl.PtrOffset(3*16))
		gl.VertexAttribDivisor(8, 1)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	})

	i.dirty = true
}

// Draw implements the core.InstancedMesh interface
func (i *InstancedMesh) Draw(ub core.UniformBuffer, st *core.State) {
	// we compile on the fly with opengl, sucks
	if i.dirty {
		i.compile()
	}

	bindRenderState(ub, *st, false)

	gl.BindVertexArray(i.vao)
	gl.DrawElementsInstanced(i.primitiveType, i.indexcount, gl.UNSIGNED_SHORT, nil, i.instanceCount)
	gl.BindVertexArray(0)
}
