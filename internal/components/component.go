package components


// func (Component) Update(dt float32) {

// }

type Component interface {
	Update(dt float32)
	Start()
}