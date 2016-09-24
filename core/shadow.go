package core

import (
	"math"

	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

// Shadower is an interface which wraps logic to implement shadowing of a light
type Shadower interface {
	// Textures returns the shadow textures used by this shadower
	Textures() []Texture

	// Render calls the shadower render implementation by assing a light and a scene camera.
	RenderStages(light *Light, camera *Camera) []RenderStage
}

// ShadowMap is a utility implementation of the Shadower interface which renders shadows by using a cascading shadow map.
type ShadowMap struct {
	size     uint16
	cameras  []*Camera
	textures []Texture
}

const (
	numCascades = 3
	maxCascades = 10
)

// NewShadowMap returns a new ShadowMap
func NewShadowMap(size uint16) *ShadowMap {
	shadowMap := &ShadowMap{size, make([]*Camera, numCascades), make([]Texture, numCascades)}
	for i := 0; i < numCascades; i++ {
		rt := renderSystem.NewRenderTarget(uint32(size), uint32(size), false, 1)
		c := NewCamera("ShadowCamera", OrthographicProjection)
		c.SetRenderTarget(rt)
		c.SetViewport(mgl32.Vec4{0.0, 0.0, float32(size), float32(size)})
		c.SetAutoReshape(false)
		c.SetRenderTechnique(nil)
		shadowMap.cameras[i] = c
		shadowMap.textures[i] = rt.ColorTexture(0)
	}
	return shadowMap
}

// RenderTarget implements the Shadower interface
func (s *ShadowMap) Textures() []Texture {
	return s.textures
}

func (s *ShadowMap) cascadeRenderStage(cascade int, light *Light, camera *Camera) (out RenderStage) {
	/*
		1-find all objects that are inside the current camera frustum
		2-find minimal aa bounding box that encloses them all
		3-transform corners of that bounding box to the light's space (using light's view matrix)
		4-find aa bounding box in light's space of the transformed (now obb) bounding box
		5-this aa bounding box is your directional light's orthographic frustum.
	*/

	// compute lightcamera viewmatrix
	lightPos64 := mgl64.Vec3{float64(light.Block.Position.X()), float64(light.Block.Position.Y()), float64(light.Block.Position.Z())}
	s.cameras[cascade].viewMatrix = mgl64.LookAtV(lightPos64, mgl64.Vec3{0.0, 0.0, 0.0}, mgl64.Vec3{0.0, 1.0, 0.0})

	// 3-transform corners of that bounding box to the light's space (using light's view matrix)
	// 4-find aa bounding box in light's space of the transformed (now obb) bounding box
	nodesBoundsLight := camera.cascadingAABBS[cascade].Transformed(s.cameras[cascade].viewMatrix)

	// 5-this aa bounding box is your directional light's orthographic frustum. except we want integer increments
	worldUnitsPerTexel := nodesBoundsLight.Max().Sub(nodesBoundsLight.Min()).Mul(1.0 / float64(s.size))
	projMinX := math.Floor(nodesBoundsLight.Min().X()/worldUnitsPerTexel.X()) * worldUnitsPerTexel.X()
	projMaxX := math.Floor(nodesBoundsLight.Max().X()/worldUnitsPerTexel.X()) * worldUnitsPerTexel.X()
	projMinY := math.Floor(nodesBoundsLight.Min().Y()/worldUnitsPerTexel.Y()) * worldUnitsPerTexel.Y()
	projMaxY := math.Floor(nodesBoundsLight.Max().Y()/worldUnitsPerTexel.Y()) * worldUnitsPerTexel.Y()

	s.cameras[cascade].projectionMatrix = mgl64.Ortho(
		projMinX, projMaxX,
		projMinY, projMaxY,
		-nodesBoundsLight.Max().Z(),
		-nodesBoundsLight.Min().Z())

	vpmatrix := s.cameras[cascade].projectionMatrix.Mul4(s.cameras[cascade].viewMatrix)
	biasvpmatrix := mgl64.Mat4FromCols(
		mgl64.Vec4{0.5, 0.0, 0.0, 0.0},
		mgl64.Vec4{0.0, 0.5, 0.0, 0.0},
		mgl64.Vec4{0.0, 0.0, 0.5, 0.0},
		mgl64.Vec4{0.5, 0.5, 0.5, 1.0}).Mul4(vpmatrix)

	// set light block
	light.Block.ZCuts[cascade] = mgl32.Vec4{float32(camera.cascadingZCuts[cascade]), 0.0, 0.0, 0.0}
	light.Block.VPMatrix[cascade] = Mat4DoubleToFloat(biasvpmatrix)

	// set camera constants
	s.cameras[cascade].constants.SetData(s.cameras[cascade].projectionMatrix, s.cameras[cascade].viewMatrix, nil)

	// create a single stage now
	out.Camera = s.cameras[cascade]
	out.Name = fmt.Sprintf("ShadowStageCascade%d", cascade)

	// create pass per bucket, opaque is default
	for state, nodeBucket := range camera.stateBuckets {
		if state.Blending == true {
			continue
		}

		for _, n := range nodeBucket {
			n.materialData.SetTexture(fmt.Sprintf("shadowTex%d", cascade), s.textures[cascade])
		}

		out.Passes = append(out.Passes, RenderPass{
			State: resourceManager.State("shadow"),
			Name:  "ShadowPass",
			Nodes: nodeBucket,
		})
	}

	return out
}

// Render implements the Shadower interface
func (s *ShadowMap) RenderStages(light *Light, cam *Camera) (out []RenderStage) {
	for c := 0; c < numCascades; c++ {
		out = append(out, s.cascadeRenderStage(c, light, cam))
	}

	return
}
