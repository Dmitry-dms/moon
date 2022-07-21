package draw

import (
	// "fmt"

	"github.com/Dmitry-dms/moon/pkg/fonts"
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/widgets"
)

const MAX_COMMANDS_COUNT = 1000

type RenderQueue struct {
	commands     []Command
	bufferWindow []Command

	CmdCount int
	Zindex   int
}

func NewRenderQueue() *RenderQueue {
	q := RenderQueue{
		commands: make([]Command, 0),
	}
	return &q
}

func (r *RenderQueue) AddCommand(cmd Command) {
	// fmt.Printf("%s \n", cmd.t)
	if cmd.Type == WindowStartCmd {
	}
	if cmd.Type == WindowCmd {
		r.commands = append(r.commands, r.bufferWindow...)
		// fmt.Println("End window ", r.CmdCount)

		return
	}
	if cmd.Type == ScrollButtonCmd || cmd.Type == ScrollbarCmd {
		if cmd.Shown {
			cmd.Type = RoundedRect
		}
	}
	r.bufferWindow = append(r.bufferWindow, cmd)
	r.CmdCount++
}

func (r *RenderQueue) Commands() []Command {
	return r.commands
}

func (r *RenderQueue) clearCommands() {
	r.commands = []Command{}
	r.bufferWindow = []Command{}
	r.CmdCount = 0
}

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
