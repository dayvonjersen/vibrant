package vibrant

const (
	BLACK_MAX_LIGHTNESS float64 = 0.05
	WHITE_MIN_LIGHTNESS float64 = 0.95
	COMPONENT_RED       int     = -3
	COMPONENT_GREEN     int     = -2
	COMPONENT_BLUE              = -1
)

type ColorCutQuantizer struct {
	//	tempHSL          [3]float64
	Colors           []int
	ColorPopulations map[int]int
	QuantizedColors  []*Swatch
}

// XXX stubs
type Swatch struct{}
func NewSwatch(_ ...interface{}) *Swatch { return &Swatch{} }
func shouldIgnoreColor(color int) bool { return true }

func NewColorCutQuantizer(bitmap Bitmap, maxColors int) *ColorCutQuantizer {
	pixels := bitmap.Pixels()
	histo := NewColorHistogram(pixels)
	colorPopulations := make(map[int]int, histo.NumberColors)
	for i, c := range histo.Colors {
		colorPopulations[c] = histo.ColorCounts[i]
	}
	validColors := make([]int, 0)
	for _, c := range histo.Colors {
		if !shouldIgnoreColor(c) {
			validColors = append(validColors, c)
		}
	}
	ccq := &ColorCutQuantizer{Colors: validColors, ColorPopulations: colorPopulations}
	if len(validColors) <= maxColors {
		for _, c := range validColors {
			ccq.QuantizedColors = append(ccq.QuantizedColors, NewSwatch(c, colorPopulations[c]))
		}
	} else {
		ccq.quantizePixels(len(validColors)-1, maxColors)
	}
	return ccq
}

func (ccq *ColorCutQuantizer) quantizePixels(maxColorIndex, maxColors int) {
}
