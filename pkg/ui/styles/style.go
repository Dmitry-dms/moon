package styles

type Style struct {

	// widget space
	Margin     float32
	Padding    float32
	TopMargin,BotMargin  float32
	LeftMargin float32

	//btn
	BtnColor        [4]float32
	BtnHoveredColor [4]float32
	BtnActiveColor  [4]float32

	//text
	TextColor [4]float32
	TextSize  int
}

var (
	DefaultStyle = Style{
		TopMargin:       10,
		BotMargin:       10,
		LeftMargin:      10,
		Padding:         10,
		BtnColor:        [4]float32{124, 90, 156, 1},
		BtnHoveredColor: [4]float32{200, 270, 30, 1},
		BtnActiveColor:  [4]float32{255, 0, 0, 1},
		TextColor:       [4]float32{255, 255, 255, 1},
		TextSize:        20,
	}
)
