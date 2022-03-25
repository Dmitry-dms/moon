package scenes

import (
	"bufio"

	"fmt"
	"os"
	"regexp"

	"github.com/Dmitry-dms/moon/internal/components"

	"github.com/Dmitry-dms/moon/internal/renderers"
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

func (w *GameWorld) Init() {
	w.loadResources()
	fmt.Printf("Init game world - %s \n", w.Name)

	sprsheet = gogl.AssetPool.GetSpriteSheet("assets/images/decorations.png")

	// g = components.NewGameObject("Obj 1",
	// 	components.NewTransform(mgl32.Vec2{400, 100}, mgl32.Vec2{256, 256}), 1)
	// spr := components.DefSpriteRenderer()

	// sprite1 := gogl.DefSprite()
	// spr.SetSprite(sprite1)
	// g.AddSpriteRenderer(spr)

	// g2 = components.NewGameObject("Obj 2",
	// 	components.NewTransform(mgl32.Vec2{200, 100}, mgl32.Vec2{256, 256}), 2)
	// //spr2 := components.NewSpriteRenderer(mgl32.Vec4{1, 1, 1, 1}, gogl.NewSprite(gogl.AssetPool.GetTexture("assets/images/blend1.png")))
	// spr2 := &components.SpriteRenderer{}
	// spr2.SetColor(mgl32.Vec4{1, 1, 1, 1})

	// sprite2 := gogl.DefSprite()
	// sprite2.SetTexture(gogl.AssetPool.GetTexture("assets/images/blend1.png"))

	// spr2.SetSprite(sprite2)
	// g2.AddSpriteRenderer(spr2)

	// w.addGameObjToWorld(g)
	// w.addGameObjToWorld(g2)

	// s, err := g.MarshalJSON()
	// fmt.Println(string(s), err)

}

type exportedWorld struct {
	GameObjects []*components.GameObject `json:"game_objects"`
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
		if i != len(w.gameObjects) -1 {
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
				//err := json.Unmarshal([]byte(spl), obj)
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

var f bool

func (w *GameWorld) Update(dt float32) {

	w.renderer.Update(dt)
}
func (w *GameWorld) Render(camera *gogl.Camera) {
	if f == false {
		// fmt.Printf("BEFORE = %v \n", g2)
		// m, err := g2.MarshalJSON()
		// fmt.Println(string(m), err)

		// var newObj components.GameObject
		// newObj.UnmarshalJSON(m)
		// fmt.Printf("AFTER = %v \n", newObj)

		f = true
	}
	w.renderer.Render(camera)
}
