package scenes

import (
	"bufio"

	"fmt"
	"os"
	"regexp"

	"github.com/Dmitry-dms/moon/internal/components"
	"github.com/go-gl/mathgl/mgl32"

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
var angle float32 = 30
func (w *GameWorld) Update(dt float32, camera *gogl.Camera) {
	renderers.DebugDraw.AddCircle(mgl32.Vec2{200,200}, 64,  mgl32.Vec3{1,0,0}, 1)
	angle += 40*dt

	renderers.UpdateGridLines(camera)
	w.renderer.Update(dt)
}
func (w *GameWorld) Render(camera *gogl.Camera) {

	w.renderer.Render(camera)
}
