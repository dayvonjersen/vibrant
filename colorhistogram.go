package vibrant

import "image/color"
import "sort"

type ColorHistogram struct {
	Colors       []int
	ColorCounts  []int
	NumberColors int
}

func colortoint(c color.Color) int {
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	return int((r << 4) | (g << 2) | b)
}

func NewColorHistogram(colorPixels []color.Color) *ColorHistogram {
	pixels := make([]int, len(colorPixels))
	for _, px := range colorPixels {
		pixels = append(pixels, colortoint(px))
	}
	sort.Ints(pixels)
	numColors := countDistinctColors(pixels)
	colors := make([]int, numColors)
	colorCounts := make([]int, numColors)

	if numColors > 0 {
		curIndex := 0
		curColor := pixels[0]
		colors[0] = curColor
		colorCounts[0] = 1

		for _, px := range pixels {
			if px == curColor {
				colorCounts[curIndex]++
			} else {
				curColor = px
				curIndex++
				colors[curIndex] = curColor
				colorCounts[curIndex] = 1
			}
		}
	}

	return &ColorHistogram{colors, colorCounts, numColors}
}

func countDistinctColors(pixels []int) int {
	if len(pixels) < 2 {
		return len(pixels)
	}
	count := 1
	current := pixels[0]
	for _, px := range pixels {
		if px != current {
			current = px
			count++
		}
	}
	return count
}
