package widgets

type Widget interface {
	WidgetId() string
	UpdatePosition([4]float32)
	Height() float32
}

type WidgetType int

const (
	ImageWidget WidgetType = iota
	ButtonWidget
	VerticalSpacingWidget
	TextWidget
)
