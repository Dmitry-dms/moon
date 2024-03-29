package scenes

import (
	"fmt"
	"time"

	// "unsafe"

	"github.com/Dmitry-dms/moon/internal/components"
	//	"github.com/Dmitry-dms/moon/internal/listeners"
	"github.com/Dmitry-dms/moon/pkg/gogl"

	//  "github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	//	mgl "github.com/go-gl/mathgl/mgl32"
	imgui "github.com/inkyblackness/imgui-go/v4"
)

type EditorScene struct {
	showDemoWindow bool
	glfw           *glfw.Window

	changeSceneCallback func(scene int)
	camera              *gogl.Camera

	activeGameWorld  *GameWorld
	activeGameObject *components.GameObject
	mouseControls    *components.MouseControls

	isRunning bool
}

func NewEditorScene(changeSceneCallback func(scene int)) *EditorScene {
	edtrScene := EditorScene{
		showDemoWindow:      true,
		changeSceneCallback: changeSceneCallback,
		camera:              gogl.NewCamera(mgl32.Vec2{0, 0}),
	}
	callback := func(g *components.GameObject) {
		edtrScene.activeGameObject = g
	}
	world := NewGameWorld("first", 64, 64, callback)
	edtrScene.activeGameWorld = world
	edtrScene.mouseControls = components.NewMouseControls(edtrScene.AddGameObjectToScene)

	return &edtrScene
}

func (e *EditorScene) GetCamera() *gogl.Camera {
	return e.camera
}

func (e *EditorScene) loadResources() {
	gogl.AssetPool.GetShader("assets/shaders/default.glsl")

	gogl.AssetPool.AddSpritesheet("assets/images/decorations.png",
		gogl.NewSpritesheet(gogl.AssetPool.GetTexture("assets/images/decorations.png"), 16, 16, 81, 0))
}

func (e *EditorScene) Load() {
	e.loadResources()
	e.activeGameWorld.loadResources()
	// e.activeGameWorld.Load()
}

func (e *EditorScene) Init() {

	fmt.Println("init editor scene")

	
	e.activeGameWorld.Init()

}


// при запуске сцены, запускаем объекты
func (e *EditorScene) Start() {
	// for _, g := range e.gameObjects {
	// 	if g != nil {
	// 		g.Start()
	// 		e.renderer.Add(g)
	// 	}
	// }
	// e.isRunning = true
}

var inc int

func (e *EditorScene) AddGameObjectToScene(obj *components.GameObject) {
	e.activeGameWorld.AddGameObjToWorld(obj)
}
func (e *EditorScene) Destroy() {
	e.activeGameWorld.Save()
}

func (e *EditorScene) Update(dt float32) {

	e.mouseControls.Update(dt)
	e.activeGameWorld.Update(dt, e.camera)
}

func (e *EditorScene) Render() {
	
	e.activeGameWorld.Render(e.camera)
}
func (e *EditorScene) Imgui() {
	{
		imgui.Begin("Test")
		imgui.Text(fmt.Sprintf("Application average %.3f ms/frame (%.1f FPS)",
			float32(time.Second.Milliseconds())/imgui.CurrentIO().Framerate(), imgui.CurrentIO().Framerate()))
		imgui.End()
	}
	imgui.Begin("Image picker")
	pos := imgui.WindowPos() //текущая позиция окна
	size := imgui.WindowSize()
	itemSpacing := imgui.CurrentStyle().ItemSpacing()

	windowX2 := pos.X + size.X
	for i := 0; i < sprsheet.Size(); i++ {
		sprite := sprsheet.GetSprite(i)
		spWidth := sprite.GetWidth() * 2
		spHeight := sprite.GetHeight() * 2
		id := sprite.GetTexture().GetId()
		texCoords := sprite.GetTextureCoords()

		imgui.PushIDInt(i)
		if imgui.ImageButtonV(imgui.TextureID(id), imgui.Vec2{X: float32(spWidth), Y: float32(spHeight)},
			imgui.Vec2{X: float32(texCoords[2].X()), Y: float32(texCoords[0].Y())},
			imgui.Vec2{X: float32(texCoords[0].X()), Y: float32(texCoords[2].Y())}, -1,
			imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0}, imgui.Vec4{X: 1, Y: 1, Z: 1, W: 1}) {

			obj := components.GenerateSpriteObject(sprite, 32, 32)
			//привязываем к курсору
			e.mouseControls.PickupObject(obj)

		}

		imgui.PopID()

		lastBtnPos := imgui.ItemRectMax()
		lastX2 := lastBtnPos.X
		nextBtnX2 := lastX2 + itemSpacing.X + float32(spWidth)
		if (i+1) < sprsheet.Size() && nextBtnX2 < windowX2 {
			//разместим следующую кнопку в этой строке
			imgui.SameLine()
		}
	}
	imgui.End()

	{
		if e.activeGameObject != nil {
			imgui.Begin("Inspector")
			e.activeGameObject.Imgui()
			imgui.End()
		}
	}

	if e.showDemoWindow {
		// Normally user code doesn't need/want to call this because positions are saved in .ini file anyway.
		// Here we just want to make the demo initial state a bit more friendly!
		const demoX = 650
		const demoY = 20
		imgui.SetNextWindowPosV(imgui.Vec2{X: demoX, Y: demoY}, imgui.ConditionFirstUseEver, imgui.Vec2{})
		imgui.ShowDemoWindow(&e.showDemoWindow)
	}
}
