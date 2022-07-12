package ui

import (
	// "fmt"

	"github.com/Dmitry-dms/moon/pkg/gogl"
)

const MAX_COMMANDS_COUNT = 1000

type RenderQueue struct {
	commands     []command
	bufferWindow []command

	CmdCount int
	Zindex   int
}

func NewRenderQueue() *RenderQueue {
	q := RenderQueue{
		commands: make([]command, 0),
	}
	return &q
}

func (r *RenderQueue) AddCommand(cmd command) {
	// fmt.Printf("%s \n", cmd.t)
	if cmd.t == WindowStartCmd {
	}
	if cmd.t == WindowCmd {
		r.commands = append(r.commands, r.bufferWindow...)
		// fmt.Println("End window ", r.CmdCount)
		
		return
	}
	if cmd.t == ScrollButtonCmd || cmd.t == ScrollbarCmd {
		if cmd.shown {
			cmd.t = RoundedRect
		}
	}
	r.bufferWindow = append(r.bufferWindow, cmd)
	r.CmdCount++
}

func (r *RenderQueue) Commands() []command {
	return r.commands
}

func (r *RenderQueue) clearCommands() {
	r.commands = []command{}
	r.bufferWindow = []command{}
	r.CmdCount = 0
}

type command struct {
	t        CmdType
	priority int
	shown    bool
	rect     *rect_command
	triangle *triangle_command
	rRect    *rounded_rect
	window   *window_command
	winStart *window_start_command
}

type rect_command struct {
	x, y, w, h float32
	clr        [4]float32
	id         string
	scaleFactor float32
	texture    *gogl.Texture
}
type triangle_command struct {
	x0, y0, x1, y1, x2, y2 float32
	clr                    [4]float32
}

type window_start_command struct {
	id   string
	x, y float32
}
type window_command struct {
	active     bool
	id         string
	x, y, w, h float32
	clr        [4]float32
	toolbar    toolbar_command
}
type toolbar_command struct {
	h   float32
	clr [4]float32
}
type rounded_rect struct {
	x, y, w, h float32
	clr        [4]float32
	radius     int
	texture    *gogl.Texture
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
)
