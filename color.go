package vibrant

//import "image/color"
import "math"

func rgb(rgba ...uint32) (r, g, b float64) {
	r = float64(rgba[0] >> 8)
	g = float64(rgba[1] >> 8)
	b = float64(rgba[2] >> 8)
	return r, g, b
}

func unpackColor(color int) (r, g, b int) {
	r = color >> 16 & 0xff
	g = color >> 8 & 0xff
	b = color >> 0 & 0xff
	return r, g, b
}

func unpackColorFloat(color int) (r, g, b float64) {
	ir, ig, ib := unpackColor(color)
	r = float64(ir)
	g = float64(ig)
	b = float64(ib)
	return r, g, b
}

func packColor(r, g, b int) int {
	return (r << 16) | (g << 8) | b
}

func RgbToHsl(color int) (h, s, l float64) {
	r, g, b := unpackColorFloat(color)
	r /= 255.0
	g /= 255.0
	b /= 255.0
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

func huetocomponent(v1, v2, h float64) float64 {
	if 6*h < 1 {
		return v1 + (v2-v1)*6*h
	}
	if 2*h < 1 {
		return v2
	}
	if 3*h < 2 {
		return v1 + (v2-v1)*((2.0/3.0)-h)*6
	}
	return v1
}

func HslToRgb(h, s, l float64) (rgb int) {
	var r, g, b int
	if s == 0 {
		r = int(l * 255)
		g = r
		b = r
	} else {
		var v1, v2 float64
		if l < 0.5 {
			v2 = l * (1 + s)
		} else {
			v2 = (l + s) - (s * l)
		}
		v1 = 2*l - v2
		if h < 0 {
			h += 1
		}
		if h > 1 {
			h -= 1
		}
		r = int(255.0 * huetocomponent(v1, v2, h+(1.0/3.0)))
		g = int(255.0 * huetocomponent(v1, v2, h))
		b = int(255.0 * huetocomponent(v1, v2, h-(1.0/3.0)))
	}
	return packColor(r, g, b)
}

func TextColor(bgColor int, contrast float64) int {
	if Contrast(0xffffff, bgColor) >= contrast {
		return 0xffffff
	}
	return 0
}

func Contrast(fg, bg int) float64 {
	lum1 := Luminance(unpackColorFloat(fg))
	lum2 := Luminance(unpackColorFloat(bg))
	return math.Max(lum1, lum2) / math.Min(lum1, lum2)
}

func Luminance(red, green, blue float64) float64 {
	red /= 255.0
	if red < 0.03928 {
		red /= 12.92
	} else {
		red = math.Pow((red+0.055)/1.055, 2.4)
	}
	green /= 255.0
	if green < 0.03928 {
		green /= 12.92
	} else {
		green = math.Pow((green+0.055)/1.055, 2.4)
	}
	blue /= 255.0
	if blue < 0.03928 {
		blue /= 12.92
	} else {
		blue = math.Pow((blue+0.055)/1.055, 2.4)
	}
	return (0.2126 * red) + (0.7152 * green) + (0.0722 * blue)
}

/*
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
}*/
