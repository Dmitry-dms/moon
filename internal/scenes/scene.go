package scenes

import (
	"github.com/Dmitry-dms/moon/internal/components"
	"github.com/Dmitry-dms/moon/pkg/gogl"
)

// type Scene struct {
// 	renderer renderers.Renderer
// }

type Scene interface {
	Init()
	Start()
	Destroy()
	Update(dt float32)
	AddGameObjectToScene(obj *components.GameObject)
	GetCamera() *gogl.Camera
	Render()
	Load()

	Imgui()
}


