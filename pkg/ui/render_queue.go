package ui



const MAX_COMMANDS_COUNT = 1000

type RenderQueue struct {
	commands []command
	CmdCount int
}

func NewRenderQueue() *RenderQueue {
	q := RenderQueue{
		commands: make([]command, MAX_COMMANDS_COUNT),
	}
	return &q
}

func (r *RenderQueue) AddCommand(cmd command) {
	c := r.CmdCount
	r.commands[c] = cmd
	r.CmdCount++
}

func (r *RenderQueue) Commands() []command {
	return r.commands
}

func (r *RenderQueue) clearCommands() {
	for i, _ := range r.commands {
		r.commands[i] = command{}
	}
	r.CmdCount = 0
}

type command struct {
	t        CmdType
	rect     *rect_command
	triangle *triangle_command
	rRect    *rounded_rect
}

type rect_command struct {
	x, y, w, h float32
	clr        [4]float32
}
type triangle_command struct {
	x0, y0, x1, y1, x2, y2 float32
	clr                    [4]float32
}
type rounded_rect struct {
	x, y, w, h float32
	clr        [4]float32
	radius     int
}

type CmdType int

const (
	Rect CmdType = iota
	Triangle
	Line
	Circle
	RoundedRect
)
