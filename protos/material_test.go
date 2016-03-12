package protos

import (
	"testing"

	"github.com/golang/protobuf/jsonpb"
)

func TestMaterialUnmarshal(t *testing.T) {
	clearMaterialJson := `{
	"programName": "",
	"culling": true,
	"cullFace": "CULL_BACK",
	"blending": false,
	"blendSrcMode": "BLEND_SRC_ALPHA",
	"blendDstMode": "BLEND_ONE_MINUS_SRC_ALPHA",
	"blendEquation": "BLEND_FUNC_ADD",
	"depthTest": true,
	"depthWrite": true,
	"depthFunc": "DEPTH_LESS_EQUAL",
	"colorWrite": true,
	"scissorTest": true
	}`

	var clearMaterial Material
	if err := jsonpb.UnmarshalString(clearMaterialJson, &clearMaterial); err != nil {
		t.Error("Cannot unmarshal material from json: ", err)
	}

	if clearMaterial.ProgramName != "" {
		t.Error("Wrong program name")
	}

	if clearMaterial.Culling != true {
		t.Error("Wrong culling")
	}

	if clearMaterial.Blending != false {
		t.Error("Wrong blending")
	}

	if clearMaterial.BlendSrcMode != Material_BLEND_SRC_ALPHA {
		t.Error("Wrong blendSrcMode")
	}

	if clearMaterial.BlendDstMode != Material_BLEND_ONE_MINUS_SRC_ALPHA {
		t.Error("Wrong blendDstMode")
	}

	if clearMaterial.BlendEquation != Material_BLEND_FUNC_ADD {
		t.Error("Wrong blendFunc")
	}

	if clearMaterial.DepthTest != true {
		t.Error("Wrong depthTest")
	}

	if clearMaterial.DepthWrite != true {
		t.Error("Wrong depthWrite")
	}

	if clearMaterial.DepthFunc != Material_DEPTH_LESS_EQUAL {
		t.Error("Wrong depthFunc")
	}

	if clearMaterial.ColorWrite != true {
		t.Error("Wrong colorWrite")
	}

	if clearMaterial.ScissorTest != true {
		t.Error("Wring scissorWrite")
	}

	t.Log(clearMaterial)
}
