package vibrant

import "fmt"
import "strings"

type Swatch struct {
	Color      int // 24-bit int
	Population int
	Name       string // might be empty
}

// Convenience method that returns CSS (minified) in the form of
// .(swatch.Name) { background-color: (swatch.RGBHex()); color: (swatch.BodyTextColor) }
// or if that wasn't clear enough:
// .vibrant{background-color:#bada55;color:#ffffff;}
func (sw *Swatch) String() string {
	return fmt.Sprintf(".%s{background-color:%s;color:%s;}", strings.ToLower(sw.Name), sw.RGBHex(), sw.BodyTextColor())
}

func rgbHex(color int) string {
	r, g, b := unpackColor(color)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

//
func (sw *Swatch) RGB() (r, g, b int) {
	return unpackColor(sw.Color)
}

// Returns Swatch.Color as a 24-bit hex, like HTML, e.g. "#bada55"
func (sw *Swatch) RGBHex() string {
	return rgbHex(sw.Color)
}

// Returns either "#ffffff" or "#000000" based on MIN_CONTRAST_TITLE_TEXT
func (sw *Swatch) TitleTextColor() string {
	return rgbHex(TextColor(sw.Color, MIN_CONTRAST_TITLE_TEXT))
}

// Returns either "#ffffff" or "#000000" based on MIN_CONTRAST_BODY_TEXT
func (sw *Swatch) BodyTextColor() string {
	return rgbHex(TextColor(sw.Color, MIN_CONTRAST_BODY_TEXT))
}
