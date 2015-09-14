package vibrant

import "fmt"
import "strings"

type Swatch struct {
	Color      int
	Population int
	Name       string
}

func (sw *Swatch) String() string {
	return fmt.Sprintf(".%s { background-color: %s; color: %s;}\n", strings.ToLower(sw.Name), sw.RGBHex(), sw.BodyTextColor())
}

func rgbHex(color int) string {
	r, g, b := unpackColor(color)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func (sw *Swatch) RGBHex() string {
    return rgbHex(sw.Color)
}

func (sw *Swatch) TitleTextColor() string {
	return rgbHex(TextColor(sw.Color, MIN_CONTRAST_TITLE_TEXT))
}

func (sw *Swatch) BodyTextColor() string {
	return rgbHex(TextColor(sw.Color, MIN_CONTRAST_BODY_TEXT))
}
