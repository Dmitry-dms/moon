package widgets

type VSpace struct {
	BoundingBox [4]float32 //x,y,w,h
}

func (s VSpace) Rectangle() [4]float32 {
	return s.BoundingBox
}
func (s VSpace) Visible() bool {
	return true
}
func (s VSpace) Height() float32 {
	return s.BoundingBox[3]
}

func (s VSpace) GetColor() [4]float32 {
	return [4]float32{}
}
func (s VSpace) GetId() string {
	return "vert_spacing"
}
