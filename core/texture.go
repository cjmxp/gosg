package core

import "unsafe"

type TextureTarget int

const (
	TextureTarget1D TextureTarget = iota
	TextureTarget1DArray
	TextureTarget2D
	TextureTarget2DArray
	TextureTargetCubemapXPositive
	TextureTargetCubemapXNegative
	TextureTargetCubemapYPositive
	TextureTargetCubemapYNegative
	TextureTargetCubemapZPositive
	TextureTargetCubemapZNegative
)

type TextureFormat int

const (
	TextureFormatR TextureFormat = iota
	TextureFormatRG
	TextureFormatRGB
	TextureFormatRGBA
	TextureFormatDEPTH
)

type TextureSizedFormat int

const (
	TextureSizedFormatR8 TextureSizedFormat = iota
	TextureSizedFormatR16F
	TextureSizedFormatR32F
	TextureSizedFormatRG8
	TextureSizedFormatRG16F
	TextureSizedFormatRG32F
	TextureSizedFormatRGB8
	TextureSizedFormatRGB16F
	TextureSizedFormatRGB32F
	TextureSizedFormatRGBA8
	TextureSizedFormatRGBA16F
	TextureSizedFormatRGBA32F
	TextureSizedFormatDEPTH32F
)

type TextureComponentType int

const (
	TextureComponentTypeUNSIGNEDBYTE TextureComponentType = iota
	TextureComponentTypeFLOAT
)

type TextureWrapMode int

const (
	TextureWrapModeClampEdge TextureWrapMode = iota
	TextureWrapModeClampBorder
)

type TextureFilter int

const (
	TextureFilterNearest TextureFilter = iota
	TextureFilterLinear
	TextureFilterMipmapLinear
)

// TextureDescriptor contains the full description of a texture and its sampling parameters
// It is used as input to texture creation functions and at runtime inside rendersystems
// to setup samplers and memory allocation
type TextureDescriptor struct {
	Width         uint32
	Height        uint32
	Mipmaps       bool
	Target        TextureTarget
	Format        TextureFormat
	SizedFormat   TextureSizedFormat
	ComponentType TextureComponentType
	Filter        TextureFilter
	WrapMode      TextureWrapMode
}

// Texture is an interface which wraps both a texture and settings for samplers sampling it
type Texture interface {
	Descriptor() TextureDescriptor

	Handle() unsafe.Pointer

	// Lt is used for sorting
	Lt(Texture) bool

	// Gt
	Gt(Texture) bool
}
