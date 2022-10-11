package styles

type Style struct {

	// widget space
	Margin               float32
	Padding              float32
	TopMargin, BotMargin float32
	LeftMargin           float32

	//btn
	BtnColor        [4]float32
	BtnHoveredColor [4]float32
	BtnActiveColor  [4]float32

	//text
	TextColor [4]float32
	TextSize  int

	//slider
	SliderColor    [4]float32
	SliderBtnColor [4]float32
	SliderHeight   float32
	SliderBtnWidth float32
}

var (
	DefaultStyle = Style{
		Margin:          10,
		Padding:         10,
		TopMargin:       10,
		BotMargin:       10,
		LeftMargin:      10,
		BtnColor:        [4]float32{124, 90, 156, 1},
		BtnHoveredColor: [4]float32{200, 270, 30, 1},
		BtnActiveColor:  [4]float32{255, 0, 0, 1},
		TextColor:       [4]float32{255, 255, 255, 1},
		TextSize:        20,
		SliderColor:     [4]float32{231, 240, 162, 0.8},
		SliderBtnColor:  [4]float32{0, 0, 0, 1},
		SliderHeight:    10,
		SliderBtnWidth:  20,
	}
)
