package listeners

import (

	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var KeyListener *keyListener

func init() {
	o := sync.Once{}
	o.Do(func() {
		KeyListener = newKeyListener()
	})
}

type keyListener struct {
	keyPressed [350]bool
}

func newKeyListener() *keyListener {
	listener := keyListener{}
	return &listener
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		KeyListener.keyPressed[key] = true
	case glfw.Release:
		KeyListener.keyPressed[key] = false
	}
}

func IsKeyPressed(key glfw.Key) bool {
	return KeyListener.keyPressed[key]
}
