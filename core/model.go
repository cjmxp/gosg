package core

import (
	"fmt"
	"path/filepath"

	"encoding/binary"
	"math"

	"github.com/fcvarela/gosg/protos"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

// LoadModel parses model data from a raw resource and returns a node ready
// to insert into the screnegraph
func LoadModel(name string, res []byte) *Node {
	var model = &protos.Model{}
	if err := proto.Unmarshal(res, model); err != nil {
		glog.Fatalln(err)
	}

	basename := filepath.Base(name)
	parentNode := NewNode(basename)
	for i := 0; i < len(model.Meshes); i++ {
		node := NewNode(basename + fmt.Sprintf("-%d", i))

		node.material = resourceManager.Material(model.Meshes[i].Material)

		// get textures
		if len(model.Meshes[i].DiffuseMap) > 0 {
			node.MaterialData().SetTexture("diffuseTex", renderSystem.NewTexture(model.Meshes[i].DiffuseMap))
		}

		if len(model.Meshes[i].NormalMap) > 0 {
			node.MaterialData().SetTexture("normalTex", renderSystem.NewTexture(model.Meshes[i].NormalMap))
		}

		// set mesh data
		mesh := renderSystem.NewMesh()
		mesh.SetName(node.name)
		mesh.SetPositions(bytesToFloat(model.Meshes[i].Positions))
		mesh.SetNormals(bytesToFloat(model.Meshes[i].Normals))
		mesh.SetTangents(bytesToFloat(model.Meshes[i].Tangents))
		mesh.SetBitangents(bytesToFloat(model.Meshes[i].Bitangents))
		mesh.SetTextureCoordinates(bytesToFloat(model.Meshes[i].Tcoords))
		mesh.SetIndices(bytesToShort(model.Meshes[i].Indices))
		mesh.SetPrimitiveType(PrimitiveTypeTriangles)

		node.SetMesh(mesh)
		parentNode.AddChild(node)
	}

	return parentNode
}

func bytesToFloat(b []byte) []float32 {
	data := make([]float32, len(b)/4)
	for i := range data {
		data[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[i*4 : (i+1)*4]))
	}
	return data
}

func bytesToShort(b []byte) []uint16 {
	data := make([]uint16, len(b)/2)
	for i := range data {
		data[i] = uint16(binary.LittleEndian.Uint16(b[i*2 : (i+1)*2]))
	}
	return data
}
