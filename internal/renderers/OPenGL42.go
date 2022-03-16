package renderers

import (
	_ "embed" // using embed for the shader sources


//	"github.com/go-gl/gl/v4.2-core/gl"
	//"github.com/pkg/errors"
)

// //go:embed gl-shader/main.vert
// var unversionedVertexShader42 string

// //go:embed gl-shader/main.frag
// var unversionedFragmentShader42 string

// OpenGL3 implements a renderer based on github.com/go-gl/gl (v4.2-core).
type OpenGL42 struct {
	//imguiIO imgui.IO

}

// NewOpenGL3 attempts to initialize a renderer.
// An OpenGL context has to be established before calling this function.
func NewOpenGL42() (*OpenGL42, error) {
	//err := gl.Init()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to initialize OpenGL 4.2-core")
	// }

	renderer := &OpenGL42{
		//glslVersion: "#version 420",
	}
	//renderer.createDeviceObjects()

	//imguiIO.SetBackendFlags(imguiIO.GetBackendFlags() | imgui.BackendFlagsRendererHasVtxOffset)

	return renderer, nil
}





// Render translates the ImGui draw data to OpenGL42 commands.
func (renderer *OpenGL42) Render() {
	
}

func (renderer *OpenGL42) Dispose() {

}

func (renderer *OpenGL42) PreRender(clearColor [3]float32) {

}


