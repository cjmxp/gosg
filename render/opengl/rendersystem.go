package opengl

import (
	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/golang/glog"
)

// RenderSystem implements the core.RenderSystem interface
type RenderSystem struct {
	renderLog string
}

func init() {
	core.SetRenderSystem(New())
}

// New returns a new RenderSystem
func New() *RenderSystem {
	r := RenderSystem{}
	return &r
}

// RenderLog implements the core.RenderSystem interface
func (r *RenderSystem) RenderLog() string {
	return r.renderLog
}

// Start implements the core.RenderSystem interface
func (r *RenderSystem) Start() {
	// load clear material
	clearState = core.GetResourceManager().State("clear")

	// set it as the active one, and force bind it
	currentState = clearState
	bindMaterialState(nil, clearState, true)

	// generate basic mesh buffers
	sharedBuffers = newBuffers()
	imguiBuffers = newBuffers()
}

// Stop implements the core.RenderSystem interface
func (r *RenderSystem) Stop() {
	// nothing
	glog.Info("Stopping")
}

func (r *RenderSystem) prepareRenderTarget(c *core.Camera) {
	// bind specific render target
	if c.RenderTarget() != nil {
		c.RenderTarget().(*Framebuffer).bind()
	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	}

	// reset viewport
	v := c.Viewport()
	gl.Viewport(int32(v[0]), int32(v[1]), int32(v[2]), int32(v[3]))

	var clearargs uint32
	cm := c.ClearMode()

	if cm&core.ClearColor > 0 {
		clearargs = clearargs | gl.COLOR_BUFFER_BIT
		cc := c.ClearColor()
		gl.ClearColor(cc[0], cc[1], cc[2], cc[3])
	}

	if cm&core.ClearDepth > 0 {
		clearargs = clearargs | gl.DEPTH_BUFFER_BIT
		cd := c.ClearDepth()
		gl.ClearDepth(cd)
	}

	if clearargs != 0 {
		bindMaterialState(nil, clearState, false)
		gl.Clear(clearargs)
	}
}

type RenderBatch struct {
	program *Program
	nodes   []*core.Node
}

// ExecuteRenderPlan implements the core.RenderSystem interface
func (r *RenderSystem) ExecuteRenderPlan(p core.RenderPlan) {
	for _, stage := range p.Stages {
		//r.renderLog += fmt.Sprintf("RenderStage: %s\n", stage.Name)
		r.prepareRenderTarget(stage.Camera)

		for _, pass := range stage.Passes {
			//r.renderLog += fmt.Sprintf("\tRenderPass: %s\n", pass.Name)
			program := bindMaterialState(stage.Camera.Constants().UniformBuffer(), pass.State, false)

			var renderBatches []RenderBatch

			var lastBatchIndex = 0
			for i := 1; i < len(pass.Nodes); i++ {
				if breaksBatch(pass.Nodes[i].MaterialData(), pass.Nodes[i-1].MaterialData()) {
					renderBatches = append(renderBatches, RenderBatch{program, pass.Nodes[lastBatchIndex:i]})
					lastBatchIndex = i
				}
			}

			// close last batch
			renderBatches = append(renderBatches, RenderBatch{program, pass.Nodes[lastBatchIndex:]})

			for _, b := range renderBatches {
				r.renderBatch(b.program, b.nodes)
			}
		}
	}
}

func (r *RenderSystem) renderBatch(program *Program, nodes []*core.Node) {
	if len(nodes) == 0 {
		return
	}

	//r.renderLog += fmt.Sprintf("\t\tBatch: %d nodes\n", len(nodes))

	// bind the textures for this batch
	bindTextures(program, nodes[0].MaterialData())

	//fixme:  build uniform buffer
	//bindUniforms(program, nodes[0].MaterialData())

	// build transform attribute buffer
	mesh := nodes[0].Mesh()
	var matrixBuckets []float32
	for _, n := range nodes {
		transform64 := n.WorldTransform()
		transform32 := core.Mat4DoubleToFloat(transform64)
		matrixBuckets = append(matrixBuckets, transform32[0:16]...)
	}

	mesh.SetInstanceCount(len(nodes))
	mesh.SetModelMatrices(matrixBuckets)
	mesh.Draw()
}
