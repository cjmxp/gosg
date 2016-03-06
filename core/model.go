package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path/filepath"

	"github.com/golang/glog"
)

// Model holds a list of ModelMesh items
type Model struct {
	Meshes []ModelMesh
}

// ModelMesh describes a 3D model
type ModelMesh struct {
	Positions      []float32
	Normals        []float32
	Tangents       []float32
	Bitangents     []float32
	TextureCoords  []float32
	BoneIds        []uint16
	BoneWeights    []float32
	Indices        []uint16
	DiffuseTexture []byte
	NormalTexture  []byte
	Name           string
	WrapMode       int
}

// LoadModel parses model data from a raw resource and returns a node ready
// to insert into the screnegraph
func LoadModel(name string, res []byte, instanced bool) *Node {
	buffer := bytes.NewBuffer(res)
	decoder := gob.NewDecoder(buffer)

	var model Model
	err := decoder.Decode(&model)
	if err != nil {
		glog.Fatalln("Error decoding model file: ", err)
	}

	basename := filepath.Base(name)
	parentNode := NewNode(basename)
	for i := 0; i < len(model.Meshes); i++ {
		node := NewNode(basename + fmt.Sprintf("-%d", i))

		// get program, support selecting program based on material property in modelfile
		//extension := filepath.Ext(basename)
		//basename_noext := basename[0 : len(basename)-len(extension)]
		if instanced {
			node.State().SetProgram(resourceManager.Program("ubershader-instanced"))
		} else {
			node.State().SetProgram(resourceManager.Program("ubershader"))
		}

		// get textures
		if len(model.Meshes[i].DiffuseTexture) > 0 {
			node.State().SetTexture(0, renderSystem.NewTexture(model.Meshes[i].DiffuseTexture))
			node.State().Uniform("diffuseTex").Set(0)
		}

		if len(model.Meshes[i].NormalTexture) > 0 {
			node.State().SetTexture(1, renderSystem.NewTexture(model.Meshes[i].NormalTexture))
			node.State().Uniform("normalTex").Set(1)
		}

		// set mesh data
		mesh := renderSystem.NewMesh()
		mesh.SetName(node.name)
		mesh.SetPositions(model.Meshes[i].Positions)
		mesh.SetNormals(model.Meshes[i].Normals)
		mesh.SetTangents(model.Meshes[i].Tangents)
		mesh.SetBitangents(model.Meshes[i].Bitangents)
		mesh.SetTextureCoordinates(3, model.Meshes[i].TextureCoords)
		mesh.SetIndices(model.Meshes[i].Indices)
		mesh.SetPrimitiveType(PrimitiveTypeTriangles)

		if instanced {
			mesh = renderSystem.NewInstancedMesh(mesh)
		}

		node.SetMesh(mesh)
		parentNode.AddChild(node)
	}

	return parentNode
}
