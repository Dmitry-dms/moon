package editor

import (

	"github.com/inkyblackness/imgui-go/v4"
)

var opend bool = true

func Imgui(aspRatio float32, frameBufferTexId uint32) {
	imgui.BeginV("Game viewport", &opend, imgui.WindowFlagsNoScrollWithMouse|imgui.WindowFlagsNoScrollbar)
	winSize := getLargestSizeViewport(aspRatio)
	winPos := getCenteredPositionViewport(winSize)

	imgui.SetCursorPos(winPos)
	// imgui.Image(imgui.TextureID(frameBufferTexId), winSize)
	imgui.ImageV(imgui.TextureID(frameBufferTexId), winSize, imgui.Vec2{0,1}, imgui.Vec2{1,0}, imgui.Vec4{1,1,1,1}, imgui.Vec4{})
	imgui.End()
}

func getLargestSizeViewport(aspRatio float32) imgui.Vec2 {
	winSize := imgui.ContentRegionAvail()
	winSize.X -= imgui.ScrollX()
	winSize.Y -= imgui.ScrollY()

	aspectWidth := winSize.X
	aspectHeight := aspectWidth / aspRatio
	if aspectHeight > winSize.Y {
		// We must switch to pillarbox mode
		aspectHeight = winSize.Y
		aspectWidth = aspectHeight * aspRatio
	}
	return imgui.Vec2{aspectWidth, aspectHeight}
}
func getCenteredPositionViewport(aspectSize imgui.Vec2) imgui.Vec2 {
	winSize := imgui.ContentRegionAvail()
	winSize.X -= imgui.ScrollX()
	winSize.Y -= imgui.ScrollY()

	viewportX := (winSize.X / 2) - (aspectSize.X / 2);
	viewportY := (winSize.Y / 2) - (aspectSize.Y / 2);

	return imgui.Vec2{viewportX + imgui.CursorPosX(), viewportY+imgui.CursorPosY()}
}
