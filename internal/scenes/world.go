package scenes

import (
	"fmt"

	"github.com/Dmitry-dms/moon/internal/components"

	"github.com/Dmitry-dms/moon/internal/renderers"
	"github.com/Dmitry-dms/moon/pkg/gogl"

	"github.com/go-gl/mathgl/mgl32"
)

type GameWorld struct {
	widthTiles, heightTiles int // tile is 32x32 pixels
	Name                    string
	Id                      int
	renderer                renderers.GameRenderer
	currGameObj             func(g *components.GameObject)
}

const (
	maxBatchSize = 1000
)

func NewGameWorld(name string, widthTiles, heightTiles int, callback func(g *components.GameObject)) *GameWorld {
	gw := GameWorld{
		Name:        name,
		widthTiles:  widthTiles,
		heightTiles: heightTiles,
		renderer:    renderers.NewRenderer(maxBatchSize),
		Id:          0,
		currGameObj: callback,
	}
	return &gw
}

func (w *GameWorld) addGameObjToWorld(obj *components.GameObject) {
	w.renderer.AddGameObj(obj)
}

var g *components.GameObject
var sprsheet *gogl.Spritesheet

func (w *GameWorld) loadResources() {
	gogl.AssetPool.GetTexture("assets/images/blend1.png")
	gogl.AssetPool.GetTexture("assets/images/blend2.png")
}

func (w *GameWorld) Init() {
	w.loadResources()
	fmt.Printf("Init game world - %s \n", w.Name)

	sprsheet = gogl.AssetPool.GetSpriteSheet("assets/images/spritesheet.png")

	g = components.NewGameObject("Obj 1",
		components.NewTransform(mgl32.Vec2{400, 100}, mgl32.Vec2{256, 256}), 1)
	spr := components.NewSpriteRenderer(mgl32.Vec4{1, 1, 1, 1}, gogl.NewSprite(nil))
	g.AddSpriteRenderer(spr)

	w.currGameObj(g)

	g2 := components.NewGameObject("Obj 2",
		components.NewTransform(mgl32.Vec2{200, 100}, mgl32.Vec2{256, 256}), 2)
	spr2 := components.NewSpriteRenderer(mgl32.Vec4{1, 1, 1, 1}, gogl.NewSprite(gogl.AssetPool.GetTexture("assets/images/blend1.png")))
	g2.AddSpriteRenderer(spr2)

	w.addGameObjToWorld(g)
	w.addGameObjToWorld(g2)

}

func (w *GameWorld) Update(dt float32) {

	w.renderer.Update(dt)
}
func (w *GameWorld) Render(camera *gogl.Camera) {
	w.renderer.Render(camera)
}
