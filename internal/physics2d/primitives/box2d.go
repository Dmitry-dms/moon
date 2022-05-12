package primitives

import (
	"github.com/Dmitry-dms/moon/internal/physics2d/rigidbody"
	"github.com/Dmitry-dms/moon/pkg/math"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Box2D struct {
	size     mgl.Vec2
	halfSize mgl.Vec2
	body     *rigidbody.RigidBody2D
}

func NewBox2D(min, max mgl.Vec2) *Box2D { //от нижневого левого угла, до правого верхнего
	size := max.Sub(min)

	box := Box2D{
		size:     size,
		body:     &rigidbody.RigidBody2D{},
		halfSize: size.Mul(0.5),
	}
	return &box
}

//bottom left
func (a *Box2D) GetMin() mgl.Vec2 {
	return a.body.GetPosition().Sub(a.halfSize)
}

//top rigth
func (a *Box2D) GetMax() mgl.Vec2 {
	return a.body.GetPosition().Add(a.halfSize)
}

func (a *Box2D) GetHalfSize() mgl.Vec2 {
	return a.halfSize
}

func (a *Box2D) GetVertices() []mgl.Vec2 {
	min := a.GetMin()
	max := a.GetMax()

	vert := []mgl.Vec2{
		{min.X(), min.Y()}, {min.X(), max.Y()},
		{max.X(), min.Y()}, {max.X(), max.Y()},
	}

	if a.body.GetRotation() != 0 {
		for _, v := range vert {
			math.Rotate(&v, a.body.GetRotation(), a.body.GetPosition())
		}
	}
	return vert
}

func (b *Box2D) GetRigidbody() *rigidbody.RigidBody2D {
	return b.body
}
