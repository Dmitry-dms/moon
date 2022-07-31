package widgets

type Widget interface {
	GetColor() [4]float32
	GetId() string
	Rectangle() [4]float32 // x,y,w,h
	Visible() bool
}

type WidgetType int

const (
	ImageWidget WidgetType = iota
	ButtonWidget
	VerticalSpacingWidget
	TextWidget
)