package ui

type UiObject struct {
	Name      string
	Transform *Transform
	Spr       *SpriteRenderer

	ZIndex    int
	uid       int
	isMovable bool
}

type Renderable interface {
	Update(dt float32)
	Spr() *SpriteRenderer
	Transform() *Transform
}

type Rectangle struct {
	X, Y, W, H float32
}
