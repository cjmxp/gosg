package core

import (
	"fmt"
)

// PrimitiveType is a raster primitive type.
type PrimitiveType uint8

// Supported primitive types
const (
	PrimitiveTypeTriangles PrimitiveType = 1 << 0
	PrimitiveTypePoints    PrimitiveType = 1 << 1
	PrimitiveTypeLines     PrimitiveType = 1 << 2
)

// Mesh is an interface which wraps handling of geometry.
type Mesh interface {
	SetPrimitiveType(PrimitiveType)

	SetPositions(positions []float32)
	SetNormals(normals []float32)
	SetTangents(tangents []float32)
	SetBitangents(tangents []float32)
	SetTextureCoordinates(size int32, coordinates []float32)
	SetIndices(indices []uint16)
	SetInstanceCount(count int)
	SetModelMatrices(matrices []float32)

	SetName(name string)
	Name() string

	Draw()

	Bounds() *AABB

	Lt(Mesh) bool
	Gt(Mesh) bool
}

// IMGUIMesh is an interface which wraps a Mesh used for IMGUI primitives.
type IMGUIMesh interface {
	Mesh
}

// NewScreenQuadMesh returns a mesh of the specified size, axis aligned to be drawn
// by an orthographic projection camera.
func NewScreenQuadMesh(width float32, height float32) Mesh {
	positions := []float32{
		0.0 * width, 0.0 * height, 0.0,
		0.0 * width, 1.0 * height, 0.0,
		1.0 * width, 1.0 * height, 0.0,
		1.0 * width, 0.0 * height, 0.0}

	tcoords := []float32{
		0.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 1.0, 0.0,
		0.0, 1.0, 0.0}

	indices := []uint16{0, 1, 2, 2, 3, 0}

	m := renderSystem.NewMesh()
	m.SetPrimitiveType(PrimitiveTypeTriangles)
	m.SetPositions(positions)
	m.SetTextureCoordinates(3, tcoords)
	m.SetIndices(indices)
	m.SetName(fmt.Sprintf("ScreenQuad-%fx%f", width, height))
	return m
}

// NewAABBMesh returns a normalized cube centered at the origin. This is
// used to draw bounding boxes by translating and scaling it according to node bounds.
func NewAABBMesh() Mesh {
	positions := []float32{
		-0.5, -0.5, -0.5,
		+0.5, -0.5, -0.5,
		+0.5, +0.5, -0.5,
		-0.5, +0.5, -0.5,
		-0.5, -0.5, +0.5,
		+0.5, -0.5, +0.5,
		+0.5, +0.5, +0.5,
		-0.5, +0.5, +0.5,
	}

	indices := []uint16{
		0, 1,
		1, 2,
		2, 3,
		3, 0,
		4, 5,
		5, 6,
		6, 7,
		7, 4,
		0, 4,
		1, 5,
		2, 6,
		3, 7}

	m := renderSystem.NewMesh()
	m.SetPrimitiveType(PrimitiveTypeLines)
	m.SetPositions(positions)
	m.SetIndices(indices)
	m.SetName("AABB")
	return m
}
