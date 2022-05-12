package listeners

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
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

func SizeCllback(w *glfw.Window, width int32, height int32) {
	mWinListener.width = width
	mWinListener.height = height
}

func SetWinHeight(h int32) {
	mWinListener.height = h
}
func SetWinWidth(w int32) {
	mWinListener.width = w
}
