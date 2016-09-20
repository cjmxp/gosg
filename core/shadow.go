package core

import (
	"math"

	"github.com/fcvarela/gosg/protos"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

// Shadower is an interface which wraps logic to implement shadowing of a light
type Shadower interface {
	// RenderTarget returns the shadower's render target.
	RenderTarget() RenderTarget

	// Render calls the shadower render implementation by assing a light and a list of nodes.
	RenderStage(light *Light, stateBuckets map[*protos.State][]*Node) RenderStage
}

// ShadowMap is a utility implementation of the Shadower interface which renders shadows by using a shadow map.
type ShadowMap struct {
	size   uint16
	camera *Camera
}

// NewShadowMap returns a new ShadowMap
func NewShadowMap(size uint16) *ShadowMap {
	rt := renderSystem.NewRenderTarget(uint32(size), uint32(size), 0, 10)
	c := NewCamera("ShadowCamera", OrthographicProjection)

	c.SetRenderTarget(rt)
	c.SetViewport(mgl32.Vec4{0.0, 0.0, float32(size), float32(size)})
	c.SetAutoReshape(false)
	c.SetRenderTechnique(nil)
	return &ShadowMap{size, c}
}

// RenderTarget implements the Shadower interface
func (s *ShadowMap) RenderTarget() RenderTarget {
	return s.camera.RenderTarget()
}

// Render implements the Shadower interface
func (s *ShadowMap) RenderStage(light *Light, stateBuckets map[*protos.State][]*Node) (out RenderStage) {
	/*
		1-find all objects that are inside the current camera frustum
		2-find minimal aa bounding box that encloses them all
		3-transform corners of that bounding box to the light's space (using light's view matrix)
		4-find aa bounding box in light's space of the transformed (now obb) bounding box
		5-this aa bounding box is your directional light's orthographic frustum.
	*/

	// 1-find all objects that are inside the current camera frustum
	// 2-find minimal aa bounding box that encloses them all
	nodesBoundsWorld := NewAABB()
	for _, nodes := range stateBuckets {
		for n := range nodes {
			nodesBoundsWorld.ExtendWithBox(nodes[n].worldBounds)

			// for now, hack the shadow texture sampler
			nodes[n].materialData.SetTexture("shadowTex", s.camera.renderTarget.ColorTexture(0))
		}
	}

	// compute lightcamera viewmatrix
	lightPos64 := mgl64.Vec3{float64(light.Block.Position.X()), float64(light.Block.Position.Y()), float64(light.Block.Position.Z())}
	s.camera.viewMatrix = mgl64.LookAtV(lightPos64, mgl64.Vec3{0.0, 0.0, 0.0}, mgl64.Vec3{0.0, 1.0, 0.0})

	// 3-transform corners of that bounding box to the light's space (using light's view matrix)
	// 4-find aa bounding box in light's space of the transformed (now obb) bounding box
	nodesBoundsLight := nodesBoundsWorld.Transformed(s.camera.viewMatrix)

	// 5-this aa bounding box is your directional light's orthographic frustum. except we want integer increments
	worldUnitsPerTexel := nodesBoundsLight.Max().Sub(nodesBoundsLight.Min()).Mul(1.0 / float64(s.size))
	projMinX := math.Floor(nodesBoundsLight.Min().X()/worldUnitsPerTexel.X()) * worldUnitsPerTexel.X()
	projMaxX := math.Floor(nodesBoundsLight.Max().X()/worldUnitsPerTexel.X()) * worldUnitsPerTexel.X()
	projMinY := math.Floor(nodesBoundsLight.Min().Y()/worldUnitsPerTexel.Y()) * worldUnitsPerTexel.Y()
	projMaxY := math.Floor(nodesBoundsLight.Max().Y()/worldUnitsPerTexel.Y()) * worldUnitsPerTexel.Y()

	s.camera.projectionMatrix = mgl64.Ortho(
		projMinX, projMaxX,
		projMinY, projMaxY,
		-nodesBoundsLight.Max().Z(),
		-nodesBoundsLight.Min().Z())

	vpmatrix := s.camera.projectionMatrix.Mul4(s.camera.viewMatrix)
	biasvpmatrix := mgl64.Mat4FromCols(
		mgl64.Vec4{0.5, 0.0, 0.0, 0.0},
		mgl64.Vec4{0.0, 0.5, 0.0, 0.0},
		mgl64.Vec4{0.0, 0.0, 0.5, 0.0},
		mgl64.Vec4{0.5, 0.5, 0.5, 1.0}).Mul4(vpmatrix)
	light.Block.VPMatrix = Mat4DoubleToFloat(biasvpmatrix)

	// set camera constants
	s.camera.constants.SetData(s.camera.projectionMatrix, s.camera.viewMatrix, nil)

	// create pass per bucket, opaque is default
	for state, nodeBucket := range stateBuckets {
		if state.Blending == true {
			continue
		}

		out.Passes = append(out.Passes, RenderPass{
			State: resourceManager.State("shadow"),
			Name:  "ShadowPass",
			Nodes: nodeBucket,
		})
	}

	out.Camera = s.camera
	out.Name = "ShadowStage"

	return
}
