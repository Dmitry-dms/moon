package widgets

type VSpace struct {
	Height float32
}

func (s VSpace) Rectangle() [4]float32 {
	return [4]float32{}
}

func (s VSpace) GetColor() [4]float32 {
	return [4]float32{}
}
func (s VSpace) GetId() string {
	return "vert_spacing"
}
