#type vertex
#version 420

layout (location=0) in vec3 aPos;
layout (location=1) in vec3 aColor;


uniform mat4 uProjection;
uniform mat4 uView;

out vec3 fColor;


void main()
{
    fColor = aColor;
    gl_Position = uProjection * uView * vec4(aPos, 1.0);
}

#type fragment
#version 420

in vec3 fColor;

out vec4 color;

// uniform sampler2D uTextures[8];
void main()
{
        color =vec4(fColor,1.0);
}