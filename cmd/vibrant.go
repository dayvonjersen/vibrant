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
    //palette, err := vibrant.Generate(bitmap,colors)
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
