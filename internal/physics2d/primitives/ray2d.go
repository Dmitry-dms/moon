package primitives

import mgl "github.com/go-gl/mathgl/mgl32"

type Ray2D struct {
	origin mgl.Vec2
	direction mgl.Vec2
}

func NewRay2D(origin,direction mgl.Vec2) *Ray2D {
	r := Ray2D{
		origin: origin,
		direction: direction,
	}
	r.direction.Normalize()
	return &r
}

func (r *Ray2D) GetDirection() mgl.Vec2 {
	return r.direction
}
func (r *Ray2D) GetOrigin() mgl.Vec2 {
	return r.origin
}