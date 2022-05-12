package primitives

import (
	"testing"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	mgl "github.com/go-gl/mathgl/mgl32"
)


func TestPointOnLine(t *testing.T) {
	testCases := []struct {
		desc  string
		line  *gogl.Line2D
		point mgl.Vec2
		res   bool
	}{
		{
			desc:  "Shoud return true",
			line:  gogl.NewLine2Df(mgl.Vec2{0, 0}, mgl.Vec2{12, 4}),
			point: mgl.Vec2{0, 0},
			res: true,
		},
		{
			desc:  "Shoud return true",
			line:  gogl.NewLine2Df(mgl.Vec2{0, 0}, mgl.Vec2{0, 10}),
			point: mgl.Vec2{0, 5},
			res: true,
		},
		{
			desc:  "Shoud return true",
			line:  gogl.NewLine2Df(mgl.Vec2{0, 0}, mgl.Vec2{12, 4}),
			point: mgl.Vec2{12, 4},
			res: true,
		},
		{
			desc:  "Shoud return false",
			line:  gogl.NewLine2Df(mgl.Vec2{5, 10}, mgl.Vec2{12, 4}),
			point: mgl.Vec2{3, 3},
			res: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			f := PointOnLine(tC.point, *tC.line)
			if f != tC.res {
				t.Error("Wrong")
			}
		})
	}
}
