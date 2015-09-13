package main

import "os"

import "image"
import _ "image/jpeg"

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
}
