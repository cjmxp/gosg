#version 330 core

uniform sampler2D diffuseTex;

in vec2 Frag_UV;
in vec4 Frag_Color;

out vec4 Out_Color;

void main() {
    Out_Color = Frag_Color * texture(diffuseTex, Frag_UV.st);
}
