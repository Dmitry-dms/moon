package widgets

type Widget interface {
	// Color() [4]float32
	WidgetId() string
	UpdatePosition([4]float32)
	// Rectangle() [4]float32 // x,y,w,h
	// Visible() bool
	Height() float32
}

type WidgetType int

const (
	ImageWidget WidgetType = iota
	ButtonWidget
	VerticalSpacingWidget
	TextWidget
)
