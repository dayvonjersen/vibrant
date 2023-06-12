package vibrant

import (
	"encoding/json"
	"errors"
	"image"
	"math"
	"strings"
)

// These constants are taken directly from the Android Palette source code,
// with the exception of DEFAULT_CALCULATE_NUMBER_COLORS which was 16.
const (
	calculateBitmapMinDimension     = 100
	DEFAULT_CALCULATE_NUMBER_COLORS = 256
	TARGET_DARK_LUMA                = 0.26
	MAX_DARK_LUMA                   = 0.45
	MIN_LIGHT_LUMA                  = 0.55
	TARGET_LIGHT_LUMA               = 0.74
	MIN_NORMAL_LUMA                 = 0.3
	TARGET_NORMAL_LUMA              = 0.5
	MAX_NORMAL_LUMA                 = 0.7
	TARGET_MUTED_SATURATION         = 0.3
	MAX_MUTED_SATURATION            = 0.4
	TARGET_VIBRANT_SATURATION       = 1
	MIN_VIBRANT_SATURATION          = 0.35
	WEIGHT_SATURATION               = 3
	WEIGHT_LUMA                     = 6
	WEIGHT_POPULATION               = 1
	MIN_CONTRAST_TITLE_TEXT         = 3.0
	MIN_CONTRAST_BODY_TEXT          = 4.5
)

type Palette struct {
	// Contains the quantized palette for a given source image
	swatches          []*Swatch
	highestPopulation int
	selected          []*Swatch
}

// CSS4
//
// Export as CSS4 variables. Use in your CSS with var()
// e.g. html { background-color: var(--vibrant); }
//
// FIXME(day): variable names not final
// FIXME(day): variable names not final
// FIXME(day): variable names not final
// --vibrant:#123456;
// --vibrantText:#ffffff;
// --darkvibrant:#123456;
// --darkvibrantText:#ffffff;
// --lightvibrant:#123456;
// --lightvibrantText:#ffffff;
// --muted:#123456;
// --mutedText:#ffffff;
// --darkmuted:#123456;
// --darkmutedText:#ffffff;
// --lightmuted:#123456;
// --lightmutedText:#ffffff;
func (p *Palette) CSS4() string {
	out := ""
	for k, sw := range p.ExtractAwesome() {
		// @consistency these variable names should match JSON output
		// @consistency these variable names should match JSON output
		// @consistency these variable names should match JSON output
		out += "--" + strings.ToLower(k) + ":" + sw.Color.String() + ";"
		out += "--" + strings.ToLower(k) + "Text:" + sw.Color.TitleTextColor().String() + ";"
	}
	return out
}

