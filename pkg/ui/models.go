package ui

type Rectangle struct {
	UiObject
}

type UiObject struct {
	Name      string
	Transform *Transform
	Spr       *SpriteRenderer

	ZIndex int
	uid    int
}

type Renderable interface {
	Update(dt float32)
	Spr() *SpriteRenderer
	Transform() *Transform
}

