package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Dmitry-dms/moon/internal/example"
	"github.com/Dmitry-dms/moon/internal/platforms"
	"github.com/Dmitry-dms/moon/internal/renderers"
	//"github.com/go-gl/glfw/v3.3/glfw"
	imgui "github.com/inkyblackness/imgui-go/v4"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()
	
	
	platform, err := platforms.NewGLFW(io, platforms.GLFWClientAPIOpenGL42)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	renderer, err := renderers.NewOpenGL42(io)
	// renderer, err := renderers.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	example.Run(platform, renderer)
}
