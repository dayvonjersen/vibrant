package vibrant

import "image/color"
import "math"

func rgb(rgba ...uint32) (r, g, b float64) {
	r = float64(rgba[0] >> 8)
	g = float64(rgba[1] >> 8)
	b = float64(rgba[2] >> 8)
	return r, g, b
}

func Hue(c color.Color) float64 {
	r, g, b := rgb(c.RGBA())

	v := math.Max(b, math.Max(r, g))
	t := math.Min(b, math.Min(r, g))

	if v == t {
		return 0
	}

	vt := v - t
	cr := (v - r) / vt
	cg := (v - g) / vt
	cb := (v - b) / vt

	var h float64
	switch {
	case r == v:
		h = cb - cg
	case g == v:
		h = 2 + cr - cb
	default:
		h = 4 + cg - cr
	}

	h /= 6
	if h < 0 {
		h++
	}
	return h
}

func Saturation(c color.Color) float64 {
	r, g, b := rgb(c.RGBA())
	v := math.Max(b, math.Max(r, g))
	t := math.Min(b, math.Min(r, g))
	if v == t {
		return 0
	}
	return (v - t) / v
}

func Brightness(c color.Color) float64 {
	r, g, b := rgb(c.RGBA())
	v := math.Max(b, math.Max(r, g))
	return v / 255
}
