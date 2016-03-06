#version 330 core

in vec2 tcoords0;

layout (location = 0) out vec3 color;

uniform sampler2D zTex;
uniform vec2 outputsize;
uniform vec2 ssaoNoise[16];

float readDepth( in vec2 coord ) {
    vec2 camerarange = vec2(0.1, 100.0);
    return (2.0 * camerarange.x) / (camerarange.y + camerarange.x - texture(zTex, coord ).x * (camerarange.y - camerarange.x));
}

float ssao(float depth) {
    float d = 0.0, aoCap = 1.0, ao = 0.0;

    float aoMultiplier = 10000.0;
    float depthTolerance = 1.0/10000.0;

    for (int j=0; j<16; j++) {
        d = readDepth(tcoords0 + ssaoNoise[j]/outputsize);
        ao += min(aoCap, max(0.0, depth-d-depthTolerance) * aoMultiplier);
    }
    return 1.0 - ao/16.0;
}

void main() {
    float depth = readDepth(tcoords0);
    color = vec3(ssao(depth));
    
    color = vec3(texture(zTex, tcoords0).r);
}
