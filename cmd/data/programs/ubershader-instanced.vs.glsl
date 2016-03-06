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

// this is the same for all our models
layout (location = 0) in vec3 position_in;
layout (location = 1) in vec3 normal_in;
layout (location = 2) in vec3 tangent_in;
layout (location = 3) in vec3 bitangent_in;
layout (location = 4) in vec3 tcoords0_in;
layout (location = 5) in mat4 mMatrix;

out vec3 position;
out vec3 cameraPosition;
out vec3 tcoords0;
out mat3 tbn;

void main() {
    // clip position
    gl_Position = vpMatrix * mMatrix * vec4(position_in, 1.0);

    // world position & camera world position
    position = (mMatrix * vec4(position_in, 1.0)).xyz;
    cameraPosition = inverse(vMatrix)[3].rgb;

    // world TBN
    vec3 normal = normalize((mMatrix * vec4(normal_in, 0.0)).xyz);
    vec3 tangent = normalize((mMatrix * vec4(tangent_in, 0.0)).xyz);
    vec3 bitangent = normalize((mMatrix * vec4(bitangent_in, 0.0)).xyz);

    tbn = mat3(tangent, bitangent, normal);

    // texture coordinates
    tcoords0 = tcoords0_in;
}
