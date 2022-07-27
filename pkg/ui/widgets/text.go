package widgets

type Text struct {
	Message      string
	Id           string
	Texture      uint32
	BoundingBox  [4]float32 //x,y,w,h
	CurrentColor [4]float32
}

func (i Text) Rectangle() [4]float32 {
	return i.BoundingBox
}

func (i Text) GetColor() [4]float32 {
	return i.CurrentColor
}
func (i Text) GetId() string {
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
