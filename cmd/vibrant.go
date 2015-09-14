package main

import "os"

import "image"
import _ "image/jpeg"

import "localhost/vibrant"
import "fmt"

//import "strconv"

func main() {
	if len(os.Args) <= 1 {
		//println("usage: vibrant [input image] [maxColors]")
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
	bitmap := vibrant.NewBitmap(img)
	//colors, _ := strconv.Atoi(os.Args[2])
	palette, err := vibrant.Generate(bitmap, 8192)
	//palette, err := vibrant.NewPalette(bitmap)
	if err != nil {
		panic(err.Error())
	}
	/*    for _, sw := range palette.Swatches {
	      fmt.Println(sw)
	  }*/
	if palette.VibrantSwatch != nil {
		fmt.Println(palette.VibrantSwatch)
	}
	if palette.DarkVibrantSwatch != nil {
		fmt.Println(palette.DarkVibrantSwatch)
	}
	if palette.LightVibrantSwatch != nil {
		fmt.Println(palette.LightVibrantSwatch)
	}
	if palette.MutedSwatch != nil {
		fmt.Println(palette.MutedSwatch)
	}
	if palette.DarkMutedSwatch != nil {
		fmt.Println(palette.DarkMutedSwatch)
	}
	if palette.LightMutedSwatch != nil {
		fmt.Println(palette.LightMutedSwatch)
	}
}
