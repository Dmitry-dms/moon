package renderers

import (

	"math"

	"github.com/Dmitry-dms/moon/internal/utils"
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/mathgl/mgl32"
)

// type GridLines struct {
// 	cameraPos mgl32.Vec2
// 	projectionSize mgl32.Vec2
// }

func UpdateGridLines(camera *gogl.Camera) {
	cameraPos := camera.Position
	projSize := camera.GetProjectionSize()

	firstX := ((cameraPos.X() / utils.GRID_WIDTH) - 1) * utils.GRID_WIDTH
	firstY := ((cameraPos.Y() / utils.GRID_HEIGHT) - 1) * utils.GRID_HEIGHT

	numVtLines := (projSize.X() / utils.GRID_WIDTH) + 2
	numHzLines := (projSize.Y() / utils.GRID_HEIGHT) + 2

	height := projSize.Y() + utils.GRID_HEIGHT*2
	width := projSize.X() + utils.GRID_WIDTH*2

	maxLines := int(math.Max(float64(numVtLines), float64(numHzLines)))
	color := mgl32.Vec3{0.2, 0.2, 0.2}
	for i := 0; i < maxLines; i++ {
		x := firstX + utils.GRID_WIDTH*float32(i)
		y := firstY + utils.GRID_HEIGHT*float32(i)



		if float32(i) < numVtLines {

			DebugDraw.AddLine2d(mgl32.Vec2{x, firstY}, mgl32.Vec2{x, firstY + height}, color, 1)
		}

		if float32(i) < numHzLines {
			DebugDraw.AddLine2d(mgl32.Vec2{firstX, y}, mgl32.Vec2{firstX + width, y}, color, 1)
		}
	}
}
