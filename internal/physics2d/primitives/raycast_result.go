package primitives

import mgl "github.com/go-gl/mathgl/mgl32"

type RaycastResult struct {
	point  mgl.Vec2
	normal mgl.Vec2
	t      float32
	hit    bool
}

func (r *RaycastResult) Init(point, normal mgl.Vec2, t float32, hit bool) {
	r.point = point
	r.normal = normal
	r.t = t
	r.hit = hit
}
func reset(result *RaycastResult) {
	if result != nil {
		result.point = mgl.Vec2{}
		result.normal = mgl.Vec2{}
		result.t = -1
		result.hit = false
	}
}