// JSON
//
// err should always be nil unless something catastrophic is happening.
//
// JSON structure:
//
//	{
//		"vibrant": 0x123456,
//		"dark-vibrant": 0x123456,
//		"light-vibrant": 0x123456,
//		"muted": 0x123456,
//		"dark-muted": 0x123456,
//		"light-muted": 0x123456,
//		"vibrant-text": 0x123456,
//		"dark-vibrant-text": 0x123456,
//		"light-vibrant-text": 0x123456,
//		"muted-text": 0x123456,
//		"dark-muted-text": 0x123456,
//		"light-muted-text": 0x123456
//	}
//
// Please note that encoding/json encodes hex numbers as decimal integers
// Please note that encoding/json encodes hex numbers as decimal integers
// Please note that encoding/json encodes hex numbers as decimal integers
func (p *Palette) JSON() ([]byte, error) {
	// NOTE(day): As I type this I realize a more elegant data structure might look like
	// {
	//		"vibrant": {"color": 0x99ccff, "text": 0xffffff},
	//		"dark-vibrant": {"color": 0x00007f, "text": 0xffffff},
	//		...
	// }
	// also using PascalCase camelCase or snake_case instead of kebab-case
	// might yield a more pleasurable user experience.
	// this library was designed for CSS webdev stuff (that's why vibrant.Color.String() returns #abcdef oh. Oh. I just realized something...
	// -day Sun 11 Jun 2023 10:28:06 PM EDT
	// email me@dayvonjersen if you disagree =)
	// -day Sun 11 Jun 2023 09:57:57 PM EDT
	type vibrantJSON struct {
		Vibrant          int `json:"vibrant"`
		DarkVibrant      int `json:"dark-vibrant"`
		LightVibrant     int `json:"light-vibrant"`
		Muted            int `json:"muted"`
		DarkMuted        int `json:"dark-muted"`
		LightMuted       int `json:"light-muted"`
		VibrantText      int `json:"vibrant-text"`
		DarkVibrantText  int `json:"dark-vibrant-text"`
		LightVibrantText int `json:"light-vibrant-text"`
		MutedText        int `json:"muted-text"`
		DarkMutedText    int `json:"dark-muted-text"`
		LightMutedText   int `json:"light-muted-text"`
	}
	this := p.ExtractAwesome()
	choose := func(sw *Swatch, defaultColor int, defaultText int) (primary, text int) {
		if sw == nil {
			return defaultColor, defaultText
		}
		return int(sw.Color), int(sw.Color.TitleTextColor())
	}
	// NOTE(day): sensible fallback colors (greyscale)
	data := &vibrantJSON{}
	data.Vibrant, data.VibrantText = choose(this["Vibrant"], 0xacaaaa, 0x0)
	data.LightVibrant, data.LightVibrantText = choose(this["LightVibrant"], 0xffffff, 0x0)
	data.DarkVibrant, data.DarkVibrantText = choose(this["DarkVibrant"], 0x2b2b2b, 0xffffff)
	data.LightMuted, data.LightMutedText = choose(this["LightMuted"], 0xdad5d5, 0x0)
	data.DarkMuted, data.DarkMutedText = choose(this["DarkMuted"], 0x32312f, 0xffffff)
	data.Muted, data.MutedText = choose(this["Muted"], 0x6d6a6a, 0xffffff)
	return json.Marshal(data)
}

// Calls NewPalette with DEFAULT_CALCULATE_NUMBER_COLORS as a default value for numColors.
func NewPaletteFromImage(img image.Image) (Palette, error) {
	return NewPalette(img, DEFAULT_CALCULATE_NUMBER_COLORS)
}

func NewPalette(img image.Image, numColors int) (Palette, error) {
	// The original comments in the Android source suggest using a number
	// between 12 and 32, however this almost always results in too few colors
	// and an incomplete or unsatisfactory (read: inaccurate) result set.
	//
	// For best results, I've found a numColors between 256 and 2048 to be
	// satisfactory. There is a minor (almost negligible) performance hit for
	// high numColors values when calling ExtractAwesome(), however.
	//
	// A numColors above the number of validColors found in the colorHistogram
	// will skip the quantization step outright.
	//
	// See also source code for colorCutQuantizer, vbox, and colorHistogram

	b := newBitmap(img)
	var p Palette
	if numColors < 1 {
		return p, errors.New("numColors must be 1 or greater")
	}
	minDim := math.Min(float64(b.Width), float64(b.Height))
	if minDim > calculateBitmapMinDimension {
		scaleRatio := calculateBitmapMinDimension / minDim
		b = newScaledBitmap(b.Source, scaleRatio)
	}
	ccq := newColorCutQuantizer(*b, numColors)
	swatches := ccq.QuantizedColors
	p.swatches = swatches
	var population float64 = 0
	for _, sw := range swatches {
		population = math.Max(population, float64(sw.Population))
	}
	p.highestPopulation = int(population)
	return p, nil
}

