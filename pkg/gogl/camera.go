package gogl

import (
	
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Proj, View, invProj, invView mgl.Mat4
	Position                     mgl.Vec2
	ProjWidth, ProjHeight        float32
	ClearColor                   mgl.Vec4
	ProjSize                     mgl.Vec2
	Zoom                         float32
}

func NewCamera(position mgl.Vec2) *Camera {
	c := Camera{
		Position: position,
		Proj:     mgl.Mat4{},
		View:     mgl.Mat4{},
		invProj:  mgl.Mat4{},
		invView:  mgl.Mat4{},
		// projSize: mgl.Vec2{32*40, 32*21},
		ProjSize: mgl.Vec2{1920, 1080},
		Zoom:     1,
	}
	c.AdjustProjection()
	return &c
}

func (c *Camera) SetPosition(vec mgl.Vec2) {
	c.Position = c.Position.Add(vec)
}
func (c *Camera) GetProjectionSize() mgl.Vec2 {
	return c.ProjSize
}
func (c *Camera) UpdateProjection(pr mgl.Vec2) {
	c.ProjSize = pr
	c.Proj = mgl.Ortho(0, pr[0], 0, pr[1], 0, 100)
	c.invProj = c.Proj.Inv()
}
func (c *Camera) AdjustProjection() {
	c.Proj = mgl.Ortho(0, c.ProjSize[0], 0, c.ProjSize[1], 0, 100)
	c.invProj = c.Proj.Inv()
}

func (c *Camera) GetViewMatrix() mgl.Mat4 {


	cameraFront := mgl.Vec3{0, 0, -1}
	cameraUp := mgl.Vec3{0, 1, 0}
	c.View = mgl.LookAtV(mgl.Vec3{c.Position.X(), c.Position.Y(), 20},
		cameraFront.Add(mgl.Vec3{c.Position.X(), c.Position.Y(), 0}),
		cameraUp)

	c.invView = c.View.Inv()

	return c.View
}
func (c *Camera) GetProjectionMatrix() mgl.Mat4 {
	return c.Proj
}

func (c *Camera) GetInverseProjection() mgl.Mat4 {
	return c.invProj
}
func (c *Camera) GetInverseView() mgl.Mat4 {
	return c.invView
}
