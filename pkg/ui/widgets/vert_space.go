package widgets

type VSpace struct {
	BoundingBox [4]float32 //x,y,w,h
	Id          string
}

func (s *VSpace) UpdatePosition(pos [4]float32) {
	s.BoundingBox = pos
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
func (s VSpace) Width() float32 {
	return s.BoundingBox[2]
}

func (s VSpace) Color() [4]float32 {
	return [4]float32{}
}
func (s VSpace) WidgetId() string {
	return s.Id
}
