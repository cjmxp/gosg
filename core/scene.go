package core

import (
	"runtime"
	"sort"

	"github.com/golang/glog"
)

// Scene represents a scenegraph and contains information about how it should be composed with other scenes on
// a scene stack. Scenes are not meant to be wrapped by users, but to be data configured for the expected behaviour.
type Scene struct {
	name string
	root *Node

	// should the scenemanager call update and draw this scene
	active bool

	// should this scene trigger a cursor hide/display toggle?
	displaysCursor bool

	cameraList []*Camera

	// per camera draw lists
	drawables map[string][]*Node

	// per scene lights list
	lights []*Light
}

func deleteScene(s *Scene) {
	glog.Info("Scene finalizer started: ", s.name)

	s.root.RemoveChildren()
	s.root = nil
	s.drawables = nil

	glog.Info("Scene finalizer finished: ", s.name)
}

// NewScene returns a new scene.
func NewScene(name string) *Scene {
	s := Scene{}

	s.name = name
	s.active = true
	s.cameraList = make([]*Camera, 0)
	s.drawables = make(map[string][]*Node)
	s.lights = make([]*Light, 0)

	runtime.SetFinalizer(&s, deleteScene)

	return &s
}

// Root returns the scene's root node
func (s *Scene) Root() *Node {
	return s.root
}

// SetRoot returns the scene's root node
func (s *Scene) SetRoot(root *Node) {
	s.root = root
}

// Name returns the scene's name
func (s *Scene) Name() string {
	return s.name
}

// SetActive sets the 'active' state of this scene
func (s *Scene) SetActive(active bool) {
	s.active = active
}

// Active returns whether this scene is active or not
func (s *Scene) Active() bool {
	return s.active
}

// AddCamera adds a camera to the scene by attaching it to the given node.
func (s *Scene) AddCamera(node *Node, camera *Camera) {
	node.AddChild(camera.node)

	s.cameraList = append(s.cameraList, camera)

	// resort camera list by renderorder
	if len(s.cameraList) > 1 {
		sort.Sort(CamerasByRenderOrder(s.cameraList))
	}
}

func (s *Scene) update(dt float64) {
	// physics update
	var physicsNodes []*Node
	if s.root.physicsComponent != nil {
		s.root.physicsComponent.Run(s.root, &physicsNodes)
	}

	physicsSystem.Update(dt, physicsNodes)

	// update transforms and bounds
	s.root.update(dt)
}

func (s *Scene) cull() {
	for _, c := range s.cameraList {
		c.Reshape(windowSystem.WindowSize())
	}

	s.lights = nil
	if s.root.lightExtractor != nil {
		s.root.lightExtractor.Run(s.root, &s.lights)
	}

	for _, camera := range s.cameraList {
		camera.constants.SetLights(s.lights)
		sceneRoot := camera.Scene()

		var bucket []*Node
		sceneRoot.CullComponent().Run(s, camera, sceneRoot, &bucket)
		s.drawables[camera.name] = bucket
	}
}

func (s *Scene) draw() {
	for _, camera := range s.cameraList {
		if camera.projectionType == PerspectiveProjection {
			for _, light := range s.lights {
				if light.Shadower != nil {
					light.Shadower.Render(light, s.drawables[camera.name])
				}
			}
		}

		camera.PrepareViewport()
		camera.Render(s.drawables[camera.name])
	}
}

// SetDisplaysCursor sets whether this scene wants the cursor to be hidden or not
func (s *Scene) SetDisplaysCursor(displaysCursor bool) {
	s.displaysCursor = displaysCursor
}

func (s *Scene) movedToFront() {
	if s.displaysCursor {
		windowSystem.SetCursorVisible(true)
	} else {
		windowSystem.SetCursorVisible(false)
	}
}
