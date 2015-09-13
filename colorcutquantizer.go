package vibrant

import "container/heap"

const (
	BLACK_MAX_LIGHTNESS float64 = 0.05
	WHITE_MIN_LIGHTNESS float64 = 0.95
)

type ColorCutQuantizer struct {
	//	tempHSL          [3]float64
	Colors           map[int]int
	ColorPopulations map[int]int
	QuantizedColors  []*Swatch
}

func shouldIgnoreColor(color int) bool {
	h, s, l := RgbToHsl(color)
	return l <= BLACK_MAX_LIGHTNESS || l >= WHITE_MIN_LIGHTNESS || (h >= 10 && h <= 37 && s <= 0.82)
}
func shouldIgnoreColorSwatch(sw *Swatch) bool {
	return shouldIgnoreColor(sw.Color)
}

func NewColorCutQuantizer(bitmap Bitmap, maxColors int) *ColorCutQuantizer {
	pixels := bitmap.Pixels()
	histo := NewColorHistogram(pixels)
	colorPopulations := make(map[int]int, histo.NumberColors)
	for i, c := range histo.Colors {
		colorPopulations[c] = histo.ColorCounts[i]
	}
	validColors := make(map[int]int, 0)
	i := 0
	for _, c := range histo.Colors {
		if !shouldIgnoreColor(c) {
			validColors[i] = c
			i++
		}
	}
	ccq := &ColorCutQuantizer{Colors: validColors, ColorPopulations: colorPopulations}
	if len(validColors) <= maxColors {
		for _, c := range validColors {
			ccq.QuantizedColors = append(ccq.QuantizedColors, &Swatch{Color: c, Population: colorPopulations[c]})
		}
	} else {
		ccq.quantizePixels(len(validColors)-1, maxColors)
	}
	return ccq
}

func (ccq *ColorCutQuantizer) quantizePixels(maxColorIndex, maxColors int) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, NewVbox(0, maxColorIndex, ccq.Colors, ccq.ColorPopulations))
	for pq.Len() < maxColors {
		v := heap.Pop(&pq).(*Vbox)
		if v.CanSplit() {
			heap.Push(&pq, v.Split())
			heap.Push(&pq, v)
		} else {
			break
		}
	}
	for pq.Len() > 0 {
		v := heap.Pop(&pq).(*Vbox)
		swatch := v.AverageColor()
		if !shouldIgnoreColorSwatch(swatch) {
			ccq.QuantizedColors = append(ccq.QuantizedColors, swatch)
        }
	}
}
