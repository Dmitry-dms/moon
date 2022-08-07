package widgets

import "github.com/Dmitry-dms/moon/pkg/ui/styles"

type Column struct {
	X, Y, W, H  float32
	CursorY     float32
	Padding     int
	ItemSpacing int
	style       *styles.Style
	Id          string
}

func NewColumn(id string, x, y float32, style *styles.Style) *Column {
	col := Column{
		X:     x,
		Y:     y,
		style: style,
		Id:    id,
	}
	return &col
}

func (col Column) Width() float32 {
	return col.W
}
func (col Column) WidgetId() string {
	return col.Id
}

func (col *Column) UpdateWidth(w float32) {
	if w > col.W {
		col.W = w
	}
}

func (col *Column) UpdatePosition(pos [4]float32) {
	col.X = pos[0]
	col.Y = pos[1]

	col.CursorY = pos[1]
}
