package main

import (
	"runtime"
	"github.com/Dmitry-dms/moon/internal/core"

)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}


func main() {
	// core, err := core.NewCore(1200, 720, platforms.GLFWClientAPIOpenGL42, 0)
	// if err != nil {
	// 	fmt.Println(errors.Unwrap(err))
	// 	os.Exit(1)
	// }
	//core.Window.Run()
	defer core.Window.Dispose()

	//Main loop
	core.Window.Run()
}
