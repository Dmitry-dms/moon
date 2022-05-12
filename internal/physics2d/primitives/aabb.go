package primitives

import (
	"github.com/Dmitry-dms/moon/internal/physics2d/rigidbody"
	mgl "github.com/go-gl/mathgl/mgl32"
)

//axis aligned bounding box (NOT ROTATED)
type AABB struct {
	// center mgl.Vec2
	size mgl.Vec2
	halfSize mgl.Vec2
	body *rigidbody.RigidBody2D
}

// func DefAABB() *AABB {
// 	aabb:= AABB{
// 		center: mgl.Vec2{},
// 		size: mgl.Vec2{},
// 	}
// 	return &aabb
// }

func NewAABB(min, max mgl.Vec2) *AABB { //от нижневого левого угла, до правого верхнего
	size := max.Sub(min)
	
	aabb := AABB{
		size: size,
		body: &rigidbody.RigidBody2D{},
		halfSize: size.Mul(0.5),
	}
	return &aabb
}

//bottom left
func (a *AABB) GetMin() mgl.Vec2 {
	return a.body.GetPosition().Sub(a.halfSize)
}
//top rigth
func (a *AABB) GetMax() mgl.Vec2 {
	return a.body.GetPosition().Add(a.halfSize)
}
