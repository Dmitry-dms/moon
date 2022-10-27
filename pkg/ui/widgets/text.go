package widgets

import (
	"github.com/Dmitry-dms/moon/pkg/fonts"
	"github.com/Dmitry-dms/moon/pkg/ui/styles"
)

type Text struct {
	base              baseWidget
	Message           string
	CurrentColor      [4]float32
	Chars             []fonts.CombinedCharInfo
	Selectable        bool
	Size              int
	Padding           int
	Scale             float32
	LastSelectedWidth float32
	LastSelectedX     float32
}

func NewText(id, text string, x, y, w, h float32, chars []fonts.CombinedCharInfo, style *styles.Style, sel bool) *Text {
	t := Text{
		Message: text,
		Chars:   chars,
		base: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h + float32(style.TextPadding)},
			backgroundColor: style.TransparentColor,
		},
		CurrentColor: style.TextColor,
		Size:         style.TextSize,
		Padding:      style.TextPadding * int(style.FontScale),
		Scale:        style.FontScale,
		Selectable:   sel,
	}
	return &t
}

func (t *Text) UpdatePosition(pos [4]float32) {
	t.base.boundingBox = pos
}

func (t *Text) FindSelectedString(x, dx float32) (float32, float32, string) {
	msg := ""
	var w, startPos float32 = 0, 0
	startFounded := false
	for _, pos := range t.Chars {
		if x >= pos.Pos.X && x <= pos.Pos.X+float32(pos.Char.Width) && !startFounded {
			w = pos.Width
			msg += string(pos.Char.Rune)
			startPos = pos.Pos.X
			startFounded = true
		} else {
			if startFounded && w < dx {
				msg += string(pos.Char.Rune)
				w += pos.Width
			}

		}
	}
	return startPos, w, msg
}

func (t *Text) SetWH(width, height float32) {
	t.base.boundingBox[2] = width
	t.base.boundingBox[3] = height + float32(t.Padding)
}

func (t *Text) SetBackGroundColor(clr [4]float32) {
	t.base.backgroundColor = clr
}

func (i *Text) BoundingBox() [4]float32 {
	return i.base.boundingBox
}
func (i *Text) BackgroundColor() [4]float32 {
	return i.base.backgroundColor
}
func (i *Text) Color() [4]float32 {
	return i.CurrentColor
}
func (i *Text) WidgetId() string {
	return i.base.id
}

func (i *Text) Height() float32 {
	return i.base.height()
}
func (i *Text) Visible() bool {
	return true
}
func (i *Text) Width() float32 {
	return i.base.width()
}
