package widgets

// TODO: add Base Widget
type Button struct {
	IsActive     bool
	CurrentColor [4]float32
	Id           string
	BoundingBox  [4]float32 //x,y,w,h
}

func NewButton(id string, x, y, w, h float32, backClr [4]float32) *Button {
	btn := Button{
		IsActive:     false,
		CurrentColor: backClr,
		Id:           id,
		BoundingBox:  [4]float32{x, y, w, h},
	}
	return &btn
}
func (b *Button) UpdatePosition(pos [4]float32) {
	b.BoundingBox = pos
}

func (b *Button) AddWidth(w float32) {
	b.BoundingBox[2] += w
}
func (b *Button) AddHeight(h float32) {
	b.BoundingBox[3] += h
}

func (b Button) WidgetId() string {
	return b.Id
}

func (b Button) Height() float32 {
	return b.BoundingBox[3]
}
func (b Button) Width() float32 {
	return b.BoundingBox[2]
}

func (b Button) Rectangle() [4]float32 {
	return b.BoundingBox
}

func (b *Button) SetColor(clr [4]float32) {
	b.CurrentColor = clr
}
