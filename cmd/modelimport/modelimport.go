package main

// make sure you're linking to libassimp.a, explicit static doesn't work on macos

// #cgo linux LDFLAGS: -Bstatic -lassimp -lstdc++
// #cgo darwin LDFLAGS: -Bstatic -lassimp -lstdc++ -lz
// #cgo pkg-config: assimp
// #cgo CFLAGS: -Wno-attributes
// #cgo CXXFLAGS: -Wno-attributes
// #cgo darwin CXXFLAGS: -I/usr/local/include -I/usr/local/include/assimp
// #cgo darwin CFLAGS: -I/usr/local/include -I/usr/local/include/assimp
// #include "modelimport.h"
import "C"
import (
	"bytes"
	"encoding/gob"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/golang/glog"

	"github.com/fcvarela/gosg/core"
)

// placeholder for memcpy to copy to
var (
	placeholder = strings.Repeat("A", int(512)+1)
	modelfile   = flag.String("modelfile", "", "Input model file")
)

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
}

func main() {
	scene := importModelfile(*modelfile)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(scene)
	if err != nil {
		glog.Fatal("Could not encode: ", err)
	}

	file, err := os.Create(*modelfile + ".model")
	if err != nil {
		glog.Fatal("Could not open file for writing: ", err)
	}
	defer file.Close()

	file.Write(buf.Bytes())
	glog.Infoln("Done")
}

func importModelfile(filename string) *core.Model {
	opts := C.aiProcessPreset_TargetRealtime_MaxQuality

	opts = opts | C.aiProcess_CalcTangentSpace
	opts = opts | C.aiProcess_FixInfacingNormals
	opts = opts | C.aiProcess_FlipUVs
	opts = opts | C.aiProcess_GenSmoothNormals
	opts = opts | C.aiProcess_ImproveCacheLocality
	opts = opts | C.aiProcess_TransformUVCoords
	opts = opts | C.aiProcess_Triangulate
	opts = opts | C.aiProcess_ValidateDataStructure
	opts = opts | C.aiProcess_OptimizeGraph
	opts = opts | C.aiProcess_OptimizeMeshes

	// run import from file
	aiscene := C.aiImportFile(C.CString(filename), C.uint(opts))
	if aiscene == nil {
		str := C.aiGetErrorString()
		glog.Fatal("Could not import file: ", C.GoString(str))
	}

	// create scene
	s := core.Model{}
	s.Meshes = make([]core.ModelMesh, 0)

	for meshIdx := 0; meshIdx < int(C.get_mesh_count(aiscene)); meshIdx++ {
		// create mesh and set its properties
		m := core.ModelMesh{}

		// set name
		meshName := C.CString(placeholder)
		C.get_mesh_name(aiscene, C.int(meshIdx), meshName)
		m.Name = C.GoString(meshName)
		glog.Infoln("Got a mesh: ", m.Name)

		// get vertex count
		vertexcount := int(C.get_vertex_count(aiscene, C.int(meshIdx))) * 3

		// get positions
		glog.Infoln("Copying positions")
		positions := (*[1 << 30]C.float)(unsafe.Pointer(C.get_positions(aiscene, C.int(meshIdx))))
		m.Positions = make([]float32, vertexcount)
		for i := 0; i < vertexcount; i++ {
			m.Positions[i] = float32(positions[i])
		}

		// get normals
		glog.Infoln("Copying normals")
		normals := (*[1 << 30]C.float)(unsafe.Pointer(C.get_normals(aiscene, C.int(meshIdx))))
		m.Normals = make([]float32, vertexcount)
		for i := 0; i < vertexcount; i++ {
			m.Normals[i] = float32(normals[i])
		}

		// get tangents
		glog.Infoln("Copying tangents")
		tangents := (*[1 << 30]C.float)(unsafe.Pointer(C.get_tangents(aiscene, C.int(meshIdx))))
		if tangents != nil {
			m.Tangents = make([]float32, vertexcount)
			for i := 0; i < vertexcount; i++ {
				m.Tangents[i] = float32(tangents[i])
			}
		} else {
			glog.Error("Model has no tangents and importer failed to generate them")
			continue
		}

		// get bitangents
		glog.Infoln("Copying bitangents")
		bitangents := (*[1 << 30]C.float)(unsafe.Pointer(C.get_bitangents(aiscene, C.int(meshIdx))))
		m.Bitangents = make([]float32, vertexcount)
		for i := 0; i < vertexcount; i++ {
			m.Bitangents[i] = float32(bitangents[i])
		}

		// get texturecoords
		glog.Infoln("Copying texture coordinates")
		texturecoords := (*[1 << 30]C.float)(unsafe.Pointer(C.get_texturecoords(aiscene, C.int(meshIdx))))
		m.TextureCoords = make([]float32, vertexcount)
		for i := 0; i < vertexcount; i++ {
			m.TextureCoords[i] = float32(texturecoords[i])
		}

		// get indexes
		glog.Infoln("Copying indexes")
		indexcount := C.get_indexcount(aiscene, C.int(meshIdx))
		indices := (*[1 << 30]C.uint)(unsafe.Pointer(C.get_indices(aiscene, C.int(meshIdx))))
		m.Indices = make([]uint16, indexcount)
		for i := 0; i < int(indexcount); i++ {
			m.Indices[i] = uint16(indices[i])
		}

		// get diffuse/normal textures
		glog.Infoln("Copying textures")
		diffusePath := C.CString(placeholder)
		normalPath := C.CString(placeholder)
		C.get_mesh_maps(aiscene, C.int(meshIdx), diffusePath, normalPath)

		// load them
		if C.GoString(diffusePath) != "" {
			glog.Infoln("Copying diffuse texture: ", C.GoString(diffusePath))
			dpath := filepath.Join(filepath.Dir(filename), C.GoString(diffusePath))
			diffuseTextureFile, e := ioutil.ReadFile(dpath)
			if e != nil {
				glog.Fatal("Could not read file: ", e)
			}
			m.DiffuseTexture = diffuseTextureFile
		}

		if C.GoString(normalPath) != "" {
			glog.Infoln("Copying normal texture: ", C.GoString(normalPath))
			npath := filepath.Join(filepath.Dir(filename), C.GoString(normalPath))
			normalTextureFile, e := ioutil.ReadFile(npath)
			if e != nil {
				glog.Fatal("Could not read file: ", e)
			}
			m.NormalTexture = normalTextureFile
		}

		s.Meshes = append(s.Meshes, m)

		// set texture wrap mode
		m.WrapMode = int(C.get_mesh_wrapmode(aiscene, C.int(meshIdx)))

		// cleanup
		glog.Infoln("Cleaning up")
		C.free(unsafe.Pointer(diffusePath))
		C.free(unsafe.Pointer(normalPath))
		C.free(unsafe.Pointer(meshName))
		C.free(unsafe.Pointer(indices))
	}

	glog.Infoln("Releasing scene")
	// release everything
	C.aiReleaseImport(aiscene)
	return &s
}
