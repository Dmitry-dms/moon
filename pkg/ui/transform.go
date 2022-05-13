package ui

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	position mgl.Vec2 `json:"position"`
	scale    mgl.Vec2 `json:"scale"`
}


func (t *Transform) IsEqual(tr *Transform) bool {
	return t.position == tr.position && t.scale == tr.scale
}

// type tranformExported struct {
// 	Position mgl.Vec2 `json:"position"`
// 	Scale    mgl.Vec2 `json:"scale"`
// }

// func (t *Transform) MarshalJSON() ([]byte, error) {
// 	tr := tranformExported{
// 		Position: t.position,
// 		Scale:    t.scale,
// 	}
// 	return json.Marshal(tr)
// }
// func (t *Transform) UnmarshalJSON(data []byte) error {
// 	var tr tranformExported
// 	err := json.Unmarshal(data, &tr)
// 	t.position = tr.Position
// 	t.scale = tr.Scale
// 	return err
// }

func NewTransform(pos, scale mgl.Vec2) *Transform {
	return &Transform{position: pos, scale: scale}
}

func (t *Transform) Copy() *Transform {
	return &Transform{position: t.position, scale: t.scale}
}
func (t *Transform) SetPosition(pos mgl.Vec2) {
	t.position = pos
}
func (t *Transform) SetScale(scale mgl.Vec2) {
	t.scale = scale
}
func (t *Transform) GetPosition() mgl.Vec2 {
	return t.position
}
func (t *Transform) GetScale() mgl.Vec2 {
	return t.scale
}

func (t *Transform) CopyTo(to *Transform) {
	to.position = t.position
	to.scale = t.scale
}
