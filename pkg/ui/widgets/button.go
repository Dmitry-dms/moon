package widgets


// TODO: add Base Widget
type Button struct {
	IsActive     bool
	CurrentColor [4]float32
	Id           string
	BoundingBox  [4]float32 //x,y,w,h
}

func (b Button) Visible() bool {
	return true
}

func (b *Button) AddWidth(w float32) {
	b.BoundingBox[2] += w
}
func (b *Button) AddHeight(h float32) {
	b.BoundingBox[3] += h
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

func (b Button) GetColor() [4]float32 {
	return b.CurrentColor
}
func (b Button) GetId() string {
	return b.Id
}
func (b *Button) SetColor(clr [4]float32) {
	b.CurrentColor = clr
}
