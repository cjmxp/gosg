package opengl

import (
	"unsafe"

	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/gl/v4.1-core/gl"
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
	modelMatrixBuffer
)

type buffers struct {
	vao           uint32
	buffers       []uint32
	bufferOffsets []int
}

var (
	currentVAO = uint32(0)
)

func (b *buffers) addData(target uint32, buffer bufferType, datalen int, buf unsafe.Pointer) {
	if b.bufferOffsets[buffer] == 0 {
		gl.BindBuffer(target, b.buffers[buffer])
		gl.BufferData(target, datalen, buf, gl.STATIC_DRAW)
	} else {
		// cpu copy: get existing data, alloc new space for everything, add old data, new data
		cpuBuf := make([]byte, b.bufferOffsets[buffer])
		gl.BindBuffer(target, b.buffers[buffer])
		gl.GetBufferSubData(target, 0, b.bufferOffsets[buffer], unsafe.Pointer(&cpuBuf[0]))
		gl.BufferData(target, datalen+b.bufferOffsets[buffer], nil, gl.STATIC_DRAW)
		gl.BufferSubData(target, 0, b.bufferOffsets[buffer], unsafe.Pointer(&cpuBuf[0]))
		gl.BufferSubData(target, b.bufferOffsets[buffer], datalen, buf)
	}
	b.bufferOffsets[buffer] += datalen
}

func newBuffers() *buffers {
	bf := &buffers{}

	// create VAO
	gl.GenVertexArrays(1, &bf.vao)

	// create buffers
	bf.buffers = make([]uint32, modelMatrixBuffer+1)
	bf.bufferOffsets = make([]int, modelMatrixBuffer+1)

	// initialize gl buffer handles
	gl.GenBuffers(int32(len(bf.buffers)), &bf.buffers[0])

	// init attributes
	bindVAO(bf.vao)

	// position
	gl.BindBuffer(gl.ARRAY_BUFFER, bf.buffers[positionBuffer])
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// normal
	gl.BindBuffer(gl.ARRAY_BUFFER, bf.buffers[normalBuffer])
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	// tangent
	gl.BindBuffer(gl.ARRAY_BUFFER, bf.buffers[tangentBuffer])
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)

	// bitangent
	gl.BindBuffer(gl.ARRAY_BUFFER, bf.buffers[bitangentBuffer])
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, 0, nil)

	// texcoord
	gl.BindBuffer(gl.ARRAY_BUFFER, bf.buffers[texCoordBuffer])
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointer(4, 3, gl.FLOAT, false, 0, nil)

	// indices
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, bf.buffers[indexBuffer])

	// model matrices, prealloc for 200 instances of 4x4 matrices with 4 bytes per float (this is our hard-max)
	gl.BindBuffer(gl.ARRAY_BUFFER, bf.buffers[modelMatrixBuffer])
	gl.BufferData(gl.ARRAY_BUFFER, 200*16*4, nil, gl.STREAM_DRAW)

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

	bindVAO(0)
	return bf
}

var (
	sharedBuffers *buffers
	imguiBuffers  *buffers
)

// Mesh implements the core.Mesh interface
type Mesh struct {
	buffers           *buffers
	indexcount        int32
	indexOffset       int32
	indexBufferOffset int32
	instanceCount     int32
	name              string
	bounds            *core.AABB
	primitiveType     uint32
}

// IMGUIMesh implements the core.IMGUIMesh interface
type IMGUIMesh struct {
	*Mesh
}

// NewMesh implements the core.RenderSystem interface
func (r *RenderSystem) NewMesh() core.Mesh {
	m := Mesh{}

	m.bounds = core.NewAABB()
	m.buffers = sharedBuffers
	return &m
}

// NewIMGUIMesh implements the core.RenderSystem interface
func (r *RenderSystem) NewIMGUIMesh() core.IMGUIMesh {
	imguiMesh := &IMGUIMesh{r.NewMesh().(*Mesh)}
	imguiMesh.buffers = imguiBuffers
	return imguiMesh
}

func bindVAO(vao uint32) {
	if currentVAO != vao {
		gl.BindVertexArray(vao)
		currentVAO = vao
	}
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
		gl.Enable(gl.PROGRAM_POINT_SIZE)
		gl.PointParameteri(gl.POINT_SPRITE_COORD_ORIGIN, gl.LOWER_LEFT)
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
	// the index offset is equal to the number of primitives already in the position buffer
	m.indexOffset = int32(m.buffers.bufferOffsets[positionBuffer]) / (4 * 3)
	m.buffers.addData(gl.ARRAY_BUFFER, positionBuffer, len(positions)*4, gl.Ptr(positions))

	// grow our bounds
	for i := 0; i < len(positions); i += 3 {
		p := mgl64.Vec3{
			float64(positions[i+0]),
			float64(positions[i+1]),
			float64(positions[i+2])}
		m.bounds.ExtendWithPoint(p)
	}
}

