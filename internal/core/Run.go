package core

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/inkyblackness/imgui-go/v4"
)

// Platform covers mouse/keyboard/gamepad inputs, cursor shape, timing, windowing.
type Platform interface {
	// ShouldStop is regularly called as the abort condition for the program loop.
	ShouldStop() bool
	// ProcessEvents is called once per render loop to dispatch any pending events.
	ProcessEvents()
	// DisplaySize returns the dimension of the display.
	DisplaySize() [2]float32
	// FramebufferSize returns the dimension of the framebuffer.
	FramebufferSize() [2]float32
	// NewFrame marks the begin of a render pass. It must update the imgui IO state according to user input (mouse, keyboard, ...)
	NewFrame()
	// PostRender marks the completion of one render pass. Typically this causes the display buffer to be swapped.
	PostRender()
	// ClipboardText returns the current text of the clipboard, if available.
	ClipboardText() (string, error)
	// SetClipboardText sets the text as the current text of the clipboard.
	SetClipboardText(text string)

	GetWindow() *glfw.Window

	Dispose()
}

// type clipboard struct {
// 	platform Platform
// }

// func (board clipboard) Text() (string, error) {
// 	return board.platform.ClipboardText()
// }

// func (board clipboard) SetText(text string) {
// 	board.platform.SetClipboardText(text)
// }

// Renderer covers rendering imgui draw data.
type Renderer interface {
	// PreRender causes the display buffer to be prepared for new output.
	PreRender(clearColor [3]float32)
	// Render draws the provided imgui draw data.
	Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData)
	Dispose()
}
