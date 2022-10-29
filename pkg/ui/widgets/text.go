package widgets

import (
	"fmt"
	"github.com/Dmitry-dms/moon/pkg/fonts"
	"github.com/Dmitry-dms/moon/pkg/ui/styles"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
)

type Text struct {
	base              baseWidget
	Message           string
	CurrentColor      [4]float32
	Chars             []fonts.CombinedCharInfo
	Flag              TextFlag
	Size              int
	Padding           int
	Scale             float32
	LastSelectedWidth float32
	LastSelectedX     float32

	startX, startY      float32
	endX, endY          float32
	startInd, startLine int
	endInd, endLine     int
}

type TextFlag uint

const (
	Selectable TextFlag = 1 << iota
	SplitWords
	SplitChars
	FitContent
	Editable
	Default

	//FitWords = FitContent | SplitWords
	//FitChars = FitContent | SplitChars
)

func NewText(id, text string, x, y, w, h float32, chars []fonts.CombinedCharInfo, style *styles.Style, flag TextFlag) *Text {
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
		Flag:         flag,
	}
	return &t
}

func (t *Text) UpdatePosition(pos [4]float32) {
	t.base.boundingBox = pos
}
func (t *Text) FindSelectedString2(x, y, dx, dy float32, lines []fonts.TextLine, gh bool) utils.Rect {
	//fmt.Println(x, y, y+dy)
	//startLineInd := 0
	//for i, line := range lines {
	//	if y+dy <= line.Height {
	//		startLineInd = i
	//		break
	//	}
	//}

	msg := ""
	var w, startPos float32 = 0, 0
	type inf struct {
		msg        string
		start, end float32
	}

	_ = w
	tmp := []inf{}
	startFounded := false
	endFounded := false
	//ind := 0
	for i := 0; i <= len(lines)-1; i++ {
		for ind, pos := range lines[i].Text {
			if !gh {
				if x >= pos.Pos.X && x <= pos.Pos.X+float32(pos.Char.Width) && !startFounded &&
					y+dy <= lines[i].StartY+lines[i].Height {
					t.startX = pos.Pos.X
					t.startY = lines[i].StartY
					t.startInd = ind
					t.startLine = i
					startFounded = true
				}
			} else {
				if x+dx >= pos.Pos.X && x+dx <= pos.Pos.X+float32(pos.Char.Width) &&
					y+dy <= lines[i].StartY+lines[i].Height && !endFounded { // start y был Height, который закоменчен
					//fmt.Println(y+dy, lines[i].StartY, i)
					//fmt.Println(y, dy, pos.Char.Height)
					//fmt.Println("here", x+dx, ind, i)
					t.endX = x + dx
					t.endInd = ind
					t.endLine = i
					endFounded = true
				}
			}
		}
		//if i != len(lines) && y+dy >= lines[0].Height {
		//	tmp = append(tmp, inf{start: startPos, msg: lines[i].Msg[ind:]})
		//} else {
		tmp = append(tmp, inf{start: startPos, msg: msg})
		//}
	}
	//fmt.Println(string(lines[t.startLine].Text[t.startInd].Char.Rune),
	//	string(lines[t.endLine].Text[t.endInd].Char.Rune))
	//fmt.Println(t.startLine, t.startInd, t.endLine, t.endInd)
	var x1, y1, w1, h1 float32
	stre := ""
	if t.startLine == t.endLine {
		x1 = t.startX
		y1 = t.startY
		h1 = lines[t.startLine].Height
		for i := t.startInd; i <= t.endInd; i++ {
			c := lines[t.startLine].Text[i]
			stre += string(c.Char.Rune)
			w1 += float32(c.Width)
		}
	}
	fmt.Println(stre)
	r := utils.NewRect(x1, y1, w1, h1)
	return r

	//fmt.Println(x, y, dx, dy)
	//msg := ""
	//var w, startPos float32 = 0, 0
	//startFounded := false
	//for _, pos := range t.Chars {
	//	if x >= pos.Pos.X && x <= pos.Pos.X+float32(pos.Char.Width) && !startFounded {
	//		w = pos.Width
	//		msg += string(pos.Char.Rune)
	//		startPos = pos.Pos.X
	//		startFounded = true
	//	} else {
	//		if startFounded && w < dx {
	//			msg += string(pos.Char.Rune)
	//			w += pos.Width
	//		}
	//	}
	//}
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
