#version 330 core

layout (location = 0) out vec4 color;

uniform vec4 in_color;

void main() {
    color = in_color;
}
