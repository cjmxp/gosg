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

/* BRDF */
float fresnelSchlick( float f0, float dot_v_h ) {
    return f0 + ( 1 - f0 ) * pow( 1.0 - dot_v_h, 5.0 );
}

float fresnelSchlick( float dot_v_h )
{
    float value  = clamp( 1 - dot_v_h, 0.0, 1.0 );
    float value2 = value * value;
    return ( value2 * value2 * value );
}

float geomSchlickGGX( float alpha, float dot_n_v )
{
    float k    = 0.5 * alpha;
    float geom = dot_n_v / ( dot_n_v * ( 1 - k ) + k );

    return geom;
}

float ndfGGX(float alpha, float NdotH )
{
    float alpha2 = alpha * alpha;
    float t      = 1 + ( alpha2 - 1 ) * NdotH * NdotH;
    return INV_PI * alpha2 / ( t * t );
}

vec3 BRDF(vec3 l, vec3 v, vec3 n, vec3 diffuse_color, float metalness, float roughness) {
    float dot_n_v = dot(n, v);
    float dot_n_l = dot(n, l);
    float dot_n_l_clamp  = max(dot_n_l, 0.0);

    vec3 h = normalize(l+v);
    float dot_v_h = dot(v, h);
    float fresnel = fresnelSchlick(dot_v_h);
    vec3 f0_color = mix( diffuse_color, vec3(1.0), (1.0 - metalness));
    vec3 specular_color = mix(f0_color, vec3(1.0), fresnel);

    float clamp_roughness = max(roughness, 0.01);
    float alpha = clamp_roughness * clamp_roughness;

    float dot_n_h = dot(n, h);
    float ndf = ndfGGX(alpha, dot_n_h);
    float geom = geomSchlickGGX(alpha, dot_n_v);
    float specular_brdf = (0.25 * ndf * geom ) / (dot_n_l * dot_n_v);
    specular_color *= specular_brdf;

    return dot_n_l_clamp * (INV_PI * (1.0 - metalness) * diffuse_color + specular_color);
}

void main() {
    vec3 N = normalize(tbn * (texture(normalTex, tcoords0.st).rgb * 2.0 - 1.0));
    vec3 E = normalize(cameraPosition - position);
    vec4 material = texture(diffuseTex, tcoords0.st);

    int lc = int(lightCount[0]);
    color.rgb = vec3(0.0);

    for (int i=0; i<lc; i++) {
        vec3 lightPosition = lights[i].position.xyz;
        vec3 L = normalize(lightPosition - position);
        vec4 shadowLightPos = lights[i].vpMatrix * vec4(position, 1.0);
        color.rgb += BRDF(L, E, N, material.rgb, 0.0, 0.7) * shadow(shadowLightPos, dot(N, L));
    }
    color.a = material.a;
    color.rgb = tonemapUncharted2(color.rgb);
    color.rgb = pow(color.rgb, vec3(1.0/2.2));
}