package vibrant

import "github.com/nfnt/resize"
import "image"
import "image/color"
import "math"

// type Bitmap is a simple wrapper for an image.Image
type Bitmap struct {
	Width  int
	Height int
	Source image.Image
}

func NewBitmap(input image.Image) *Bitmap {
	bounds := input.Bounds()
	return &Bitmap{bounds.Dx(), bounds.Dy(), input}
}

// Scales input image.Image by ratio using github.com/nfnt/resize
func NewScaledBitmap(input image.Image, ratio float64) *Bitmap {
	bounds := input.Bounds()
	w := math.Ceil(float64(bounds.Dx()) * ratio)
	h := math.Ceil(float64(bounds.Dy()) * ratio)
	return &Bitmap{int(w), int(h), resize.Resize(uint(w), uint(h), input, resize.Bilinear)}
}

// Returns all of the pixels of this Bitmap.Source as a 1D array of color.Color's
func (b *Bitmap) Pixels() []color.Color {
	c := make([]color.Color, 0)
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			c = append(c, b.Source.At(x, y))
		}
	}
	return c
}
