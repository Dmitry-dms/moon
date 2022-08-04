package widgets

type Text struct {
	Message         string
	Id              string

	BoundingBox     [4]float32 //x,y,w,h
	CurrentColor    [4]float32
	BackgroundColor [4]float32

	Size int
}

func NewText(id, text string, x, y, w, h float32, size int, clr [4]float32) *Text {
	t := Text{
		Message:      text,
		Id:           id,
		BoundingBox:  [4]float32{x, y, w, h},
		CurrentColor: clr,
		Size: size,
	}
	return &t
}

func (t *Text) UpdatePosition(pos [4]float32) {
	t.BoundingBox = pos
}

func (t *Text) SetBackGroundColor(clr [4]float32) {
	t.BackgroundColor = clr
}

func (i Text) Rectangle() [4]float32 {
	return i.BoundingBox
}

func (i Text) Color() [4]float32 {
	return i.CurrentColor
}
func (i Text) WidgetId() string {
	return i.Id
}

func (i Text) Height() float32 {
	return i.BoundingBox[3]
}
func (i Text) Visible() bool {
	return true
}
func (i Text) Width() float32 {
	return i.BoundingBox[2]
}
