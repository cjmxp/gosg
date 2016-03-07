package demoapp

import (
	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func makeDemo2Scene() *core.Scene {
	s := core.NewScene("Demo2")

	camera := core.NewCamera("GeometryPassCamera", core.PerspectiveProjection)
	camera.SetAutoReshape(true)
	camera.SetVerticalFieldOfView(60.0)
	camera.SetClearColor(mgl32.Vec4{0.0, 0.0, 0.0, 0.0})
	camera.SetClearMode(core.ClearColor | core.ClearDepth)
	camera.SetClipDistance(mgl64.Vec2{1.0, 1000.0})
	camera.Node().SetInputComponent(core.NewMouseCameraInputComponent(100.0))
	camera.SetRenderOrder(0)

	node := core.NewNode("GeometryRoot")
	//physicsSystem := core.GetPhysicsSystem()

	for i := -5; i < 5; i++ {
		for j := -5; j < 5; j++ {
			randomVec := mgl64.Vec3{float64(i) * 9.96 * 2.0, float64(j) * 9.96 * 2.0, 0.0}

			// load model
			f16 := core.GetResourceManager().InstancedModel("f16.model")
			f16.Translate(randomVec)
			f16.Rotate(float64(i), randomVec)

			// add rigid body
			//ss := physicsSystem.NewSphereShape(9.96)
			//rb := physicsSystem.CreateRigidBody(12.150, ss)
			//f16.SetRigidBody(rb)

			// add body to physics simulation
			//physicsSystem.AddRigidBody(rb)

			// add model to scenegraph
			node.AddChild(f16)
		}
	}

	lightNode1 := core.NewNode("Light1")
	lightNode2 := core.NewNode("Light2")
	lightNode3 := core.NewNode("Light3")

	lightNode1.Translate(mgl64.Vec3{-1000.0, 0.0, 0.0})
	lightNode2.Translate(mgl64.Vec3{+1000.0, 0.0, 0.0})
	lightNode3.Translate(mgl64.Vec3{0.0, 0.0, +1000.0})

	light1 := &core.Light{
		Block: core.LightBlock{
			// position is only used to determine light type (w component)
			Position: mgl32.Vec4{0.0, 0.0, 0.0, 1.0},
			Ambient:  mgl32.Vec4{0.0, 0.0, 0.0, 1.0},
			Diffuse:  mgl32.Vec4{0.0, 0.0, 1.0, 1.0},
			Specular: mgl32.Vec4{0.0, 0.0, 1.0, 1.0},
		},
		Shadower: nil,
	}

	light2 := &core.Light{
		Block: core.LightBlock{
			// position is only used to determine light type (w component)
			Position: mgl32.Vec4{0.0, 0.0, 0.0, 1.0},
			Ambient:  mgl32.Vec4{0.0, 0.0, 0.0, 1.0},
			Diffuse:  mgl32.Vec4{1.0, 0.0, 0.0, 1.0},
			Specular: mgl32.Vec4{1.0, 0.0, 0.0, 1.0},
		},
		Shadower: nil,
	}

	light3 := &core.Light{
		Block: core.LightBlock{
			// position is only used to determine light type (w component)
			Position: mgl32.Vec4{0.0, 0.0, 0.0, 1.0},
			Ambient:  mgl32.Vec4{0.0, 0.0, 0.0, 1.0},
			Diffuse:  mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
			Specular: mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
		},
		Shadower: nil,
	}

	lightNode1.SetLight(light1)
	lightNode2.SetLight(light2)
	lightNode2.SetLight(light3)

	node.AddChild(lightNode1)
	node.AddChild(lightNode2)
	node.AddChild(lightNode3)

	// add ground plane
	ps := core.GetPhysicsSystem().NewStaticPlaneShape(mgl64.Vec3{0.0, 1.0, 0.0}, -120.0)
	pb := core.GetPhysicsSystem().CreateRigidBody(0.0, ps)
	core.GetPhysicsSystem().AddRigidBody(pb)

	// set gravity for the whole system
	core.GetPhysicsSystem().SetGravity(mgl64.Vec3{0.0, -9.8, 0.0})

	// push camera back
	camera.Node().Translate(mgl64.Vec3{0.0, 0.0, 200.0})

	// set camera's scene
	camera.SetScene(node)

	// add node graph to scene root
	s.SetRoot(node)

	// add camera to root node
	s.AddCamera(node, camera)

	s.SetActive(true)
	s.SetDisplaysCursor(false)

	return s
}
