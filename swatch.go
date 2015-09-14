package vibrant

import "fmt"
import "strings"

type Swatch struct {
	Color      int
	Population int
	Name       string
}

func (sw *Swatch) String() string {
	r, g, b := unpackColor(sw.Color)
	bg := fmt.Sprintf("#%02x%02x%02x", r, g, b)
	//return fmt.Sprintf("(Swatch) Population: %d, Color: %s, Name: %s", sw.Population, bg, sw.Name)
	//bg := fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	/*	tcolor := sw.TitleTextColor()
		r, g, b = unpackColor(tcolor)
		tt := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		//tt := fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
		bcolor := sw.BodyTextColor()
		r, g, b = unpackColor(bcolor)
		bt := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		//bt := fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
    */
	return fmt.Sprintf(".%s { background-color: %s;}\n", strings.ToLower(sw.Name), bg)
}

func (sw *Swatch) TitleTextColor() int {
	return TextColor(sw.Color, MIN_CONTRAST_TITLE_TEXT)
}

func (sw *Swatch) BodyTextColor() int {
	return TextColor(sw.Color, MIN_CONTRAST_BODY_TEXT)
}
