package vibrant

import "github.com/nfnt/resize"
import "image"
import "image/color"
import "math"

type Bitmap struct {
	Width  int
	Height int
	Source image.Image
}

func NewBitmap(input image.Image) *Bitmap {
	bounds := input.Bounds()
	return &Bitmap{bounds.Dx(), bounds.Dy(), input}
}

func NewScaledBitmap(input image.Image, ratio float64) *Bitmap {
	bounds := input.Bounds()
	w := math.Ceil(float64(bounds.Dx()) * ratio)
	h := math.Ceil(float64(bounds.Dy()) * ratio)
	return &Bitmap{int(w), int(h), resize.Resize(uint(w), uint(h), input, resize.Bilinear)}
}

func (b *Bitmap) Pixels() []color.Color {
	c := make([]color.Color, 0)
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			c = append(c, b.Source.At(x, y))
		}
	}
	return c
}
