package draw

import (
	// "fmt"

	"github.com/Dmitry-dms/moon/pkg/fonts"
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/widgets"
)


type Command struct {
	Type     CmdType
	WidgetId string
	// priority int
	Shown    bool
	Rect     *Rect_command
	triangle *triangle_command
	RRect    *Rounded_rect
	Window   *Window_command
	WinStart *Window_start_command
	Toolbar  *Toolbar_command
	Text     *Text_command
}

type Text_command struct {
	Widget *widgets.Text
	Text   string
	Font   fonts.Font
	X, Y   float32
	Clr    [4]float32
	Id     string
	// Scale int
	Size  int
	TexId uint32
}

type Rect_command struct {
	X, Y, W, H  float32
	Clr         [4]float32
	Id          string
	ScaleFactor float32
	Texture     *gogl.Texture
}
type triangle_command struct {
	x0, y0, x1, y1, x2, y2 float32
	clr                    [4]float32
}

type Window_start_command struct {
	id   string
	x, y float32
}
type Window_command struct {
	Active     bool
	Id         string
	X, Y, W, H float32
	Clr        [4]float32
	Toolbar    Toolbar_command
}
type Toolbar_command struct {
	X, Y, W, H float32
	Clr        [4]float32
}
type Rounded_rect struct {
	X, Y, W, H float32
	Clr        [4]float32
	Radius     int
	Texture    *gogl.Texture
}

type CmdType int

const (
	RectType CmdType = iota
	RectTypeT
	Triangle
	TriangleT
	Line
	Circle
	RoundedRect
	RoundedRectT
	WindowCmd
	ToolbarCmd
	WindowStartCmd
	ScrollbarCmd
	ScrollButtonCmd
	Text
)