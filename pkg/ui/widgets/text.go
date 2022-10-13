package widgets

import "github.com/Dmitry-dms/moon/pkg/ui/styles"

type Text struct {
	base         baseWidget
	Message      string
	CurrentColor [4]float32

	Size    int
	Padding int
	Scale   float32
}

func NewText(id, text string, x, y, w, h float32, style *styles.Style) *Text {
	t := Text{
		Message: text,
		base: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h + float32(style.TextPadding)},
			backgroundColor: style.TransparentColor,
		},
		CurrentColor: style.TextColor,
		Size:         style.TextSize,
		Padding:      style.TextPadding * int(style.FontScale),
		Scale:        style.FontScale,
	}
	return &t
}

func (t *Text) UpdatePosition(pos [4]float32) {
	t.base.boundingBox = pos
}

func (t *Text) SetWH(width, height float32) {
	t.base.boundingBox[2] = width
	t.base.boundingBox[3] = height + float32(t.Padding)
}

func (t *Text) SetBackGroundColor(clr [4]float32) {
	t.base.backgroundColor = clr
}

func (i *Text) Rectangle() [4]float32 {
	return i.base.boundingBox
}
func (i *Text) BackgroundColor() [4]float32 {
	return i.base.backgroundColor
}
func (i *Text) Color() [4]float32 {
	return i.CurrentColor
}
func (i *Text) WidgetId() string {
	return i.base.id
}

func (i *Text) Height() float32 {
	return i.base.height()
}
func (i *Text) Visible() bool {
	return true
}
func (i *Text) Width() float32 {
	return i.base.width()
}
