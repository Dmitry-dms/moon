package gogl

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Line2D struct {
	from     mgl.Vec2
	to       mgl.Vec2
	color    mgl.Vec3
	lifetime int //frames
}

func NewLine2D(from, to mgl.Vec2, color mgl.Vec3, lifetime int) *Line2D {
	l := Line2D{
		from:     from,
		to:       to,
		color:    color,
		lifetime: lifetime,
	}
	return &l
}

//for physics
func NewLine2Df(from, to mgl.Vec2) *Line2D {
	l := Line2D{
		from: from,
		to:   to,
	}
	return &l
}
func (l *Line2D) BeginFrame() int {
	l.lifetime--
	return l.lifetime
}

func (l *Line2D) From() mgl.Vec2 {
	return l.from
}
func (l *Line2D) To() mgl.Vec2 {
	return l.to
}
func (l *Line2D) Color() mgl.Vec3 {
	return l.color
}

func (l *Line2D) LenSqr() float32 {
	return l.to.Sub(l.from).LenSqr()
}
