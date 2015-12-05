// go port of the awesome Android Palette class:
//
// https://android.googlesource.com/platform/frameworks/support/+/b14fc7c/v7/palette/src/android/support/v7/graphics/
//
// which I translated from this beautifully cleaned up version:
//
// https://github.com/Infinity/Iris
//
// and was first made aware of by this Google I/O presentation:
//
// https://www.youtube.com/watch?v=ctzWKRlTYHQ?t=451
//
// and of course
//
// https://github.com/jariz/vibrant.js
//
// and last but not least
//
// https://github.com/akfish/node-vibrant
//
// Why:
//
// I was dissatisfied with the performance of the above JavaScript ports.
// One day, I mocked up an HTML thumbnail view gallery thing and wanted to
// use "Vibrancy" to decorate the buttons, links, titles, etc. The page had
// about 25 wallpaper images I used just as placeholders and using
// Vibrant.js brought my browser to its knees.
//
// Yes, I could have just thumbnailed the wallpapers and carried on like a
// normal person, but this act of stupidity got me thinking: why do this
// calculation in the browser? If it kills my browser like this, imagine
// what it must do to mobiles. And it's not like the images are dynamic
// (in my use case, anyway). They're dynamic in the sense that users upload
// them but once they're there, they don't change and neither does their palette.
//
// So why not do it server-side and cache the result? As a CSS file, which I
// can then use with <style scoped> e.g.:
//    <section class="card">
//        <style scoped>@import "image1-vibrant.css";</style>
//        <h1 class="darkvibrant">~Fancy Title~</h1>
//        <div class="card-border muted">
//            <img src="image1-thumb.jpg">
//        </div>
//        <div class="card-caption darkmuted">
//            <button class="vibrant">Call to action!</button>
//        </div>
//    </section>
//
// It's probably not a good idea to over-do it like that, but the point is
// that you can do it in the first place. This approach also allows for
// graceful fallback to a default palette set
//
// Example of scoped CSS (view source):
//
// http://var.abl.cl/scoped-css
//
//
// Usage:
//	package main
//	
//	import (
//	        "fmt"
//	        "image"
//	        _ "image/jpeg"
//	        "log"
//	        "os"
//	        "github.com/generaltso/vibrant"
//	)
//	
//	func checkErr(err error) {
//	        if err != nil {
//	                log.Fatalln(err)
//	        }
//	}
//	
//	func main() {
//	        file, err := os.Open("some_image.jpg")
//	        checkErr(err)
//	        img, _, err := image.Decode(file)
//	        checkErr(err)
//	        palette, err := vibrant.NewPaletteFromImage(img)
//	        checkErr(err)
//	        for name, swatch := range palette.ExtractAwesome() {
//	                fmt.Printf("name: %- 12s ", name /* or swatch.Name */)
//	                fmt.Printf("color: %s ", swatch.RGBHex())
//	                fmt.Printf("population: %d\n", swatch.Population)
//	        }
//	}
// 
// Example Output: 
//	name: DarkMuted    color: #44422f population: 1050
//	name: Vibrant      color: #f2966c population: 153
//	name: LightVibrant color: #f4b38a population: 87
//	name: DarkVibrant  color: #662f2a population: 185
//	name: Muted        color: #a2555f population: 2722
//	name: LightMuted   color: #d3aeaf population: 3469
//
//
// There is also a command-line tool in vibrant/:
//  $ cd vibrant && go install
//  $ vibrant some_image.png
//
// Example Output:
//  .lightmuted{background-color:#bec8c8;color:#000000;}
//  .darkmuted{background-color:#444230;color:#ffffff;}
//  .vibrant{background-color:#f2966c;color:#000000;}
//  .lightvibrant{background-color:#e27768;color:#000000;}
//  .darkvibrant{background-color:#783832;color:#ffffff;}
//  .muted{background-color:#a2555f;color:#ffffff;}
//
// And there is a simple demo web application which is a little more visual:
//  $ cd demo && go run app.go
//  Listening on 0.0.0.0:8080...
//
// Screenshot:
//
// http://var.abl.cl/awesome.png
//
// This API and the tools provided are a work-in-progress proof-of-concept
// and certainly could use improvement, better naming, etc.
//
// Comments, feedback and pull requests are welcome!
//
// tso@teknik.io
//
// https://github.com/generaltso/vibrant
package vibrant
