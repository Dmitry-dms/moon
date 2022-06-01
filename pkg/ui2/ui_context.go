package ui2

import (
	// "io/ioutil"
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	// "os"

	// "github.com/go-gl/gltext"
	"github.com/Dmitry-dms/moon/pkg/ui2/fonts"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	// "github.com/golang/freetype"
	// "github.com/golang/freetype/truetype"

	"github.com/pkg/errors"
)

var context *UIContext = nil

type UIContext struct {
	Initialized      bool
	Io               *ImIO
	Viewports        []*ImGuiViewportP
	InputEventsQueue []ImGuiInputEvent

	MouseCursor ImGuiMouseCursor

	//new frame
	Hooks                []*ImGuiContextHook[any]
	Time                 float32
	WithinFrameScope     bool
	FrameCount           uint
	WindowsActiveCount   uint
	TooltipOverrideCount float64

	//font
	DefaultFont *truetype.Font
}

func LoadFontFromFile(path string, scale int32) (*truetype.Font,error) {
	// fd, err := os.Open(path)
	// if err != nil {
	// 	return nil, err
	// }

	// defer fd.Close()

	// return gltext.LoadTruetype(fd, scale, 32, 127, gltext.LeftToRight)
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}
	// makeImageFont(f)

	fonts.NewFont(path, 12)

	return f,nil
	//load font (fontfile, font scale, window width, window height
	// font, err := glfont.LoadFont(path, int32(52), 1920, 1080)
	// if err != nil {
	// 	panic(err)
	// }
	// return font, nil

}
var text = []string{
	"’Twas brillig, and the slithy toves",
	"Did gyre and gimble in the wabe;",
	"All mimsy were the borogoves,",
	"And the mome raths outgrabe.",
	"",
	"“Beware the Jabberwock, my son!",
	"The jaws that bite, the claws that catch!",
	"Beware the Jubjub bird, and shun",
	"The frumious Bandersnatch!”",
	"",
	"He took his vorpal sword in hand:",
	"Long time the manxome foe he sought—",
	"So rested he by the Tumtum tree,",
	"And stood awhile in thought.",
	"",
	"And as in uffish thought he stood,",
	"The Jabberwock, with eyes of flame,",
	"Came whiffling through the tulgey wood,",
	"And burbled as it came!",
	"",
	"One, two! One, two! and through and through",
	"The vorpal blade went snicker-snack!",
	"He left it dead, and with its head",
	"He went galumphing back.",
	"",
	"“And hast thou slain the Jabberwock?",
	"Come to my arms, my beamish boy!",
	"O frabjous day! Callooh! Callay!”",
	"He chortled in his joy.",
	"",
	"’Twas brillig, and the slithy toves",
	"Did gyre and gimble in the wabe;",
	"All mimsy were the borogoves,",
	"And the mome raths outgrabe.",
}
var size float64 = 12
var spacing   float64 = 1.5

func makeImageFont(f *truetype.Font) {
	// Initialize the context.
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	// if *wonb {
	// 	fg, bg = image.White, image.Black
	// 	ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	// }
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(144)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	hinting := "none"
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(12)>>6))
	for _, s := range text {
		_, err := c.DrawString(s, pt)
		if err != nil {
			fmt.Println(err)
			return
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}

func GetCurrentContext() *UIContext {
	return context
}

type ImGuiMouseCursor int

const (
	ImGuiMouseCursor_None ImGuiMouseCursor = iota
	ImGuiMouseCursor_Arrow
	ImGuiMouseCursor_TextInput
	ImGuiMouseCursor_ResizeAll
	ImGuiMouseCursor_ResizeNS
	ImGuiMouseCursor_ResizeEW
	ImGuiMouseCursor_ResizeNESW
	ImGuiMouseCursor_ResizeNWSE
	ImGuiMouseCursor_Hand
	ImGuiMouseCursor_NotAllowed
	ImGuiMouseCursor_COUNT
)

func initializeContext() {
	viewport := NewImGuiViewportP()
	context.Viewports = append(context.Viewports, viewport)

}

func (i *UIContext) GetMouseCursor() ImGuiMouseCursor {
	return i.MouseCursor
}

func (i *UIContext) pushEvent(e ImGuiInputEvent) {
	i.InputEventsQueue = append(i.InputEventsQueue, e)
}

func CreateContext() *UIContext {
	prevCtx := GetCurrentContext()
	if prevCtx == nil {
		ctx := &UIContext{
			Initialized: true,
			Io: &ImIO{
				DisplaySize:             Vec2{},
				BackendPlatformUserData: newData(),
			},
			Viewports:        make([]*ImGuiViewportP, 0),
			InputEventsQueue: make([]ImGuiInputEvent, 0),
			Hooks:            make([]*ImGuiContextHook[any], 0),
		}
		context = ctx
	}
	initializeContext()
	return context
}

func (c *UIContext) CallContextHooks(hook_type ImGuiContextHookType) {
	for _, hook := range c.Hooks {
		if hook.Type == hook_type {
			hook.Callback(c, hook)
		}
	}
}

func (c *UIContext) NewFrame() {
	// Remove pending delete hooks before frame start.
	// This deferred removal avoid issues of removal while iterating the hook vector
	for i, hook := range c.Hooks {
		if hook.Type == ImGuiContextHookType_PendingRemoval_ {
			c.Hooks = removeIndex(c.Hooks, i)
		}
	}
	c.CallContextHooks(ImGuiContextHookType_NewFramePre)

	err := ErrorCheckNewFrameSanityChecks()
	if err != nil {
		panic(err)
	}

	c.Time += c.Io.DeltaTime
	c.WithinFrameScope = true
	c.FrameCount += 1
	c.TooltipOverrideCount = 0
	c.WindowsActiveCount = 0

	c.UpdateViewportsNewFrame()

}

//проверка на ошибки
func ErrorCheckNewFrameSanityChecks() error {
	ctx := GetCurrentContext()

	if !ctx.Initialized {
		return errors.New("Gui is not initialized")
	}
	if ctx.Io.DeltaTime > 0 || ctx.FrameCount == 0 {
		return errors.New("Need a positive DeltaTime!")
	}

	return nil
}
func (c *UIContext) UpdateViewportsNewFrame() {
	if len(c.Viewports) != 1 {
		panic("Viewports len != 1")
	}
	// Update main viewport with current platform position.
	main_viewport := c.Viewports[0]
	main_viewport.Flags = ImGuiViewportFlags_IsPlatformWindow | ImGuiViewportFlags_OwnedByApp
	main_viewport.Pos = Vec2{0, 0}
	main_viewport.Size = c.Io.DisplaySize

	for _, v := range c.Viewports {
		v.WorkOffsetMin = v.BuildWorkOffsetMin
		v.WorkOffsetMax = v.BuildWorkOffsetMax
		v.BuildWorkOffsetMax, v.BuildWorkOffsetMin = Vec2{0, 0}, Vec2{0, 0}
		v.UpdateWorkRect()
	}
}

func removeIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
