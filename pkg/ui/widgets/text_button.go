package widgets

import "github.com/Dmitry-dms/moon/pkg/ui/styles"

type TextButton struct {
	Id    string
	align TextAlign
	style *styles.Style
	*Text
	*Button
}

type TextAlign uint32

const (
	Center TextAlign = iota
	Left
	Right
)

func NewTextButton(id string, x, y, w, h float32,
	text string, align TextAlign, style *styles.Style) *TextButton {
	tb := TextButton{
		Id:    id,
		align: align,
		style: style,
	}
	txt := NewText("", text, x, y, w, h, style)
	btn := NewButton("", x, y, 2*style.Padding+w, 2*style.Padding+h, style.BtnColor)
	tb.Button = btn
	tb.Text = txt

	tb.UpdateTextPos(tb.Button.BoundingBox()[0], tb.Button.BoundingBox()[1])

	return &tb
}

func (tb *TextButton) UpdateTextPos(x, y float32) {
	// x := tb.Button.BoundingBox[0]
	// y := tb.Button.BoundingBox[1]
	switch tb.align {
	case Center:
		tb.Text.UpdatePosition([4]float32{x + tb.style.Padding, y + tb.style.Padding, tb.Text.Width(), tb.Text.Height()})
	case Left:
		tb.Text.UpdatePosition([4]float32{x, y + tb.style.Padding, tb.Text.Width(), tb.Text.Height()})
	case Right:
		tb.Text.UpdatePosition([4]float32{x + tb.Button.Width() - tb.Text.Width(), y + tb.style.Padding, tb.Text.Width(), tb.Text.Height()})
	}
}

func (tb *TextButton) Active() bool {
	return tb.Button.IsActive
}

func (tb *TextButton) ChangeActive() {
	tb.Button.ChangeActive()
}

func (tb *TextButton) UpdatePosition(pos [4]float32) {
	tb.Button.UpdatePosition(pos)
	tb.UpdateTextPos(tb.Button.BoundingBox()[0], tb.Button.BoundingBox()[1])
}
func (tb *TextButton) WidgetId() string {
	return tb.Id
}
func (tb *TextButton) Height() float32 {
	return tb.Button.Height()
}

func (tb *TextButton) BoundingBox() [4]float32 {
	return tb.Button.BoundingBox()
}
func (tb *TextButton) Width() float32 {
	return tb.Button.Width()
}
