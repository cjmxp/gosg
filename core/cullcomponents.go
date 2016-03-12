package core

// CullComponent is an interface that wraps culling of a scenegraph.
type CullComponent interface {
	// Run culls a scenegraph node. A scene and camera are provided for visibility/frustum checks.
	// If the policy dictates the node is to be drawn, then it should be added to the nodeBucket.
	Run(*Scene, *Camera, *Node, *[]*Node)
}

// DefaultCullComponent implements a scenegraph culler. The policy for this culler is to
// mark all nodes in frustum for drawing. The node's modelMatrix state uniform is also set
// from the nodes worldtransform. This may change as we transition away from individual uniforms
// for instanced/indirect drawing.
type DefaultCullComponent struct{}

// Run implements the CullComponent interface
func (cc *DefaultCullComponent) Run(scene *Scene, camera *Camera, node *Node, nodeBucket *[]*Node) {
	if node.worldBounds.InFrustum(camera.Frustum()) == false {
		return
	}

	if node.active == false {
		return
	}

	// the default implementation is to add ourselves to the bucket
	if node.mesh != nil {
		node.materialData.Uniform("mMatrix").Set(node.worldTransform)
		*nodeBucket = append(*nodeBucket, node)
	}

	for _, c := range node.children {
		c.cullComponent.Run(scene, camera, c, nodeBucket)
	}
}

// AlwaysPassCullComponent implements a scenegraph culler. The policy for this culler is to
// mark all nodes for drawing regardless of visibility. This is useful for scenes where visibility
// check is not necessary (ie: screen quads for deferred rendering).
type AlwaysPassCullComponent struct{}

// Run implements the CullComponent interface
func (apcc *AlwaysPassCullComponent) Run(scene *Scene, camera *Camera, node *Node, nodeBucket *[]*Node) {
	// the default implementation is to add ourselves to the bucket
	if node.mesh != nil {
		node.materialData.Uniform("mMatrix").Set(node.worldTransform)
		*nodeBucket = append(*nodeBucket, node)
	}

	for _, ch := range node.children {
		ch.cullComponent.Run(scene, camera, ch, nodeBucket)
	}
}
