package main

import "os"

import "image"
import _ "image/jpeg"

import "localhost/vibrant"
import "fmt"

func main() {
	f, err := os.Open("test4.jpg")
	if err != nil {
		panic(err.Error())
	}
	img, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		panic(err.Error())
	}
    bitmap := vibrant.NewBitmap(img)
    palette, err := vibrant.NewPalette(bitmap)
    if err != nil {
        panic(err.Error())
    }
/*    for _, sw := range palette.Swatches {
        fmt.Println(sw)
    }*/
    fmt.Println(palette.VibrantSwatch)
    fmt.Println(palette.DarkVibrantSwatch)
    fmt.Println(palette.LightVibrantSwatch)
    fmt.Println(palette.MutedSwatch)
    fmt.Println(palette.DarkMutedSwatch)
    fmt.Println(palette.LightMutedSwatch)
}
