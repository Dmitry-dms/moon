package listeners

import (
	"fmt"
	"sync"

	"github.com/go-gl/gl/v4.2-core/gl"
	// "github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	s := sync.Once{}
	s.Do(func() {
		mWinListener = &windowListener{}
	})
}

var mWinListener *windowListener

type windowListener struct {
	width, height int32

}

func GetWindowWidth() int32 {
	return mWinListener.width
}
func GetWindowHeight() int32 {
	return mWinListener.height
}

func SizeCllback(w *glfw.Window, width int32, height int32) {
	mWinListener.width = width
	mWinListener.height = height

	gl.Viewport(0,0,width,height)

	

	s := mgl32.Vec2{float32(width),float32(height)}
	posx,posy := w.GetPos()
	
	s2 := mgl32.Vec2{float32(posx),float32(posy)}
	SetGameViewPortSize(s)
	SetGameViewPortPos(s2)
	
	fmt.Println(width,height)
}

func SetWinHeight(h int32) {
	mWinListener.height = h
}
func SetWinWidth(w int32) {
	mWinListener.width = w
}
