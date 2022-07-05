package ui

const MAX_COMMANDS_COUNT = 1000

type RenderQueue struct {
	commands     []command
	bufferWindow []command
	activeBuffer []command
	lastBuffer   []command
	CmdCount     int
	Zindex       int
	Windows      map[string]pair
	LastPar      pair
}
type pair struct {
	x, y int
}

func NewRenderQueue() *RenderQueue {
	q := RenderQueue{
		commands: make([]command, 0),
		Windows:  make(map[string]pair),
	}
	return &q
}

func (r *RenderQueue) AddCommand(cmd command) {
	if cmd.t == WindowStartCmd {
	}
	if cmd.t == WindowCmd {
		r.commands = append(r.commands, r.bufferWindow...)
		return
	}
	r.bufferWindow = append(r.bufferWindow, cmd)
	r.CmdCount++
}

func (r *RenderQueue) Commands() []command {
	r.commands = append(r.commands, r.activeBuffer...)
	return r.commands
}

func (r *RenderQueue) clearCommands() {
	r.commands = []command{}
	r.bufferWindow = []command{}
	r.activeBuffer = []command{}
	r.CmdCount = 0
}

type command struct {
	t        CmdType
	priority int
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
}

type CmdType int

const (
	RectType CmdType = iota
	Triangle
	Line
	Circle
	RoundedRect
	WindowCmd
	ToolbarCmd
	WindowStartCmd
)
