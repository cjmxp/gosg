package core

import "github.com/go-gl/mathgl/mgl32"

// LightBlock holds a light's properties. It is embedded in a sceneblock and
// passed to every program.
type LightBlock struct {
	VPMatrix mgl32.Mat4
	Position mgl32.Vec4
	Ambient  mgl32.Vec4
	Diffuse  mgl32.Vec4
	Specular mgl32.Vec4
}

// Light represents a light. It contains a properties block and an optional shadower.
type Light struct {
	Block    LightBlock
	Shadower Shadower
}

// LightExtractor is an interface which extracts a light from a node and adds it to a bucket.
type LightExtractor interface {
	// Run extracts light from a node and adds it to lightBucket.
	Run(node *Node, lightBucket *[]*Light)
}

// DefaultLightExtractor is a LightExtractor which adds a nodes's light to the bucket if the node is active.
type DefaultLightExtractor struct{}

// Run implements the LightExtractor interface
func (lc *DefaultLightExtractor) Run(node *Node, lightBucket *[]*Light) {
	if node.active == false {
		return
	}

	// the default implementation is to add ourselves to the bucket
	if node.light != nil {
		// set xyz from node transform
		lPos := node.WorldPosition().Vec4(float64(node.light.Block.Position.W()))
		node.light.Block.Position = Vec4DoubleToFloat(lPos)
		*lightBucket = append(*lightBucket, node.light)
	}

	// and then recurse
	for _, c := range node.children {
		c.lightExtractor.Run(c, lightBucket)
	}
}
