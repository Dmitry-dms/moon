package gogl

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	proj, view, invProj, invView mgl.Mat4
	Position                     mgl.Vec2
	projWidth, projHeight        float32
	ClearColor                   mgl.Vec4
	projSize                     mgl.Vec2
	zoom                         float32
}

func NewCamera(position mgl.Vec2) *Camera {
	c := Camera{
		Position: position,
		proj:     mgl.Mat4{},
		view:     mgl.Mat4{},
		invProj:  mgl.Mat4{},
		invView:  mgl.Mat4{},
		projSize: mgl.Vec2{32, 32},
		zoom:     1,
	}
	c.AdjustProjection()
	return &c
}

func (c *Camera) SetPosition(vec mgl.Vec2) {
	c.Position = c.Position.Add(vec)
}
func (c *Camera) AdjustProjection() {
	//c.proj = mgl.Ortho(0, c.projSize[0]*c.zoom, 0, c.projSize[1]*c.zoom, 0, 100)
	c.proj = mgl.Ortho(0, 32*40, 0, 32*21, 0, 100)
	c.invProj = c.proj.Inv()
}

func (c *Camera) GetViewMatrix() mgl.Mat4 {
	cameraFront := mgl.Vec3{0, 0, -1}
	cameraUp := mgl.Vec3{0, 1, 0}
	//c.view = mgl.Ident4()
	// := c.Position
	//c.view = mgl.LookAt(p[0], p[1], 20, cameraFront[0]+p[0], cameraFront[1]+p[1], cameraFront[2]+0, 0, 1, 0)
	c.view = mgl.LookAtV(mgl.Vec3{c.Position.X(), c.Position.Y(), 20},
		cameraFront.Add(mgl.Vec3{c.Position.X(), c.Position.Y(), 0}),
		cameraUp)

	c.invView = c.view.Inv()
	return c.view
}
func (c *Camera) GetProjectionMatrix() mgl.Mat4 {
	return c.proj
}
