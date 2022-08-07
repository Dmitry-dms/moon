package widgets

import "github.com/Dmitry-dms/moon/pkg/ui/styles"

type Row struct {
	X, Y, W, H  float32
	CursorX     float32
	Padding     int
	ItemSpacing int
	style       *styles.Style
	Id          string
}

func NewRow(id string, x, y float32, style *styles.Style) *Row {
	r := Row{
		X:     x,
		Y:     y,
		style: style,
		Id:    id,
	}
	return &r
}

func (r Row) Height() float32 {
	return r.H
}
func (r Row) WidgetId() string {
	return r.Id
}

func (r *Row) UpdateHeight(h float32) {
	if h > r.H {
		r.H = h
	}
}

func (r *Row) UpdatePosition(pos [4]float32) {
	r.X = pos[0]
	r.Y = pos[1]

	r.CursorX = r.X
}
