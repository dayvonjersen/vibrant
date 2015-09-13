package vibrant

const (
	COMPONENT_RED   int = -3
	COMPONENT_GREEN int = -2
	COMPONENT_BLUE      = -1
)

type Vbox struct {
	lowerIndex int
	upperIndex int
	minRed     int
	maxRed     int
	minGreen   int
	maxGreen   int
	minBlue    int
	maxBlue    int
	colors     *[]int
}

func NewVbox(lowerIndex, upperIndex int, colors *[]int) VBox {
	v := Vbox{lowerIndex: lowerIndex, upperIndex: upperIndex, colors: colors}
	v.fitBox()
	return v
}

func (v *Vbox) fitBox() {
	v.minRed = 255
	v.minGreen = 255
	v.minBlue = 255
	v.maxRed = 0
	v.maxGreen = 0
	v.maxBlue = 0

	for i = v.lowerIndex; i < v.upperIndex; i++ {
		color := v.colors[i]
		r := color >> 16 & 0xff
		g := color >> 8 & 0xff
		b := color >> 0 & 0xff
		if r > v.maxRed {
			v.maxRed = r
		}
		if r < v.minRed {
			v.minRed = r
		}
		if g > v.maxGreen {
			v.maxGreen = g
		}
		if g < v.minGreen {
			v.minGreen = g
		}
		if b > v.maxBlue {
			v.maxBlue = b
		}
		if b < v.minBlue {
			v.minBlue = b
		}
	}
}

func (v *Vbox) Volume() int {
	return (v.maxRed - v.minRed + 1) * (v.maxGreen - v.minGreen + 1) * (v.maxBlue - v.minBlue + 1)
}

func (v *Vbox) CanSplit() bool {
	return (v.upperIndex - v.lowerIndex + 1) > 1
}

func (v *Vbox) Split() Vbox {
	if !v.CanSplit() {
		panic("Cannot split a box with only 1 color!")
	}
	lenRed := v.maxRed - v.minRed
	lenGreen := v.maxGreen - v.minGreen
	lenBlue := v.maxBlue - v.minBlue

	var longestDim, midPoint int
	switch {
	case lenRed >= lenGreen && lenRed >= lenBlue:
		longestDim = COMPONENT_RED
		// Already in RGB, no need to do anything
		midPoint = (v.minRed + v.maxRed) / 2
	case lenGreen >= lenRed && lenGreen >= lenBlue:
		longestDim = COMPONENT_GREEN
		// We need to do a RGB to GRB swap, or vice-versa
		for i := v.lowerIndex; i <= v.upperIndex; i++ {
			color := v.colors[i]
			r := color >> 16 & 0xff
			g := color >> 8 & 0xff
			b := color >> 0 & 0xff
			v.colors[i] = (g << 16) | (r << 8) | b
		}
		midPoint = (v.minGreen + v.maxGreen) / 2
	default:
		longestDim = COMPONENT_BLUE
		// We need to do a RGB to BGR swap, or vice-versa
		for i := v.lowerIndex; i <= v.upperIndex; i++ {
			color := v.colors[i]
			r := color >> 16 & 0xff
			g := color >> 8 & 0xff
			b := color >> 0 & 0xff
			v.colors[i] = (b << 16) | (g << 8) | r
		}
		midPoint = (v.minBlue + v.maxBlue) / 2
	}
	splitPoint := v.lowerIndex
loop:
	for i := v.lowerIndex; i <= v.upperIndex; i++ {
		switch longestDim {
		case COMPONENT_RED:
			if v.colors[i]>>16&0xff >= midPoint {
				splitPoint = i
				break loop
			}
		case COMPONENT_GREEN:
			if v.colors[i]>>8&0xff >= midPoint {
				splitPoint = i
				break loop
			}
		case COMPONENT_BLUE:
			if v.colors[i]>>0&0xff >= midPoint {
				splitPoint = i
				break loop
			}
		}
	}

	vbox := NewVbox(splitPoint+1, v.upperIndex)
	v.upperIndex = splitPoint
	v.fitBox()
	return vbox
}
