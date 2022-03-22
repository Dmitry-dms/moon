package listeners

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var MouseListener *mouseListener

func init() {
	o := sync.Once{}
	o.Do(func() {
		MouseListener = newMouseListener()
	})
}

type mouseListener struct {
	scrollX, scrollY         float64
	xPos, yPos, lastX, lastY float64
	mouseBtnPressed          [9]bool
	isDragging               bool
}

func newMouseListener() *mouseListener {
	listener := mouseListener{}
	return &listener
}

func MousePositionCallback(w *glfw.Window, xpos float64, ypos float64) {
	MouseListener.lastX = MouseListener.xPos
	MouseListener.lastY = MouseListener.yPos

	MouseListener.xPos = xpos
	MouseListener.yPos = ypos

	MouseListener.isDragging = MouseListener.mouseBtnPressed[0] ||
		MouseListener.mouseBtnPressed[1] ||
		MouseListener.mouseBtnPressed[2]
}

//mods - это сочетания клавиш
func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		//если клавишей на мышке больше, игнорируем нажатие
		if int(button) < len(MouseListener.mouseBtnPressed) {
			MouseListener.mouseBtnPressed[button] = true
		}
	case glfw.Release:
		if int(button) < len(MouseListener.mouseBtnPressed) {
			MouseListener.mouseBtnPressed[button] = false
			MouseListener.isDragging = false
		}
	}
}

func MouseScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	MouseListener.scrollX = xoff
	MouseListener.scrollY = yoff
}

func EndFrame() {
	MouseListener.scrollX, MouseListener.scrollY = 0, 0
	MouseListener.lastX, MouseListener.lastY = MouseListener.xPos, MouseListener.yPos
}

func GetX() float64 {
	return MouseListener.xPos
}
func GetY() float64 {
	return MouseListener.yPos
}

func GetDx() float64 {
	return (MouseListener.lastX - MouseListener.xPos)
}
func GetDy() float64 {
	return (MouseListener.lastY - MouseListener.yPos)
}

func GetScrollX() float64 {
	return MouseListener.scrollX
}
func GetScrollY() float64 {
	return MouseListener.scrollY
}
func IsDragging() bool {
	return MouseListener.isDragging
}

func MouseButtonDown(button glfw.MouseButton) bool {
	if int(button) < len(MouseListener.mouseBtnPressed) {
		return MouseListener.mouseBtnPressed[button]
	} else {
		return false
	}
}
