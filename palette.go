package vibrant

import "errors"
import "math"

const (
	CALCULATE_BITMAP_MIN_DIMENSION  float64 = 100
	DEFAULT_CALCULATE_NUMBER_COLORS int     = 16
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

type Palette struct {
	Swatches           []*Swatch
	HighestPopulation  int
	VibrantSwatch      *Swatch
	DarkVibrantSwatch  *Swatch
	LightVibrantSwatch *Swatch
	MutedSwatch        *Swatch
	DarkMutedSwatch    *Swatch
	LightMutedSwatch   *Swatch
}

func NewPalette(b *Bitmap) (Palette, error) {
	return Generate(b, DEFAULT_CALCULATE_NUMBER_COLORS)
}

func Generate(b *Bitmap, numColors int) (Palette, error) {
	var p Palette
	if numColors < 1 {
		return p, errors.New("numColors mustbe 1 or greater")
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
	p.VibrantSwatch = p.findColor(TARGET_NORMAL_LUMA, MIN_NORMAL_LUMA, MAX_NORMAL_LUMA, TARGET_VIBRANT_SATURATION, MIN_VIBRANT_SATURATION, 1)
	p.LightVibrantSwatch = p.findColor(TARGET_LIGHT_LUMA, MIN_LIGHT_LUMA, 1, TARGET_VIBRANT_SATURATION, MIN_VIBRANT_SATURATION, 1)
	p.DarkVibrantSwatch = p.findColor(TARGET_DARK_LUMA, 0, MAX_DARK_LUMA, TARGET_VIBRANT_SATURATION, MIN_VIBRANT_SATURATION, 1)
	p.MutedSwatch = p.findColor(TARGET_NORMAL_LUMA, MIN_NORMAL_LUMA, MAX_NORMAL_LUMA, TARGET_MUTED_SATURATION, 0, MAX_MUTED_SATURATION)
	p.LightMutedSwatch = p.findColor(TARGET_LIGHT_LUMA, MIN_LIGHT_LUMA, 1, TARGET_MUTED_SATURATION, 0, MAX_MUTED_SATURATION)
	p.DarkMutedSwatch = p.findColor(TARGET_DARK_LUMA, 0, MAX_DARK_LUMA, TARGET_MUTED_SATURATION, 0, MAX_MUTED_SATURATION)

	if p.VibrantSwatch == nil {
		if p.DarkVibrantSwatch != nil {
			h, s, l := RgbToHsl(p.DarkVibrantSwatch.Color)
			l = TARGET_NORMAL_LUMA
			p.VibrantSwatch = &Swatch{Color: HslToRgb(h, s, l)}
		}
	}
	if p.DarkVibrantSwatch == nil {
		if p.VibrantSwatch != nil {
			h, s, l := RgbToHsl(p.VibrantSwatch.Color)
			l = TARGET_DARK_LUMA
			p.VibrantSwatch = &Swatch{Color: HslToRgb(h, s, l)}
		}
	}

	if p.VibrantSwatch != nil {
		p.VibrantSwatch.Name = "Vibrant"
	}
	if p.DarkVibrantSwatch != nil {
		p.DarkVibrantSwatch.Name = "DarkVibrant"
	}
	if p.LightVibrantSwatch != nil {
		p.LightVibrantSwatch.Name = "LightVibrant"
	}
	if p.MutedSwatch != nil {
		p.MutedSwatch.Name = "Muted"
	}
	if p.DarkMutedSwatch != nil {
		p.DarkMutedSwatch.Name = "DarkMuted"
	}
	if p.LightMutedSwatch != nil {
		p.LightMutedSwatch.Name = "LightMuted"
	}
	return p, nil
}

func (p *Palette) isAlreadySelected(swatch *Swatch) bool {
	return (p.VibrantSwatch == swatch ||
		p.DarkVibrantSwatch == swatch ||
		p.LightVibrantSwatch == swatch ||
		p.MutedSwatch == swatch ||
		p.DarkMutedSwatch == swatch ||
		p.LightMutedSwatch == swatch)
}

func (p *Palette) findColor(targetLuma, minLuma, maxLuma, targetSaturation, minSaturation, maxSaturation float64) *Swatch {
	var swatch *Swatch
	var maxValue float64 = 0
	for _, sw := range p.Swatches {
		_, sat, luma := RgbToHsl(sw.Color)
		if sat >= minSaturation && sat <= maxSaturation && luma >= minLuma && luma <= maxLuma && !p.isAlreadySelected(sw) {
			value := weightedMean(invertDiff(sat, targetSaturation), WEIGHT_SATURATION, invertDiff(luma, targetLuma), WEIGHT_LUMA, float64(sw.Population)/float64(p.HighestPopulation), WEIGHT_POPULATION)
			if swatch == nil || value > maxValue {
				swatch = sw
				maxValue = value
			}
		}
	}
	return swatch
}
