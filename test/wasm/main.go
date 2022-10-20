//go:build js && wasm

package main

import (
	"github.com/Dmitry-dms/moon/pkg/ui"
	"github.com/nuberu/webgl"
	"syscall/js"
)

var uiCtx *ui.UiContext
var vert = `layout (location=0) in vec3 aPos;
layout (location=1) in vec4 aColor;
layout (location=2) in vec2 aTexCoords;
layout (location=3) in float aTexId;

uniform mat4 uProjection;

out vec4 fColor;
out vec2 fTexCoords;
out float fTexId;

void main()
{
    fColor = aColor;
    fTexCoords = aTexCoords;
    fTexId = aTexId;
    gl_Position = uProjection * vec4(aPos,1.0);
}
`

var frag = `
in vec4 fColor;
in vec2 fTexCoords;
in float fTexId;
out vec4 color;

uniform sampler2D Texture;

void main()
{
    if (fTexId > 0) {
        vec4 tC = texture(Texture,fTexCoords);
        color =  fColor * tC;
    } else {
        color = fColor;
    }
}`

func main() {
	doc := js.Global().Get("document")
	canvasEl := doc.Call("createElement", "canvas")
	doc.Get("body").Call("appendChild", canvasEl)
	width := 800
	height := 600
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	//
	gl, err := webgl.FromCanvas(canvasEl)
	if err != nil {
		panic(err)
	}

	front := newGl(gl)
	//
	uiCtx.Initialize(front)

	//renderFrame := js.NewCallback(func(args []js.Value) {
	//	// Calculate rotation rate
	//	now := float32(args[0].Float())
	//	tdiff := now - tmark
	//	tmark = now
	//	rotation = rotation + float32(tdiff)/500
	//
	//	// Do new model matrix calculations
	//	movMatrix = mgl32.HomogRotate3DX(0.5 * rotation)
	//	movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DY(0.3 * rotation))
	//	movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DZ(0.2 * rotation))
	//
	//	// Convert model matrix to a JS TypedArray
	//	var modelMatrixBuffer *[16]float32
	//	modelMatrixBuffer = (*[16]float32)(unsafe.Pointer(&movMatrix))
	//
	//	// Apply the model matrix
	//	gl.UniformMatrix4fv(ModelMatrix, false, []float32((*modelMatrixBuffer)[:]))
	//
	//	// Clear the screen
	//	gl.Enable(webgl.DEPTH_TEST)
	//	gl.Clear(uint32(webgl.COLOR_BUFFER_BIT) & uint32(webgl.DEPTH_BUFFER_BIT))
	//
	//	// Draw the cube
	//	gl.DrawElements(webgl.TRIANGLES, len(indices), webgl.UNSIGNED_SHORT, 0)
	//
	//	// Call next frame
	//	js.Global().Call("requestAnimationFrame", renderFrame)
	//})
	//defer renderFrame.Release()

	//js.Global().Call("requestAnimationFrame", renderFrame)

	done := make(chan struct{}, 0)
	<-done
	//gl.Enable(gl.BLEND)
	//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	//
	//tex, _ = tex.Init("assets/images/mario.png")
	//tex2, _ = tex2.Init("assets/images/goomba.png")
	//
	//for ; ; {
	//
	//	gl.ClearColor(1, 1, 1, 1)
	//	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	//
	//	uiCtx.NewFrame([2]float32{float32(Width), float32(Height)})
	//
	//	firstWindow()
	//
	//	uiCtx.EndFrame([2]float32{float32(Width), float32(Height)})
	//
	//	window.SwapBuffers()
	//}
}

func firstWindow() {
	uiCtx.BeginWindow("first wnd")

	if uiCtx.ButtonT("Нажать", "Press") {
		//	ish = !ish
		//
	}
	//if ish {
	//	uiCtx.Text("#er", "Wdff213 ello world!", 14)
	//	uiCtx.Text("#fgfgd", "hello world!", 14)
	//}

	uiCtx.TreeNode("tree1", "Configuration", func() {
		uiCtx.Text("text-ttp-1", "Обычная картинка, которая  ничего не делает", 14)
		uiCtx.Text("#t3j", "hello world!", 14)
		uiCtx.TreeNode("tree1yuy2", "Настройки", func() {
			uiCtx.Text("texiyt-ttp-1", "Обычная картинка, которая  ничего не делает", 14)
			uiCtx.Text("#tiy3j", "hello world!", 14)
		})
	})

	//uiCtx.VSpace("#vs1fdgdf")

	//uiCtx.Row("row 13214", func() {
	//	uiCtx.Image("#im4kjdg464", tex)
	//	uiCtx.Column("col fdfd", func() {
	//		uiCtx.Image("#im76", tex2)
	//		uiCtx.Image("#im4", tex)
	//	})
	//
	//	uiCtx.Column("col fdfdвава", func() {
	//		uiCtx.Button("ASsfdffb")
	//		uiCtx.Button("ASsfdffbbb")
	//		uiCtx.Slider("slider-1", &slCounter, 0, 255)
	//	})
	//
	//	uiCtx.Image("#im4kj", tex)
	//})
	//uiCtx.SubWidgetSpace("widhsp-1", ui.Default, func() {
	//	uiCtx.Image("#im4kjdg464tht", tex2)
	//	uiCtx.Image("#im76erewr", tex)
	//	uiCtx.Text("#t3df", "world!", 24)
	//})
	////uiCtx.VSpace("#vs1")
	//uiCtx.Image("#imgj4", tex2)
	//
	//if uiCtx.ActiveWidget == "#imgj4" {
	//	uiCtx.Tooltip("ttp-1", func() {
	//		uiCtx.Text("text-ttp-1", "Обычная картинка, которая  ничего не делает", 14)
	//		uiCtx.Text("text-ttp-2", "Hello World", 16)
	//		uiCtx.Text("text-ttp-3", "Hello World", 16)
	//	})
	//}
	//if uiCtx.ActiveWidget == "widgId" {
	//
	//}

	uiCtx.EndWindow()
}
