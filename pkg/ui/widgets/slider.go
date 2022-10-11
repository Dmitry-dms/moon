package widgets

import (
	"github.com/Dmitry-dms/moon/pkg/ui/styles"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
)

type Slider struct {
	base                   baseWidget
	min, max               float32
	sliderHeight, btnWidth float32
	mainSliderPos          [4]float32
	btnSliderPos           [4]float32
	currentPos             float32
	btnCaptured            bool
}

func NewSlider(id string, x, y, w, h, min, max float32, style *styles.Style) *Slider {
	s := Slider{
		base: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h},
			backgroundColor: [4]float32{123, 32, 12, 0},
		},
		min:          min,
		max:          max,
		sliderHeight: style.SliderHeight,
		btnWidth:     style.SliderBtnWidth,
	}
	return &s
}

func (s *Slider) HandleMouseDrag(delta float32, f utils.CapturedDragBehavior) {
	f(utils.NewRectS(s.btnSliderPos), &s.btnCaptured)
	if s.btnCaptured {

		s.currentPos += delta
		k := (s.Width() - s.btnSliderPos[2]) / s.max
		if s.currentPos <= s.min {
			s.currentPos = s.min
		} else if s.currentPos >= s.max*k {
			s.currentPos = s.max * k
		}
	}
}

func (s *Slider) CalculateNumber(i *float32) {
	baseLength := s.Width() - s.btnSliderPos[2]
	*i = s.max * (s.currentPos / baseLength)
}

func (s *Slider) calculateSliderPos() {
	v := (s.Height() - s.sliderHeight) / 2
	s.mainSliderPos = [4]float32{s.base.boundingBox[0], s.base.boundingBox[1] + v, s.Width(), s.sliderHeight}
	s.btnSliderPos = [4]float32{s.base.boundingBox[0] + s.currentPos, s.base.boundingBox[1], s.btnWidth, s.Height()}
}

func (s *Slider) MainSliderPos() [4]float32 {
	return s.mainSliderPos
}
func (s *Slider) BtnSliderPos() [4]float32 {
	return s.btnSliderPos
}
func (s *Slider) BoundingBox() [4]float32 {
	return s.base.boundingBox
}

func (s *Slider) UpdatePosition(pos [4]float32) {
	s.base.boundingBox = pos
	s.calculateSliderPos()
}

func (s *Slider) WidgetId() string {
	return s.base.id
}

func (s *Slider) Height() float32 {
	return s.base.height()
}
func (s *Slider) Width() float32 {
	return s.base.width()
}
