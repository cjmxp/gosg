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

const float A = 0.15;
const float B = 0.50;
const float C = 0.10;
const float D = 0.20;
const float E = 0.02;
const float F = 0.30;
const float W = 11.2;

const float PI = 3.14159265358979323846;
const float INV_PI = 1.0/PI;

vec3 Uncharted2Tonemap(vec3 x) {
   return ((x*(A*x+C*B)+D*E)/(x*(A*x+B)+D*F))-E/F;
}

vec3 tonemapUncharted2(vec3 color) {
    float ExposureBias = 2.0;
    vec3 curr = Uncharted2Tonemap(ExposureBias * color);
    vec3 whiteScale = 1.0 / Uncharted2Tonemap(vec3(W));
    return curr * whiteScale;
}

// beckmann
float distribution(vec3 n, vec3 h, float roughness) {
    float m_Sq = roughness * roughness;
    float NdotH_Sq = max(dot(n, h), 0.0);
    NdotH_Sq = NdotH_Sq * NdotH_Sq;
    return exp((NdotH_Sq - 1.0) / (m_Sq*NdotH_Sq)) / (3.14159265 * m_Sq * NdotH_Sq * NdotH_Sq);
}

// cook-torrance
float geometry(vec3 n, vec3 h, vec3 v, vec3 l, float roughness) {
    float NdotH = dot(n, h);
    float NdotL = dot(n, l);
    float NdotV = dot(n, v);
    float VdotH = dot(v, h);
    float NdotL_clamped = max(NdotL, 0.0);
    float NdotV_clamped = max(NdotV, 0.0);
    return min(min(2.0 * NdotH * NdotV_clamped / VdotH, 2.0 * NdotH * NdotL_clamped / VdotH), 1.0);
}

// schlich
float fresnel(float f0, vec3 n, vec3 l) {
    return f0 + (1.0 - f0) * pow(1.0 - dot(n, l), 5.0);
}

// fresnel diff
float diffuse_energy_ratio(float f0, vec3 n, vec3 l) {
    return 1.0 - fresnel(f0, n, l);
}

void main() {
    // init
    color.rgb = vec3(0.0);

    // init materials
    vec4 albedo = texture(diffuseTex, tcoords0.st);

    // will come from textures later
    float metalness = 0.0;
    float roughness = 0.3;

    // adjust f0 from 0.118 to 0.818, this will normally be discrete
    float f0 = 0.118 + metalness * 0.7; //max 0.818

    // normal, eye/view
    vec3 N = normalize(tbn * (texture(normalTex, tcoords0.st).rgb * 2.0 - 1.0));
    vec3 V = normalize(cameraPosition - position);

    // shared products
    float NdotV = dot(N, V);
    float NdotV_clamped = max(NdotV, 0.0000000001);

    int lc = int(lightCount[0]);
    for (int i=0; i<lc; i++) {
        // lightdir, halfvec
        vec3 L = normalize(lights[i].position.xyz - position);
        vec3 H = normalize(L + V);

        float NdotL = dot(N, L);
        float NdotL_clamped = max(NdotL, 0.0);

        float fres = fresnel(f0, H, L);
        float geom = geometry(N, H, V, L, roughness);
        float ndf = distribution(N, H,  roughness);

        float brdf_spec = (0.25 * fres * geom * ndf) / (NdotL_clamped * NdotV_clamped);

        vec3 color_spec = NdotL_clamped * brdf_spec * lights[i].diffuse.rgb;
        vec3 color_diff = NdotL_clamped * diffuse_energy_ratio(f0, N, L) * albedo.rgb * lights[i].diffuse.rgb;
        color.rgb += (color_diff + color_spec) * shadow(lights[i].vpMatrix * vec4(position, 1.0), dot(N, L));
    }

    color.a = albedo.a;
    color.rgb = tonemapUncharted2(color.rgb);
    color.rgb = pow(color.rgb, vec3(1.0/2.2));
}