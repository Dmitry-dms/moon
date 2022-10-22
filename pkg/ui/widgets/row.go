package widgets

import (
	"github.com/Dmitry-dms/moon/pkg/ui/styles"
)

type HybridLayout struct {
	X, Y, W, H            float32
	CursorX, CursorY      float32
	InitY                 float32
	LastWidth, LastHeight float32

	RequiereColumn           bool
	CurrentColH, CurrentColW float32
	// CurrentCol     *Column

	// ColumnW, ColumnH float32
	Padding     int
	ItemSpacing int
	style       *styles.Style
	Id          string
}

type Column struct {
	Id   string
	W, H float32
}

func (c *HybridLayout) UpdateColWidth(w float32) {
	if c.CurrentColW < w {
		c.CurrentColW = w
	}
}
func (c *HybridLayout) AddColHeight(h float32) {
	c.CurrentColH += h
}

func NewHLayout(id string, x, y float32, style *styles.Style) *HybridLayout {
	r := HybridLayout{
		X:     x,
		Y:     y,
		style: style,
		Id:    id,
	}
	return &r
}

func (r *HybridLayout) Height() float32 {
	return r.LastHeight
}
func (r *HybridLayout) Width() float32 {
	return r.LastWidth
}
func (r *HybridLayout) WidgetId() string {
	return r.Id
}

func (r *HybridLayout) UpdateHeight(h float32) {
	if h > r.H {
		r.H = h
	}
}

func (r *HybridLayout) UpdatePosition(pos [4]float32) {
	r.X = pos[0]
	r.Y = pos[1]

	r.CursorX = pos[0]
	r.CursorY = pos[1]

	r.InitY = pos[1]
}

func (r *HybridLayout) BoundingBox() [4]float32 {
	return [4]float32{}
}
