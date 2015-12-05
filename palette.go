package vibrant

import "errors"
import "math"
import "image"

// these constants are taken directly from the Android Palette source code
const (
	CALCULATE_BITMAP_MIN_DIMENSION float64 = 100
	//DEFAULT_CALCULATE_NUMBER_COLORS int     = 16
	DEFAULT_CALCULATE_NUMBER_COLORS int     = 256
	TARGET_DARK_LUMA                float64 = 0.26
	MAX_DARK_LUMA                   float64 = 0.45
	MIN_LIGHT_LUMA                  float64 = 0.55
	TARGET_LIGHT_LUMA               float64 = 0.74
	MIN_NORMAL_LUMA                 float64 = 0.3
	TARGET_NORMAL_LUMA              float64 = 0.5
	MAX_NORMAL_LUMA                 float64 = 0.7
	TARGET_MUTED_SATURATION         float64 = 0.3
	MAX_MUTED_SATURATION            float64 = 0.4
	TARGET_VIBRANT_SATURATION       float64 = 1
	MIN_VIBRANT_SATURATION          float64 = 0.35
	WEIGHT_SATURATION               float64 = 3
	WEIGHT_LUMA                     float64 = 6
	WEIGHT_POPULATION               float64 = 1
	MIN_CONTRAST_TITLE_TEXT         float64 = 3.0
	MIN_CONTRAST_BODY_TEXT          float64 = 4.5
)

// Contains the quantized palette for a given source image
type Palette struct {
	Swatches          []*Swatch
	HighestPopulation int
	selected          []*Swatch
}

// Convenience function, uses DEFAULT_CALCULATE_NUMBER_COLORS.
//
// To specify a different number of colors to quantize to, use NewPalette
func NewPaletteFromImage(img image.Image) (Palette, error) {
	bitmap := NewBitmap(img)
	return NewPalette(bitmap, DEFAULT_CALCULATE_NUMBER_COLORS)
}

// The original comments in the Android source suggest using a number
// between 12 and 32, however this almost always results in too few colors
// and an incomplete or unsatisfactory (read: inaccurate) result set.
//
// For best results, I've found a numColors between 256 and 2048 to be
// satisfactory. There is a minor (almost negligible) performance hit for
// high numColors values when calling ExtractAwesome(), however.
//
// A numColors above the number of validColors found in the ColorHistogram
// will skip the quantization step outright.
//
// See also source code for ColorCutQuantizer, Vbox, and ColorHistogram
func NewPalette(b *Bitmap, numColors int) (Palette, error) {
	var p Palette
	if numColors < 1 {
		return p, errors.New("numColors must be 1 or greater")
	}
	minDim := math.Min(float64(b.Width), float64(b.Height))
	if minDim > CALCULATE_BITMAP_MIN_DIMENSION {
		scaleRatio := CALCULATE_BITMAP_MIN_DIMENSION / minDim
		b = NewScaledBitmap(b.Source, scaleRatio)
	}
	ccq := NewColorCutQuantizer(*b, numColors)
	swatches := ccq.QuantizedColors
	p.Swatches = swatches
	var population float64 = 0
	for _, sw := range swatches {
		population = math.Max(population, float64(sw.Population))
	}
	p.HighestPopulation = int(population)
	return p, nil
}

//
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
			h, s, l := RgbToHsl(darkvib.Color)
			l = TARGET_NORMAL_LUMA
			res["Vibrant"] = &Swatch{Color: HslToRgb(h, s, l)}
		}
	}
	if _, darkvib := res["DarkVibrant"]; !darkvib {
		if vib, ok := res["Vibrant"]; ok {
			h, s, l := RgbToHsl(vib.Color)
			l = TARGET_DARK_LUMA
			res["DarkVibrant"] = &Swatch{Color: HslToRgb(h, s, l)}
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

// Finds a Swatch which best matches the specified parameters, taking into
// consideration also the population of colors in the source image.
//
// The returned Swatch.Population value is the sum of the populations it
// represents.
func (p *Palette) FindColor(targetLuma, minLuma, maxLuma, targetSaturation, minSaturation, maxSaturation float64) *Swatch {
	var swatch *Swatch
	var maxValue float64 = 0
	population := 0
	for _, sw := range p.Swatches {
		_, sat, luma := RgbToHsl(sw.Color)
		if sat >= minSaturation && sat <= maxSaturation && luma >= minLuma && luma <= maxLuma && !p.isAlreadySelected(sw) {
			population += sw.Population
			value := weightedMean(invertDiff(sat, targetSaturation), WEIGHT_SATURATION, invertDiff(luma, targetLuma), WEIGHT_LUMA, float64(sw.Population)/float64(p.HighestPopulation), WEIGHT_POPULATION)
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
