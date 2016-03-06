package dearimgui

// #include "gosg_imgui.h"
// #cgo windows LDFLAGS: -Wl,--allow-multiple-definition -limm32
import "C"
import (
	"unsafe"

	"github.com/fcvarela/gosg/core"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/glog"
)

func init() {
	core.SetIMGUISystem(New())
}

type IMGUISystem struct {
	texture core.Texture
}

type textureData struct {
	width   int
	height  int
	payload []byte
}

var (
	displaySize mgl32.Vec2
)

func New() *IMGUISystem {
	return &IMGUISystem{}
}

func (i *IMGUISystem) getTextureData() textureData {
	var width, height C.int
	payload := unsafe.Pointer(C.get_texture_data(&width, &height))

	bufSize := int(width) * int(height) * 4
	return textureData{int(width), int(height), C.GoBytes(payload, C.int(bufSize))}
}

func (i *IMGUISystem) Start() {
	tdata := i.getTextureData()
	i.texture = core.GetRenderSystem().NewRawTexture(tdata.width, tdata.height, tdata.payload)
	if i.texture == nil {
		glog.Fatal("Cannot set nil texture")
	}
	C.set_texture_id(i.texture.Handle())
}

func (i *IMGUISystem) Stop() {

}

func (i *IMGUISystem) Begin(name string, flags core.WindowFlags) bool {
	return int(C.begin(C.CString(name), C.int(flags))) == 1
}

func (i *IMGUISystem) End() {
	C.end()
}

func (i *IMGUISystem) CollapsingHeader(name string) bool {
	return int(C.collapsing_header(C.CString(name))) == 1
}

func (i *IMGUISystem) PlotHistogram(name string, values []float32, minScale, maxScale float32, size mgl32.Vec2) {
	C.plot_histogram(C.CString(name), (*C.float)(unsafe.Pointer(&values[0])), C.int(len(values)), C.float(minScale), C.float(maxScale), (*C.float)(unsafe.Pointer(&size[0])))

}

func (i IMGUISystem) Image(texture core.Texture, size mgl32.Vec2) {
	if texture == nil {
		glog.Fatal("Cannot draw nil texture")
	}
	C.image(texture.Handle(), (*C.float)(unsafe.Pointer(&size[0])))
}

func (i *IMGUISystem) SetNextWindowPos(pos mgl32.Vec2) {
	C.set_next_window_pos(C.float(pos[0]), C.float(pos[1]))
}

func (i *IMGUISystem) SetNextWindowSize(size mgl32.Vec2) {
	C.set_next_window_size(C.float(size[0]), C.float(size[1]))
}

func (i *IMGUISystem) StartFrame(dt float64) {
	state := core.GetInputManager().State()
	size := core.GetWindowSystem().WindowSize()
	i.SetDisplaySize(size)
	i.SetMousePosition(state.Mouse.Position.X, state.Mouse.Position.Y)
	i.SetMouseButtons(
		state.Mouse.Buttons.Active[core.MouseButton1],
		state.Mouse.Buttons.Active[core.MouseButton2],
		state.Mouse.Buttons.Active[core.MouseButton3])
	i.SetMouseScrollPosition(state.Mouse.Scroll.X, state.Mouse.Scroll.Y)

	C.set_dt(C.double(dt))
	C.frame_new()
}

// DisplaySize returns the currently set display size
func (i *IMGUISystem) DisplaySize() mgl32.Vec2 {
	return displaySize
}

func (i *IMGUISystem) SetDisplaySize(s mgl32.Vec2) {
	displaySize = s
	C.set_display_size(C.float(s[0]), C.float(s[1]))
}

func (i *IMGUISystem) SetMousePosition(x, y float64) {
	C.set_mouse_position(C.double(x), C.double(y))
}

func (i *IMGUISystem) SetMouseButtons(b0, b1, b2 bool) {
	var ib0, ib1, ib2 int

	if b0 {
		ib0 = 1
	}

	if b1 {
		ib1 = 1
	}

	if b2 {
		ib2 = 1
	}

	C.set_mouse_buttons(C.int(ib0), C.int(ib1), C.int(ib2))
}

func (i *IMGUISystem) SetMouseScrollPosition(xoffset, yoffset float64) {
	C.set_mouse_scroll_position(C.double(xoffset), C.double(yoffset))
}

func (i *IMGUISystem) EndFrame() {
	C.render()

	state := core.GetInputManager().State()
	state.SetMouseValid(false)
	state.SetKeysValid(false)
}

func (i *IMGUISystem) WantsCaptureMouse() bool {
	return int(C.wants_capture_mouse()) == 1
}

func (i *IMGUISystem) WantsCaptureKeyboard() bool {
	return int(C.wants_capture_keyboard()) == 1
}

type DrawData struct {
	drawData unsafe.Pointer
}

func (i *IMGUISystem) GetDrawData() core.IMGUIDrawData {
	return &DrawData{C.get_draw_data()}
}

func (d *DrawData) CommandListCount() int {
	return int(C.get_cmdlist_count(d.drawData))
}

func (d *DrawData) GetCommandList(index int) *core.IMGUICommandList {
	c_cmdList := C.get_cmdlist(d.drawData, C.int(index))

	cmdList := &core.IMGUICommandList{
		CmdBufferSize:    int(c_cmdList.commandBufferSize),
		VertexBufferSize: int(c_cmdList.vertexBufferSize),
		IndexBufferSize:  int(c_cmdList.indexBufferSize),
		VertexPointer:    unsafe.Pointer(c_cmdList.vertexPointer),
		IndexPointer:     unsafe.Pointer(c_cmdList.indexPointer),
		Commands:         make([]core.IMGUICommand, int(c_cmdList.commandBufferSize)),
	}

	for c := 0; c < cmdList.CmdBufferSize; c++ {
		cmd := core.IMGUICommand{}

		userTexturePtr := C.get_cmdlist_cmd(d.drawData, C.int(index), C.int(c),
			(*C.int)(unsafe.Pointer(&cmd.ElementCount)),
			(*C.float)(unsafe.Pointer(&cmd.ClipRect[0])))

		cmd.TextureID = userTexturePtr
		cmdList.Commands[c] = cmd
	}
	return cmdList
}
