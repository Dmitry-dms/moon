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
}

const (
	maxBatchSize = 1000
)

func NewGameWorld(name string, widthTiles, heightTiles int) *GameWorld {
	gw := GameWorld{
		Name:        name,
		widthTiles:  widthTiles,
		heightTiles: heightTiles,
		renderer:    renderers.NewRenderer(maxBatchSize),
		Id:          0,
	}
	return &gw
}

func (w *GameWorld) addGameObjToWorld(obj *components.GameObject) {
	w.renderer.AddGameObj(obj)
}

var g *components.GameObject
var sprsheet *gogl.Spritesheet

func (w *GameWorld) Init() {
	fmt.Printf("Init game world - %s \n", w.Name)

	sprsheet = gogl.AssetPool.GetSpriteSheet("assets/images/spritesheet.png")

	g = components.NewGameObject("Obj 1",
		components.NewTransform(mgl32.Vec2{100, 100}, mgl32.Vec2{256, 256}))
	spr := components.NewSpriteRenderer(mgl32.Vec4{1, 1, 1, 1}, sprsheet.GetSprite(0))
	g.AddSpriteRenderer(spr)

	g2 := components.NewGameObject("Obj 2",
		components.NewTransform(mgl32.Vec2{400, 100}, mgl32.Vec2{256, 256}))
	spr2 := components.NewSpriteRenderer(mgl32.Vec4{1, 1, 1, 1}, sprsheet.GetSprite(16))
	g2.AddSpriteRenderer(spr2)

	w.addGameObjToWorld(g)
	w.addGameObjToWorld(g2)

}
var spriteInd int
var spriteFlipTime, spriteFlipTimeLeft float32 = 0.2, 0
func (w *GameWorld) Update(dt float32) {
	spriteFlipTimeLeft-= dt
	if spriteFlipTimeLeft <= 0 {
		spriteFlipTimeLeft = spriteFlipTime
		spriteInd++
		if spriteInd > 4  {
			spriteInd = 0
		}
		g.SetSprite(sprsheet.GetSprite(spriteInd))
	}




	w.renderer.Update(dt)
}
func (w *GameWorld) Render(camera *gogl.Camera) {
	w.renderer.Render(camera)
}
