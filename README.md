# vibrant

go port of the [Android awesome Palette class](https://android.googlesource.com/platform/frameworks/support/+/b14fc7c/v7/palette/src/android/support/v7/graphics/)

which I translated from this beautifully cleaned up version:
https://github.com/Infinity/Iris

and was first made aware of by [this Google I/O 2014 presentation](https://www.youtube.com/watch?v=ctzWKRlTYHQ?t=451)

and of course
https://github.com/jariz/vibrant.js

and last but not least
https://github.com/akfish/node-vibrant

## Why

I was dissatisfied with the performance of the above JavaScript ports. One day, I mocked up an HTML thumbnail view gallery thing and wanted to use "Vibrancy" to decorate the buttons, links, titles, etc. The page had about 25 wallpaper images I used just as placeholders and using Vibrant.js brought my browser to its knees.

Yes, I could have just thumbnailed the wallpapers and carried on like a normal person, but this act of stupidity got me thinking: why do this calculation in the browser? If it kills my browser like this, imagine what it must do to mobiles. And it's not like the images are dynamic (in my use case, anyway). They're dynamic in the sense that users upload them but once they're there, they don't change and neither does their palette.

So why not do it server-side and cache the result? As a CSS file, which I can then use with `<style scoped>` e.g.:
```HTML
<section class="card">
   <style scoped>@import "image1-vibrant.css";</style>
   <h1 class="darkvibrant">~Fancy Title~</h1>
   <div class="card-border muted">
       <img src="image1-thumb.jpg">
   </div>
   <div class="card-caption darkmuted">
       <button class="vibrant">Call to action!</button>
   </div>
</section>
```


It's probably not a good idea to over-do it like that with the colors, but the point is that you *can* do it in the first place.

This approach also allows for graceful fallback to a default palette set

Perhaps I should stop trying to words **here is an example of scoped css**, view source until it makes sense:
http://var.abl.cl/scoped-css

## Usage
```bash
$ go get github.com/generaltso/vibrant
```

```go
package main

import (
	"fmt"
	"image"
    _ "image/jpeg"
	"log"
	"os"
)

import "github.com/generaltso/vibrant"

func main() {
	file, err := os.Open("some_image.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalln(err)
	}
	palette, err := vibrant.NewPaletteFromImage(img)
	if err != nil {
		log.Fatalln(err)
	}
	for name, swatch := range palette.ExtractAwesome() {
		fmt.Printf("name: %s, color: %s, population: %d\n", name /* or swatch.Name */, swatch.RGBHex(), swatch.Population)
	}
}
```

There is also a command-line tool in vibrant/:
```bash
$ cd vibrant && go install
$ vibrant some_image.png
```
(it prints CSS to stdout)

And there is a simple demo application which is a little more visual:
```bash
$ cd demo && go run app.go
Listening on 0.0.0.0:8080...
```
![screenshot of app.go](http://var.abl.cl/better.png) 
([screenshot of app.go](http://var.abl.cl/better.png))


**This API and the tools provided are a work-in-progress proof-of-concept and certainly could use improvement, better naming, etc.**

Comments, feedback and pull requests are welcome!
[email me](mailto:tso@teknik.io)
[find this project on github](https://github.com/generaltso/vibrant)
