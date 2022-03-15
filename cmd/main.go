package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Dmitry-dms/moon/internal/core"
	"github.com/Dmitry-dms/moon/internal/platforms"
	"github.com/pkg/errors"

)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	core, err := core.NewCore(1200, 720, platforms.GLFWClientAPIOpenGL42)
	if err != nil {
		fmt.Println(errors.Unwrap(err))
		os.Exit(1)
	}
	defer core.Dispose()

	//Main loop
	core.Run()
}
