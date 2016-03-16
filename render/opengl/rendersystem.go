package opengl

import (
	"time"

	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/glog"
)

// RenderSystem implements the core.RenderSystem interface
type RenderSystem struct {
	lastPassMeshCount int
	timerQuery        uint32
	waitingResult     bool
	frameDuration     time.Duration
	cpuDuration       time.Duration
	startTime         float64
}

func init() {
	core.SetRenderSystem(New())
}

// New returns a new RenderSystem
func New() *RenderSystem {
	r := RenderSystem{waitingResult: false}
	return &r
}

// Start implements the core.RenderSystem interface
func (r *RenderSystem) Start() {
	// load clear material
	clearMaterial = core.GetResourceManager().Material("clear")

	// set it as the active one, and force bind it
	currentMaterial = clearMaterial
	bindMaterialState(nil, clearMaterial, true)

	// create timers
	gl.GenQueries(1, &r.timerQuery)
}

// StartTimer implements the core.RenderSystem interface
func (r *RenderSystem) StartTimer() {
	if r.waitingResult {
		return
	}

	gl.BeginQuery(gl.TIME_ELAPSED, r.timerQuery)
	r.startTime = core.GetTimerManager().GetTime()
}

// EndTimer implements the core.RenderSystem interface
func (r *RenderSystem) EndTimer() {
	if !r.waitingResult {
		gl.EndQuery(gl.TIME_ELAPSED)
		cpuDuration := uint64((core.GetTimerManager().GetTime() - r.startTime) * 1E9)
		r.cpuDuration = time.Duration(cpuDuration)
	}

	r.waitingResult = true
	var timerResultReady int32
	gl.GetQueryObjectiv(r.timerQuery, gl.QUERY_RESULT_AVAILABLE, &timerResultReady)
	if timerResultReady == 0 {
		return
	}

	var timerResult uint32
	gl.GetQueryObjectuiv(r.timerQuery, gl.QUERY_RESULT, &timerResult)

	r.waitingResult = false
	// both gl and durations use nanoseconds as base unit
	r.frameDuration = time.Duration(timerResult)
}

// RenderTime implements the core.RenderSystem interface
func (r *RenderSystem) RenderTime() time.Duration {
	return r.frameDuration
}

// CPUTime implements the core.RenderSystem interface
func (r *RenderSystem) CPUTime() time.Duration {
	return r.cpuDuration
}

// Stop implements the core.RenderSystem interface
func (r *RenderSystem) Stop() {
	// nothing
	glog.Info("Stopping")
}

// PrepareViewport implements the core.RenderSystem interface
func (r *RenderSystem) PrepareRenderTarget(c *core.Camera) {
	// bind specific render target
	if c.RenderTarget() != nil {
		c.RenderTarget().(*RenderTarget).bind()
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
		bindMaterialState(nil, clearMaterial, true)
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
		r.PrepareRenderTarget(stage.Camera)

		for _, pass := range stage.Passes {
			program := bindMaterialState(stage.Camera.Constants().UniformBuffer(), pass.Material, false)

			var renderBatches []RenderBatch
			var lastBatchIndex = 0
			for i := 0; i < len(pass.Nodes); i++ {
				if i == 0 {
					continue
				}

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

	// bind the textures for this batch
	bindTextures(program, nodes[0].MaterialData())

	//fixme:  build uniform buffer
	//bindUniforms(program, nodes[0].MaterialData())

	// build transform attribute buffer
	mesh := nodes[0].Mesh()
	matrixBuckets := make([]float32, 0)
	for _, n := range nodes {
		transform64 := n.WorldTransform()
		transform32 := core.Mat4DoubleToFloat(transform64)
		matrixBuckets = append(matrixBuckets, transform32[0:16]...)
	}

	mesh.SetInstanceCount(len(nodes))
	mesh.SetModelMatrices(matrixBuckets)
	mesh.Draw()
}
