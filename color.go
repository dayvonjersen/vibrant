package vibrant

import colorconv "code.google.com/p/sadbox/color"
import "image/color"
import "math"

func rgb(rgba ...uint32) (r, g, b float64) {
	r = float64(rgba[0] >> 8)
	g = float64(rgba[1] >> 8)
	b = float64(rgba[2] >> 8)
	return r, g, b
}

func colorToRgb(c color.Color) (int, int, int) {
	r, g, b, _ := c.RGBA()
	return int(r >> 8), int(g >> 8), int(b >> 8)
}

func rgbToColor(r, g, b int) color.Color {
	rgba := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff}
	return rgba
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
    r,g,b := unpackColor(color)
    h, s, l = colorconv.RGBToHSL(uint8(r),uint8(g),uint8(b))
    return
}

func HslToRgb(h, s, l float64) (rgb int) {
    r, g, b := colorconv.HSLToRGB(h, s, l)
    return packColor(int(r),int(g),int(b))
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
