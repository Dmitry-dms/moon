package listeners

import (
	// "fmt"
	// "fmt"
	// "fmt"
	// "fmt"

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

	gameViewPortPos, gameViewPortSize mgl32.Vec2
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
	// fmt.Printf("%.1f %.1f \n", xpos, ypos)

	mMouseListener.isDragging = mMouseListener.mouseBtnPressed[0] ||
		mMouseListener.mouseBtnPressed[1] ||
		mMouseListener.mouseBtnPressed[2]
	// if mMouseListener.isDragging {
	// 	fmt.Println("dragging")
	// }
}

//mods - это сочетания клавиш
func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		//если кнопок на мышке больше, игнорируем нажатие
		if int(button) < len(mMouseListener.mouseBtnPressed) {
			mMouseListener.mouseBtnPressed[button] = true
		}
		// fmt.Printf("Pressed %.1f %.1f \n", mMouseListener.xPos, mMouseListener.yPos)
	case glfw.Release:
		if int(button) < len(mMouseListener.mouseBtnPressed) {
			mMouseListener.mouseBtnPressed[button] = false
			mMouseListener.isDragging = false
		}
		// fmt.Printf("Released %.1f %.1f \n", mMouseListener.xPos, mMouseListener.yPos)
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
	// fmt.Println(float64(mWinListener.height) - mMouseListener.yPos)
	// return float64(mWinListener.height) - mMouseListener.yPos
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
func GetOrthoX() float64 { //мировые координаты
	currentX := GetX() - float64(mMouseListener.gameViewPortPos[0])
	// b := (currentX/float64(mWinListener.width))
	currentX = (currentX/float64(mWinListener.width))*2 - 1 // интервал от 0 до 1
	// fmt.Printf("CURRENT X = %.0f WIDTH =  %.0f DIM =  %.0f\n",
	// 	currentX, float64(mWinListener.width), b)

	tmp := mgl32.Vec4{float32(currentX), 0, 0, 1}
	inv := mMouseListener.camera.GetInverseProjection().Mul4x1(tmp)
	inv2 := mMouseListener.camera.GetInverseView().Mul4x1(inv)

	currentX = float64(inv2.X())
	// fmt.Printf("ORTHO X = %.0f GETX = %.0f WIDTH =  %d \n",
	// 	currentX, GetX(), mWinListener.width)

	// return float64(100)
	return currentX
}
func GetOrthoY() float64 {
	// currentY := GetY() - float64(mMouseListener.gameViewPortPos[1])
	currentY := float64(mWinListener.height) - GetY()
	// fmt.Printf("VIEWPORT SIZE %.1f %.1f \n", mMouseListener.gameViewPortSize[0], mMouseListener.gameViewPortSize[1])

	currentY = ((currentY/float64(mMouseListener.gameViewPortSize[1]))*2 - 1)

	tmp := mgl32.Vec4{0, float32(currentY), 0, 1}
	inv := mMouseListener.camera.GetInverseProjection().Mul4x1(tmp)
	inv2 := mMouseListener.camera.GetInverseView().Mul4x1(inv)

	currentY = float64(inv2.Y())

	return currentY
}

func SetGameViewPortPos(pos mgl32.Vec2) {
	mMouseListener.gameViewPortPos = pos
}
func SetGameViewPortSize(size mgl32.Vec2) {
	mMouseListener.gameViewPortSize = size
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

func RegionHit(x, y, w, h float32) bool {

	ys := float32(mWinListener.height)-float32(GetY())
	return float32(GetX()) >= x && ys >= y && float32(GetX()) <= x+w && ys <= y+h
}
