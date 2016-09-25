package core

// Framebuffer is an interface which wraps a render target. This contains
// information about depth and color attachments, dimensions.
type Framebuffer interface {
	// ColorTextureCount returns the number of color attachment textures.
	ColorTextureCount() uint8

	// DepthTexture returns the depth texture attachment.
	DepthTexture() Texture

	// ColorTexture returns the color texture attachment.
	ColorTexture(idx uint8) Texture
}
