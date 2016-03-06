package filesystem

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fcvarela/gosg/core"
	"github.com/golang/glog"
)

// ResourceSystem implements the resource system interface
type ResourceSystem struct {
	paths map[string]string
}

var (
	basePath = flag.String("data", "./data", "Data directory")
)

func init() {
	flag.Parse()
	core.GetResourceManager().SetSystem(New(*basePath))
}

// New returns a new ResourceSystem
func New(basePath string) *ResourceSystem {
	var bp string

	if runtime.GOOS == "darwin" && strings.HasSuffix(filepath.Dir(os.Args[0]), "MacOS") {
		glog.Info("Looking for data directory in same folder")
		path, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			glog.Fatalf("Could not create data path from provided: %s\n", basePath)
		}
		bp = filepath.Join(path, basePath)
	} else {
		path, err := filepath.Abs(basePath)
		if err != nil {
			glog.Fatalf("Could not create data path from provided: %s\n", basePath)
		}
		bp = path
	}

	paths := make(map[string]string)
	paths["base"] = bp
	paths["programs"] = filepath.Join(bp, "programs")
	paths["models"] = filepath.Join(bp, "models")
	paths["textures"] = filepath.Join(bp, "textures")

	r := ResourceSystem{paths: paths}

	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			glog.Fatalf("No such file or directory: %v\n", p)
		}
	}

	return &r
}

// Start implements the core.ResourceSystem interface
func (r *ResourceSystem) Start() {
	glog.Info("Starting")
}

// Stop implements the core.ResourceSystem interface
func (r *ResourceSystem) Stop() {
	glog.Info("Stopping")
}

func (r *ResourceSystem) resourceWithFullpath(fullpath string) []byte {
	data, e := ioutil.ReadFile(fullpath)
	if e != nil {
		glog.Fatalf("Could not read file: %v\n", e)
	}
	return data
}

// Model implements the core.ResourceSystem interface
func (r *ResourceSystem) Model(filename string) []byte {
	fullpath := filepath.Join(r.paths["models"], filename)
	res := r.resourceWithFullpath(fullpath)
	return res
}

// Texture implements the core.ResourceSystem interface
func (r *ResourceSystem) Texture(filename string) []byte {
	fullpath := filepath.Join(r.paths["textures"], filename)
	res := r.resourceWithFullpath(fullpath)
	return res
}

// Program implements the core.ResourceSystem interface
func (r *ResourceSystem) Program(name string) []byte {
	fullpath := filepath.Join(r.paths["programs"], name) + "." + core.GetRenderSystem().ProgramExtension()
	res := r.resourceWithFullpath(fullpath)
	return res
}

// ProgramData implements the core.ResourceSystem interface
func (r *ResourceSystem) ProgramData(name string) []byte {
	fullpath := filepath.Join(r.paths["programs"], name)
	res := r.resourceWithFullpath(fullpath)
	return res
}
