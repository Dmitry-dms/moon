package listeners

import (
	"sync"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var mMouseListener *mouseListener

func init() {
	o := sync.Once{}
	o.Do(func() {
		mMouseListener = newMouseListener()
	})
}

type mouseListener struct {
	scrollX, scrollY         float64
	xPos, yPos, lastX, lastY float64
	mouseBtnPressed          [9]bool
	isDragging               bool
	camera                   *gogl.Camera
}

func newMouseListener() *mouseListener {
	listener := mouseListener{}
	return &listener
}

func MousePositionCallback(w *glfw.Window, xpos float64, ypos float64) {
	mMouseListener.lastX = mMouseListener.xPos
	mMouseListener.lastY = mMouseListener.yPos

	mMouseListener.xPos = xpos
	mMouseListener.yPos = ypos

	mMouseListener.isDragging = mMouseListener.mouseBtnPressed[0] ||
		mMouseListener.mouseBtnPressed[1] ||
		mMouseListener.mouseBtnPressed[2]
}

//mods - это сочетания клавиш
func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		//если кнопок на мышке больше, игнорируем нажатие
		if int(button) < len(mMouseListener.mouseBtnPressed) {
			mMouseListener.mouseBtnPressed[button] = true
		}
	case glfw.Release:
		if int(button) < len(mMouseListener.mouseBtnPressed) {
			mMouseListener.mouseBtnPressed[button] = false
			mMouseListener.isDragging = false
		}
	}
}

func MouseScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	mMouseListener.scrollX = xoff
	mMouseListener.scrollY = yoff
}

func EndFrame() {
	mMouseListener.scrollX, mMouseListener.scrollY = 0, 0
	mMouseListener.lastX, mMouseListener.lastY = mMouseListener.xPos, mMouseListener.yPos
}

func GetX() float64 {
	return mMouseListener.xPos
}
func GetY() float64 {
	return mMouseListener.yPos
}

func GetDx() float64 {
	return (mMouseListener.lastX - mMouseListener.xPos)
}
func GetDy() float64 {
	return (mMouseListener.lastY - mMouseListener.yPos)
}

func GetScrollX() float64 {
	return mMouseListener.scrollX
}
func GetScrollY() float64 {
	return mMouseListener.scrollY
}

func SetCamera(camera *gogl.Camera) {
	mMouseListener.camera = camera
}
func GetOrthoX() float64 {//мировые координаты
	currentX := GetX()

	currentX = (currentX/float64(mWinListener.width))*2 - 1

	tmp := mgl32.Vec4{float32(currentX), 0, 0, 1}
	inv := mMouseListener.camera.GetInverseProjection().Mul4x1(tmp)
	inv2 := mMouseListener.camera.GetInverseView().Mul4x1(inv)

	currentX = float64(inv2.X())
	
	return currentX
}
func GetOrthoY() float64 {
	currentY := GetY()

	currentY = (currentY/float64(mWinListener.height))*2 - 1

	tmp := mgl32.Vec4{0, float32(currentY), 0, 1}
	inv := mMouseListener.camera.GetInverseProjection().Mul4x1(tmp)
	inv2 := mMouseListener.camera.GetInverseView().Mul4x1(inv)

	currentY = float64(inv2.Y())
	
	return currentY
}
func IsDragging() bool {
	return mMouseListener.isDragging
}

func MouseButtonDown(button glfw.MouseButton) bool {
	if int(button) < len(mMouseListener.mouseBtnPressed) {
		return mMouseListener.mouseBtnPressed[button]
	} else {
		return false
	}
}
