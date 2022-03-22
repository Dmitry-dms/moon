package renderers

import (
	"sort"
	"sync"

	"github.com/Dmitry-dms/moon/internal/components"
	"github.com/Dmitry-dms/moon/pkg/gogl"
)

// // Renderer covers rendering imgui draw data.
// type Renderer interface {
// 	// PreRender causes the display buffer to be prepared for new output.
// 	PreRender(clearColor [3]float32)
// 	// Render draws the provided imgui draw data.
// 	Render()

// 	Dispose()
// }

// интерфейс, который отвечает за отрисовку всех игровых объектов
type GameRenderer interface {
	Update(dt float32)
	Render(camera *gogl.Camera)
	AddGameObj(obj *components.GameObject)
}

type Renderer struct {
	maxBatchSize     int
	gameObjetcsCount uint
	batches          []*RenderBatch
	//GameObjects      []GameObject
	wg *sync.WaitGroup
}

func NewRenderer(maxBatchSize int) *Renderer {
	//gObjs := make([]GameObject, 0, 1000)
	r := Renderer{
		maxBatchSize: maxBatchSize,
		batches:      make([]*RenderBatch, 0, 1000),
		//GameObjects:  gObjs,
		wg: &sync.WaitGroup{},
	}
	return &r
}

func (r *Renderer) Render(camera *gogl.Camera) {
	for _, b := range r.batches {
		b.Render(camera)
	}
}
func (r *Renderer) Update(dt float32) {
	for _, b := range r.batches {
		r.wg.Add(1)
		go b.Update(dt, r.wg)
	}
	r.wg.Wait()
}

func (r *Renderer) AddGameObj(obj *components.GameObject) {
	added := false
	tex := obj.Spr.GetTexture()
	for _, b := range r.batches {
		if b.hasRoom && (b.GetZIndex() == obj.GetZIndex()) {
			if tex == nil || (b.HasTexture(tex) || b.HasTextureRoom()) {
				b.AddGameObject(obj)
				return
			}
		}
	}

	if !added {
		newBatch := NewRenderBatch(r.maxBatchSize, obj.GetZIndex())
		newBatch.Start()
		r.batches = append(r.batches, newBatch)
		newBatch.AddGameObject(obj)
		sort.Slice(r.batches, func(i, j int) bool {
			return r.batches[i].GetZIndex() < r.batches[j].GetZIndex()
		})
	}

}
