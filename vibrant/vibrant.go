package main

import (
	"fmt"
	"github.com/generaltso/vibrant"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		println("usage: vibrant [input image]")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	img, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	palette, err := vibrant.NewPaletteFromImage(img)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	for _, sw := range palette.ExtractAwesome() {
		fmt.Println(sw.String())
	}
}
