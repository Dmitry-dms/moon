package widgets

import "github.com/Dmitry-dms/moon/pkg/gogl"

type Image struct {
	Id      string
	Texture *gogl.Texture
	BoundingBox  [4]float32 //x,y,w,h
	CurrentColor [4]float32
}

func (i Image) Rectangle() [4]float32 {
	return i.BoundingBox
}

func (i Image) GetColor() [4]float32 {
	return i.CurrentColor
}
func (i Image) GetId() string {
	return i.Id
}

func (i Image) Visible() bool {
	return true
}

func (i Image) Height() float32 {
	return i.BoundingBox[3] 
}
func (i Image) Width() float32 {
	return i.BoundingBox[2] 
}