// SetNormals implements the core.Mesh interface
func (m *Mesh) SetNormals(normals []float32) {
	m.buffers.addData(gl.ARRAY_BUFFER, normalBuffer, len(normals)*4, gl.Ptr(normals))
}

// SetTangents implements the core.Mesh interface
func (m *Mesh) SetTangents(tangents []float32) {
	m.buffers.addData(gl.ARRAY_BUFFER, tangentBuffer, len(tangents)*4, gl.Ptr(tangents))
}

// SetBitangents implements the core.Mesh interface
func (m *Mesh) SetBitangents(bitangents []float32) {
	m.buffers.addData(gl.ARRAY_BUFFER, bitangentBuffer, len(bitangents)*4, gl.Ptr(bitangents))
}

// SetTextureCoordinates implements the core.Mesh interface
func (m *Mesh) SetTextureCoordinates(texcoords []float32) {
	m.buffers.addData(gl.ARRAY_BUFFER, texCoordBuffer, len(texcoords)*4, gl.Ptr(texcoords))
}

// SetIndices implements the core.Mesh interface
func (m *Mesh) SetIndices(indices []uint16) {
	m.indexBufferOffset = int32(m.buffers.bufferOffsets[indexBuffer])
	m.indexcount = int32(len(indices))

	m.buffers.addData(gl.ELEMENT_ARRAY_BUFFER, indexBuffer, len(indices)*2, gl.Ptr(indices))
}

// Draw implements the core.Mesh interface
func (m *Mesh) Draw() {
	bindVAO(m.buffers.vao)
	gl.DrawElementsInstancedBaseVertex(
		m.primitiveType, m.indexcount, gl.UNSIGNED_SHORT,
		gl.PtrOffset(int(m.indexBufferOffset)), m.instanceCount, m.indexOffset)
}

// Draw implements the core.IMGUIMesh interface
func (m *IMGUIMesh) Draw() {
	bindVAO(m.buffers.vao)

	imguiSystem := core.GetIMGUISystem()
	drawData := imguiSystem.GetDrawData()

	m.SetPositions([]float32{0.0, 0.0, 0.0})
	m.SetIndices([]uint16{0})

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)

	var lastTexture int32
	var lastMipmapMode int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GetIntegerv(gl.TEXTURE_MIN_FILTER, &lastMipmapMode)
	gl.ActiveTexture(gl.TEXTURE0)

	for i := 0; i < drawData.CommandListCount(); i++ {
		cmdlist := drawData.GetCommandList(i)

		gl.BindBuffer(gl.ARRAY_BUFFER, m.buffers.buffers[positionBuffer])
		gl.BufferData(gl.ARRAY_BUFFER, cmdlist.VertexBufferSize*5*4, cmdlist.VertexPointer, gl.STREAM_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.buffers.buffers[indexBuffer])
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, cmdlist.IndexBufferSize*2, cmdlist.IndexPointer, gl.STREAM_DRAW)

		// position = 0, tcoords = 1, normals/color = 2
		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(8))
		gl.VertexAttribPointer(2, 4, gl.UNSIGNED_BYTE, true, 5*4, gl.PtrOffset(16))

		var elementIndex int
		for _, cmd := range cmdlist.Commands {
			if tex := (*Texture)(cmd.TextureID); tex != nil {
				gl.BindTexture(gl.TEXTURE_2D, tex.id)
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
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
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, lastMipmapMode)
	bindVAO(0)
}

// Lt implements the core.Mesh interface
func (m *Mesh) Lt(other core.Mesh) bool {
	return m.buffers.vao < other.(*Mesh).buffers.vao
}

// Gt implements the core.Mesh interface
func (m *Mesh) Gt(other core.Mesh) bool {
	return m.buffers.vao > other.(*Mesh).buffers.vao
}

// SetInstanceCount implements the core.InstancedMesh interface
func (m *Mesh) SetInstanceCount(count int) {
	m.instanceCount = int32(count)
}

// SetModelMatrices implements the core.InstancedMesh interface
func (m *Mesh) SetModelMatrices(matrices []float32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, m.buffers.buffers[modelMatrixBuffer])

	// orphaning technique #1, map, invalidate, write, unmap
	//buf := gl.MapBufferRange(gl.ARRAY_BUFFER, 0, len(matrices)*4,
	//	gl.MAP_WRITE_BIT|gl.MAP_INVALIDATE_BUFFER_BIT|gl.MAP_UNSYNCHRONIZED_BIT)
	//copy((*[1 << 30]float32)(buf)[:], matrices)
	//gl.UnmapBuffer(gl.ARRAY_BUFFER)

	// orphaning technique #2: realloc. old buf is still used by in-flight draw calls, should be fastest
	gl.BufferData(gl.ARRAY_BUFFER, len(matrices)*4, gl.Ptr(matrices), gl.STREAM_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
