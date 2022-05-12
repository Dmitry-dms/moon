package primitives

import "github.com/go-gl/mathgl/mgl32"

type Collider2D struct {
	offset mgl32.Vec2
}

func DefCollider2D() *Collider2D {
	coll := Collider2D{
		offset: mgl32.Vec2{},
	}
	return &coll
}
