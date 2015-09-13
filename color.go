package vibrant

import "image/color"
import "math"

func rgb(rgba ...uint32) (r, g, b float64) {
	r = float64(rgba[0] >> 8)
	g = float64(rgba[1] >> 8)
	b = float64(rgba[2] >> 8)
	return r, g, b
}

func RgbToHsl(color int) (h, s, l float64) {
	r := float64(color>>16&0xff) / 255
	g := float64(color>>8&0xff) / 255
	b := float64(color>>0&0xff) / 255
	min := math.Min(r, math.Min(g, b))
	max := math.Max(r, math.Max(g, b))
	delta := max - min

	l = (max + min) / 2

	if delta == 0 {
		h = 0
		s = 0
	} else {
		switch max {
		case r:
			h = math.Mod((g-b)/delta, 6)
		case g:
			h = ((b - r) / delta) + 2
		case b:
			h = ((r - g) / delta) + 4
		}
		s = delta / (1 - math.Abs(2*l-1))
		h = math.Mod(h*60, 360)
	}
	return h, s, l
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
