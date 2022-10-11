package widgets

type Widget interface {
	WidgetId() string
	UpdatePosition([4]float32)
	Height() float32
}

type baseWidget struct {
	id              string
	boundingBox     [4]float32
	backgroundColor [4]float32
}

func (b *baseWidget) height() float32 {
	return b.boundingBox[3]
}
func (b *baseWidget) width() float32 {
	return b.boundingBox[2]
}

type WidgetType int

const (
	ImageWidget WidgetType = iota
	ButtonWidget
	VerticalSpacingWidget
	TextWidget
)
