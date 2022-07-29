#type vertex
#version 420
uniform mat4 ProjMtx;

in vec2 Position;
in vec2 UV;
in vec4 Color;

out vec2 Frag_UV;
out vec4 Frag_Color;

void main()
{
    Frag_UV = UV;
    Frag_Color = Color;
    gl_Position = ProjMtx * vec4(Position.xy, 0, 1);
}

#type fragment
#version 420
uniform sampler2D Texture;


in vec2 Frag_UV;
in vec4 Frag_Color;


out vec4 color;

void main()
{
    color = vec4(Frag_Color.rgb, Frag_Color.a * texture(Texture, Frag_UV.st).r);
    
    // if (fTexId > 0) {
    //    int id = int(fTexId);
    //     color = Frag_Color * texture(uTextures[id],fTexCoords);
    //     //color = vec4(fTexCoords, 0, 1);//чтобы узнать u v координаты
    // } else {
    //     color = Frag_Color;
    // }
}