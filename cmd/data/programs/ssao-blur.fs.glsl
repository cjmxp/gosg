#version 330 core

in vec2 tcoords0;

layout (location = 0) out vec3 color;

uniform sampler2D ssaoTex;

float blurSSAO(int blursize) {
   vec2 texelSize = 1.0 / vec2(textureSize(ssaoTex, 0));
   float result = 0.0;
   vec2 hlim = vec2(float(-blursize) * 0.5 + 0.5);
   for (int i = 0; i < blursize; ++i) {
      for (int j = 0; j < blursize; ++j) {
         vec2 offset = (hlim + vec2(float(i), float(j))) * texelSize;
         result += texture(ssaoTex, tcoords0 + offset).b;
      }
   }

   return result / float(blursize * blursize);
}

void main() {
    color = vec3(blurSSAO(4));
}
