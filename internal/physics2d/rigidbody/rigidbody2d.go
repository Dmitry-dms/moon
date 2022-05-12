package rigidbody

import mgl "github.com/go-gl/mathgl/mgl32"

//Components
type RigidBody2D struct {
	position mgl.Vec2
	rotation float32
}

func (r *RigidBody2D) SetPosition(pos mgl.Vec2) {
	r.position = pos
}
func (r *RigidBody2D) SetRotation(angle float32) {
	r.rotation = angle
}
func (r *RigidBody2D) GetPosition() mgl.Vec2 {
	return r.position
}
func (r *RigidBody2D) GetRotation() float32 {
	return r.rotation
}