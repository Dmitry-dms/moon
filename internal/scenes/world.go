package scenes

import (
	"bufio"

	"fmt"
	"os"
	"regexp"

	"github.com/Dmitry-dms/moon/internal/components"
	"github.com/Dmitry-dms/moon/internal/listeners"

	// "golang.org/x/image/colornames"

	// "golang.org/x/image/colornames"

	// "github.com/Dmitry-dms/moon/internal/listeners"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/Dmitry-dms/moon/internal/renderers"
	"github.com/Dmitry-dms/moon/pkg/fonts"
	"github.com/Dmitry-dms/moon/pkg/gogl"

)

type GameWorld struct {
	widthTiles, heightTiles int // tile is 32x32 pixels
	Name                    string
	Id                      int
	renderer                renderers.SceneRenderer
	currGameObj             func(g *components.GameObject)
	gameObjects             []*components.GameObject
}

const (
	maxBatchSize = 1000
)

// var uir *ui.UiRenderer

func NewGameWorld(name string, widthTiles, heightTiles int, callback func(g *components.GameObject)) *GameWorld {
	gw := GameWorld{
		Name:        name,
		widthTiles:  widthTiles,
		heightTiles: heightTiles,
		renderer:    renderers.NewRenderer(maxBatchSize),
		Id:          0,
		currGameObj: callback,
		gameObjects: make([]*components.GameObject, 0),
	}
	// uir = ui.NewUIRenderer(maxBatchSize, 100)
	return &gw
}

func (w *GameWorld) AddGameObjToWorld(obj *components.GameObject) {
	w.gameObjects = append(w.gameObjects, obj)
	w.renderer.AddGameObj(obj)
}

var g *components.GameObject
var g2 *components.GameObject
var sprsheet *gogl.Spritesheet

func (w *GameWorld) loadResources() {
	gogl.AssetPool.GetTexture("assets/images/blend1.png")
	gogl.AssetPool.GetTexture("assets/images/blend2.png")
}




var batch *fonts.TextBatch

func (w *GameWorld) Init() {
	w.loadResources()
	fmt.Printf("Init game world - %s \n", w.Name)
	sprsheet = gogl.AssetPool.GetSpriteSheet("assets/images/decorations.png")
	// uir.Start()

	// font := fonts.NewFont("C:/Windows/Fonts/times.ttf", 20, true)
	// batch = fonts.NewTextBatch(font)
	// batch.Init()

	// for _, v := range w.gameObjects {
	// 	if v.Spr != nil {
	// 		if v.Spr.GetTexture() != nil {
	// 			v.Spr.SetTExture(gogl.AssetPool.GetTexture(v.Spr.GetTexture().GetFilepath()))
	// 		}
	// 	}
	// }

	// g = components.NewGameObject("Obj 1",
	// 	components.NewTransform(mgl32.Vec2{0, 0}, mgl32.Vec2{100, 100}), 1)
	// spr := components.DefSpriteRenderer()

	// sprite1 := gogl.DefSprite()
	// sprite1.SetTexture(gogl.AssetPool.GetTexture("assets/images/blend2.png"))
	// spr.SetSprite(sprite1)
	// g.AddSpriteRenderer(spr)
	// w.AddGameObjToWorld(g)


	// spr := ui.DefSpriteRenderer()

	// sprite1 := gogl.DefSprite()
	// sprite1.SetTexture(gogl.AssetPool.GetTexture("assets/images/blend1.png"))
	// spr.SetSprite(sprite1)
	// com := ui.Button{
	// 	UiObject: &ui.UiObject{
	// 		Transform: ui.NewTransform(mgl32.Vec2{300, 200}, mgl32.Vec2{100, 100}),
	// 		Name: "1",
	// 		Spr: spr,
	// 		ZIndex: 1,
	// 	},
	// }
	// uir.AddUIComponent(&com)

	// g2 = components.NewGameObject("Obj 2",
	// 	components.NewTransform(mgl32.Vec2{200, 100}, mgl32.Vec2{256, 256}), 2)
	// //spr2 := components.NewSpriteRenderer(mgl32.Vec4{1, 1, 1, 1}, gogl.NewSprite(gogl.AssetPool.GetTexture("assets/images/blend1.png")))
	// spr2 := &components.SpriteRenderer{}
	// spr2.SetColor(mgl32.Vec4{1, 1, 1, 1})

	// sprite2 := gogl.DefSprite()
	// sprite2.SetTexture(gogl.AssetPool.GetTexture("assets/images/blend1.png"))

	// spr2.SetSprite(sprite2)
	// g2.AddSpriteRenderer(spr2)


	// w.AddGameObjToWorld(g2)

	

}

func (w *GameWorld) Save() {
	file, err := os.Create(fmt.Sprintf("saves/%s.txt", w.Name))
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString("[\n")
	for i, v := range w.gameObjects {
		m, _ := v.MarshalJSON()
		s := string(m)
		writer.WriteString(s)
		if i != len(w.gameObjects)-1 {
			writer.WriteString(",\n")
		}
	}
	writer.WriteString("]")
	writer.Flush()
}

func (w *GameWorld) Load() {
	file, err := os.Open(fmt.Sprintf("saves/%s.txt", w.Name))
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	reg := regexp.MustCompile("\\{.*\\}")

	maxGoID := -1
	for reader.Scan() {
		t := reader.Text()
		if len([]rune(t)) != 1 {
			spitted := reg.FindAllString(t, 1)
			for _, spl := range spitted {
				obj := &components.GameObject{}

				err := obj.UnmarshalJSON([]byte(spl))
				if err != nil {
					fmt.Println(err)
				} else {
					w.AddGameObjToWorld(obj)

					if obj.GetUid() > maxGoID {
						maxGoID = obj.GetUid()
					}
				}
			}

		}
	}
	maxGoID++
	components.Init(maxGoID)
	//	w.currGameObj(w.gameObjects[0])

}

func (w *GameWorld) Update(dt float32, camera *gogl.Camera) {
	camera.UpdateProjection(mgl32.Vec2{float32(listeners.GetWindowWidth()),float32(listeners.GetWindowHeight())})


	// fmt.Println(float32(listeners.GetOrthoX()), float32(listeners.GetOrthoY()))
	// g.SetPosition(mgl32.Vec2{float32(listeners.GetOrthoX()), float32(listeners.GetOrthoY())})
	// g.SetPosition(mgl32.Vec2{float32(listeners.GetX()), float32(listeners.GetY())})
	// g.AddPosition(mgl32.Vec2{dt*10,0})
	// renderers.DebugDraw.AddBox2D(mgl32.Vec2{1000,500}, mgl32.Vec2{100,100}, 0, mgl32.Vec3{1,0,0}, 200000)
	// renderers.UpdateGridLines(camera) // Debug draw
	// w.renderer.Update(dt)
}
func (w *GameWorld) Render(camera *gogl.Camera) {

	// batch.AddText("My name is Dmitry", 0, 100, 2, colornames.Black)
	// batch.AddText("Привет, мир!\n920043 ~hghguij Progress #$@\n[A-Za-z] {0-9_20-33}", 450, 600, 1, colornames.Magenta)
	// batch.FlushBatch()
	// w.renderer.Render(camera)
}
