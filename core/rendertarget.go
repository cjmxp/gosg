package core

// RenderTarget is an interface which wraps a render target. This contains
// information about framebuffers and their attachments.
type RenderTarget interface {
	// ColorTextureCount returns the number of color attachment textures.
	ColorTextureCount() uint8

	// DepthTexture returns the depth texture attachment.
	DepthTexture() Texture

	// ColorTexture returns the color texture attachment.
	ColorTexture(idx uint8) Texture
}
