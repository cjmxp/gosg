#version 330 core

// global uniforms
struct light {
    mat4 vpMatrix;
    vec4 position;
    vec4 ambient;
    vec4 diffuse;
    vec4 specular;
};

layout (std140) uniform sceneBlock {
    mat4 vMatrix;
    mat4 pMatrix;
    mat4 vpMatrix;
    vec4 lightCount;
    light lights[16];
};

layout (location = 0) in vec3 position_in;

uniform mat4 mMatrix;

void main() {
    gl_Position = vpMatrix * mMatrix * vec4(position_in, 1.0);
}
