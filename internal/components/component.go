package components


// func (Component) Update(dt float32) {

// }

var ID_COUNTER int = 0

func Init(maxId int) {
	ID_COUNTER = maxId
}

type Component interface {
	Update(dt float32)
	Start()
}