// Possible map keys are:
//
//	Vibrant
//	LightVibrant
//	DarkVibrant
//	Muted
//	LightMuted
//	DarkMuted
//
// Some or all of these keys might not be set depending on
// the source image and numColors parameters used when creating the Palette. YMMV
func (p *Palette) ExtractAwesome() map[string]*Swatch {
	profiles := map[string][]float64{
		"Vibrant":      {TARGET_NORMAL_LUMA, MIN_NORMAL_LUMA, MAX_NORMAL_LUMA, TARGET_VIBRANT_SATURATION, MIN_VIBRANT_SATURATION, 1},
		"LightVibrant": {TARGET_LIGHT_LUMA, MIN_LIGHT_LUMA, 1, TARGET_VIBRANT_SATURATION, MIN_VIBRANT_SATURATION, 1},
		"DarkVibrant":  {TARGET_DARK_LUMA, 0, MAX_DARK_LUMA, TARGET_VIBRANT_SATURATION, MIN_VIBRANT_SATURATION, 1},
		"Muted":        {TARGET_NORMAL_LUMA, MIN_NORMAL_LUMA, MAX_NORMAL_LUMA, TARGET_MUTED_SATURATION, 0, MAX_MUTED_SATURATION},
		"LightMuted":   {TARGET_LIGHT_LUMA, MIN_LIGHT_LUMA, 1, TARGET_MUTED_SATURATION, 0, MAX_MUTED_SATURATION},
		"DarkMuted":    {TARGET_DARK_LUMA, 0, MAX_DARK_LUMA, TARGET_MUTED_SATURATION, 0, MAX_MUTED_SATURATION},
	}
	res := make(map[string]*Swatch)
	for name, args := range profiles {
		sw := p.findColor(args...)
		if sw != nil {
			sw.Name = name
			res[name] = sw
		}
	}

	if _, vib := res["Vibrant"]; !vib {
		if darkvib, ok := res["DarkVibrant"]; ok {
			h, s, l := rgbToHsl(int(darkvib.Color))
			l = TARGET_NORMAL_LUMA
			res["Vibrant"] = &Swatch{Name: "Vibrant", Color: Color(hslToRgb(h, s, l))}
		}
	}
	if _, darkvib := res["DarkVibrant"]; !darkvib {
		if vib, ok := res["Vibrant"]; ok {
			h, s, l := rgbToHsl(int(vib.Color))
			l = TARGET_DARK_LUMA
			res["DarkVibrant"] = &Swatch{Name: "DarkVibrant", Color: Color(hslToRgb(h, s, l))}
		}
	}
	return res
}

func (p *Palette) isAlreadySelected(swatch *Swatch) bool {
	for _, sw := range p.selected {
		if swatch == sw {
			return true
		}
	}
	return false
}

// don't worry about it
func (p *Palette) findColor(params ...float64) *Swatch {
	if len(params) != 6 {
		panic("not enough arguments in call to p.FindColor")
	}

	targetLuma := params[0]
	minLuma := params[1]
	maxLuma := params[2]
	targetSaturation := params[3]
	minSaturation := params[4]
	maxSaturation := params[5]

	return p.FindColor(targetLuma, minLuma, maxLuma, targetSaturation, minSaturation, maxSaturation)
}

// Finds a Swatch which best matches the specified parameters.
//
// See also package constants.
func (p *Palette) FindColor(targetLuma, minLuma, maxLuma, targetSaturation, minSaturation, maxSaturation float64) *Swatch {
	var swatch *Swatch
	var maxValue float64 = 0
	population := 0
	for _, sw := range p.swatches {
		_, sat, luma := rgbToHsl(int(sw.Color))
		if sat >= minSaturation && sat <= maxSaturation && luma >= minLuma && luma <= maxLuma && !p.isAlreadySelected(sw) {
			population += sw.Population
			value := weightedMean(invertDiff(sat, targetSaturation), WEIGHT_SATURATION, invertDiff(luma, targetLuma), WEIGHT_LUMA, float64(sw.Population)/float64(p.highestPopulation), WEIGHT_POPULATION)
			if swatch == nil || value > maxValue {
				swatch = sw
				maxValue = value
			}
		}
	}
	if swatch != nil {
		swatch.Population = population
		p.selected = append(p.selected, swatch)
	}
	return swatch
}

// Returns a value in the range 0-1.
// 1 is returned when value equals target and decreases
// as the absolute difference between value and target increases.
func invertDiff(value, target float64) float64 {
	return 1 - math.Abs(value-target)
}

func weightedMean(values ...float64) float64 {
	var sum float64 = 0
	var sumWeight float64 = 0
	for i := 0; i < len(values); i += 2 {
		value := values[i]
		weight := values[i+1]
		sum += value * weight
		sumWeight += weight
	}
	return sum / sumWeight
}
