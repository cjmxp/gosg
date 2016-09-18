package demoapp

import (
	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func getDemo1SceneShadowTexture(s *core.Scene) core.Texture {
	geoRoot := s.Root().Children()[0]
	lightNode := geoRoot.Children()[len(geoRoot.Children())-2]
	return lightNode.Light().Shadower.RenderTarget().DepthTexture()
}

func makeGeometrySubscene() (*core.Node, *core.Camera) {
	// geometry camera
	geometryCamera := core.NewCamera("GeometryPassCamera", core.PerspectiveProjection)
	geometryCamera.SetAutoReshape(true)
	geometryCamera.SetVerticalFieldOfView(60.0)
	geometryCamera.SetClearColor(mgl32.Vec4{135.0 / 255.0, 206.0 / 255.0, 250.0 / 255.0, 0.0})
	geometryCamera.SetClearMode(core.ClearColor | core.ClearDepth)
	geometryCamera.SetClipDistance(mgl64.Vec2{1.0, 1000.0})
	geometryCamera.Node().SetInputComponent(core.NewMouseCameraInputComponent(100.0))
	geometryCamera.SetRenderOrder(0)

	geometryNode := core.NewNode("GeometryRoot")

	for i := -5; i < 5; i++ {
		for j := -5; j < 5; j++ {
			randomVec := mgl64.Vec3{float64(i) * 9.96 * 2.0, float64(j) * 9.96 * 2.0, 0.0}

			// load model
			f16 := core.GetResourceManager().Model("f16.model")
			f16.Translate(randomVec)
			f16.Rotate(float64(i), randomVec)

			// set aabb on subnodes
			//for _, c := range f16.Children() {
			//	c.State().AABB = true
			//}

			// add model to scenegraph
			geometryNode.AddChild(f16)
		}
	}

	// attach a light
	shadowMap := core.NewShadowMap(2048)
	lightNode1 := core.NewNode("Light1")
	lightNode1.Translate(mgl64.Vec3{+1000.0, 0.0, +1000.0})
	light1 := &core.Light{
		Block: core.LightBlock{
			// position is only used to determine light type (w component)
			Position: mgl32.Vec4{0.0, 0.0, 0.0, 1.0},
			Color:    mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
		},
		Shadower: shadowMap,
	}
	lightNode1.SetLight(light1)
	geometryNode.AddChild(lightNode1)
	geometryCamera.Node().Translate(mgl64.Vec3{0.0, 0.0, 50.0})
	geometryCamera.SetScene(geometryNode)

	return geometryNode, geometryCamera
}

func makeDemo1Scene() *core.Scene {
	s := core.NewScene("Demo1")
	s.SetRoot(core.NewNode("ROOT"))

	geoRoot, geoCamera := makeGeometrySubscene()

	// add geometry camera to geometry root node
	s.AddCamera(geoRoot, geoCamera)

	// visibility and cursor mode
	s.SetActive(true)

	s.Root().AddChild(geoRoot)

	return s
}
