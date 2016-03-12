#version 330 core

// global uniforms
struct light {
    mat4 vpMatrix;
    vec4 position;
    vec4 ambient;
    vec4 diffuse;
    vec4 specular;
};

layout (std140) uniform cameraConstants {
    mat4 vMatrix;
    mat4 pMatrix;
    mat4 vpMatrix;
    vec4 lightCount;
    light lights[16];
};

in vec3 position;
in vec3 cameraPosition;
in vec3 tcoords0;
in mat3 tbn;

layout (location = 0) out vec4 color;

uniform sampler2D diffuseTex;
uniform sampler2D normalTex;
uniform sampler2D shadowTex;

float shadowMap(vec2 coords, float compare) {
    return step(compare, texture(shadowTex, coords).r);
}

float shadowMapLinear(vec2 coords, float compare) {
    vec2 texelSize = vec2(1.0/2048.0);
    
    vec2 pixelPos = coords/texelSize + vec2(0.5);
    vec2 fracPart = fract(pixelPos);
    vec2 startTexel = (pixelPos - fracPart) * texelSize;
    
    float blTexel = shadowMap(startTexel, compare);
    float brTexel = shadowMap(startTexel + vec2(texelSize.x, 0.0), compare);
    float tlTexel = shadowMap(startTexel + vec2(0.0, texelSize.y), compare);
    float trTexel = shadowMap(startTexel + texelSize, compare);
    
    float mixA = mix(blTexel, tlTexel, fracPart.y);
    float mixB = mix(brTexel, trTexel, fracPart.y);
    
    return mix(mixA, mixB, fracPart.x);
}

float shadowMapPCF(vec2 coords, float compare) {
    const float NUM_SAMPLES = 3.0f;
    const float SAMPLES_START = (NUM_SAMPLES-1.0)/2.0;
    const float NUM_SAMPLES_SQUARED = NUM_SAMPLES*NUM_SAMPLES;
    const vec2 texelSize = vec2(1.0/2048.0);
    
    float result = 0.0;
    
    for (float y=-SAMPLES_START; y<=SAMPLES_START; y+=1.0) {
        for (float x=-SAMPLES_START; x<=SAMPLES_START; x+=1.0) {
            vec2 coordsOffset = vec2(x, y)*texelSize;
            result += shadowMapLinear(coords + coordsOffset, compare);
        }
    }
    
    return result/NUM_SAMPLES_SQUARED;
}

float shadow(vec4 coords, float dotNL) {
    float cosTheta = clamp(dotNL, 0.0, 1.0);
    float bias = clamp(0.005*tan(acos(cosTheta)), 0.0, 0.01);

    vec3 shadowMapCoords = coords.xyz/coords.w;
    float compare = shadowMapCoords.z - 1.0/512.0;
    return shadowMapPCF(shadowMapCoords.xy, compare);
}

void main() {
    vec3 N = normalize(tbn * (texture(normalTex, tcoords0.st).rgb * 2.0 - 1.0));
    vec3 E = normalize(cameraPosition - position);

    vec4 material = texture(diffuseTex, tcoords0.st);
    vec3 materialDiffuse = material.rgb;
    vec3 materialSpecular = materialDiffuse;
    float materialShininess = 255.0;

    int lc = int(lightCount[0]);
    color.rgb = vec3(0.0);
    
    for (int i=0; i<lc; i++) {
        vec3 diffuseTerm = materialDiffuse * lights[i].diffuse.rgb;
        vec3 specularTerm = materialSpecular * lights[i].specular.rgb;

        vec3 lightPosition = lights[i].position.xyz;
        vec3 L = normalize(lightPosition - position);
        vec3 R = normalize(-reflect(L, N));

        vec3 colorDiffuse = diffuseTerm * max(dot(N, L), 0.0);
        vec3 colorSpecular = pow(max(dot(R, E), 0.0), materialShininess) * specularTerm;

        // shadow
        vec4 lightPos = lights[i].vpMatrix * vec4(position, 1.0);
        color.rgb += (colorDiffuse + colorSpecular) * shadow(lightPos, dot(N, L));
    }
    
    color.rgb = 1.0f - exp2(-color.rgb * 2.2);
    color.a = material.a;
}
