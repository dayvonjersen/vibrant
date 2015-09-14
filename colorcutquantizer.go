package vibrant

import "container/heap"
//import "fmt"

const (
	BLACK_MAX_LIGHTNESS float64 = 0.05
	WHITE_MIN_LIGHTNESS float64 = 0.95
)

type ColorCutQuantizer struct {
	Colors           []int
	ColorPopulations map[int]int
	QuantizedColors  []*Swatch
}

func shouldIgnoreColor(color int) bool {
	h, s, l := RgbToHsl(color)
	return l <= BLACK_MAX_LIGHTNESS || l >= WHITE_MIN_LIGHTNESS || (h >= 0.0278 && h <= 0.1028 && s <= 0.82)
}
func shouldIgnoreColorSwatch(sw *Swatch) bool {
	return shouldIgnoreColor(sw.Color)
}

func NewColorCutQuantizer(bitmap Bitmap, maxColors int) *ColorCutQuantizer {
	pixels := bitmap.Pixels()
	histo := NewColorHistogram(pixels)
	colorPopulations := make(map[int]int, histo.NumberColors)
	//fmt.Printf("histogram colors: %d\n", histo.NumberColors)
	for i, c := range histo.Colors {
		colorPopulations[c] = histo.ColorCounts[i]
	}
	validColors := make([]int, 0)
	i := 0
	for _, c := range histo.Colors {
		if !shouldIgnoreColor(c) {
			validColors = append(validColors, c)
			i++
		}
	}
	//fmt.Printf("valid colors: %d\n", len(validColors))
	validCount := len(validColors)
	// XXX complete arbitrary and temporary XXX
/*	switch {
	case validCount < 5000:
		maxColors = 1024
	case validCount >= 5000 && validCount < 10000:
		maxColors = 2048
	case validCount >= 10000 && validCount < 20000:
		maxColors = 4096
	}*/
	ccq := &ColorCutQuantizer{Colors: validColors, ColorPopulations: colorPopulations}
	if validCount <= maxColors {
		for _, c := range validColors {
			ccq.QuantizedColors = append(ccq.QuantizedColors, &Swatch{Color: c, Population: colorPopulations[c]})
		}
	} else {
		ccq.quantizePixels(validCount-1, maxColors)
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
