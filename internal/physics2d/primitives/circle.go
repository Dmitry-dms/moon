package primitives

import (
	"github.com/Dmitry-dms/moon/internal/physics2d/rigidbody"
	"github.com/go-gl/mathgl/mgl32"
)

type Circle struct {
	radius float32
	body   *rigidbody.RigidBody2D
}

func DefCircle() *Circle {
	return &Circle{
		radius: 1,
		body: &rigidbody.RigidBody2D{},
	}
}

func (c *Circle) GetRadius() float32 {
	return c.radius
}
func (c *Circle) GetCenter() mgl32.Vec2 {
	return c.body.GetPosition()
}
func (c *Circle) SetRadius(rad float32) {
	c.radius = rad
}