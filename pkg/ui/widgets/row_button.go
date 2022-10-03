package widgets

import "github.com/Dmitry-dms/moon/pkg/ui/styles"

type RowButton struct {
	Id  string
	btn *Button
}

func NewRowButton(id string, x, y, w float32, style *styles.Style) *RowButton {
	rb := RowButton{
		Id:  id,
		btn: NewButton(id, x, y, w, 0, style.BtnColor),
	}

	return &rb
}

func (br *RowButton) UpdatePosition(pos [4]float32) {
	br.btn.BoundingBox = pos
}
func (br RowButton) WidgetId() string {
	return br.Id
}

func (br RowButton) Height() float32 {
	return br.btn.Height()
}
