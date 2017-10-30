package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/generaltso/vibrant"
)

var (
	input_stdin,
	output_json,
	output_css,
	output_compress,
	output_lowercase,
	output_rgb bool
)

func usage() {
	fmt.Fprintln(os.Stderr,
		`usage: vibrant [options] file
       cat image.jpg | vibrant -i [options]

options:`)
	flag.PrintDefaults()
	os.Exit(2)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.Usage = usage

	flag.BoolVar(
		&input_stdin,
		"i",
		false,
		"Read image data from stdin",
	)
	flag.BoolVar(
		&output_compress,
		"compress",
		false,
		"Strip whitespace from output.",
	)
	flag.BoolVar(
		&output_css,
		"css",
		false,
		"Output results in CSS.",
	)
	flag.BoolVar(
		&output_json,
		"json",
		false,
		"Output results in JSON.",
	)
	flag.BoolVar(
		&output_lowercase,
		"lowercase",
		true,
		"Use lowercase only for all output.",
	)
	flag.BoolVar(
		&output_rgb,
		"rgb",
		false,
		"Output RGB e.g. rgb(255,255,255) instead of HTML hex e.g. #ffffff.",
	)

	flag.Parse()

	var (
		img image.Image
		err error
	)

	if input_stdin {
		img, _, err = image.Decode(os.Stdin)
	} else {
		filename := flag.Arg(0)
		if filename == "" {
			usage()
		}

		f, err := os.Open(filename)
		checkErr(err)

		img, _, err = image.Decode(f)
		f.Close()
	}
	checkErr(err)

	palette, err := vibrant.NewPaletteFromImage(img)
	checkErr(err)

	switch {
	case output_json && output_css:
		usage()
	case output_json:
		print_json(palette)
	case output_css:
		print_css(palette)
	default:
		print_plain(palette)
	}
}

type swatch struct {
	Color, Text string
}

func print_json(palette vibrant.Palette) {
	out := map[string]interface{}{}
	for name, sw := range palette.ExtractAwesome() {
		if output_rgb {
			r, g, b := sw.Color.RGB()
			out[name] = map[string]int{"r": r, "g": g, "b": b}
		} else {
			out[name] = swatch{sw.Color.RGBHex(), sw.Color.TitleTextColor().RGBHex()}
		}
	}
	var (
		b   []byte
		err error
	)
	if output_compress {
		b, err = json.Marshal(out)
	} else {
		b, err = json.MarshalIndent(out, "", "  ")
	}
	checkErr(err)

	str := string(b)
	if output_lowercase {
		str = strings.ToLower(str)
	}
	fmt.Println(str)
}

func rgb(r ...int) string {
	return fmt.Sprintf("rgb(%d,%d,%d)", r[0], r[1], r[2])
}

func print_css(palette vibrant.Palette) {
	sp := " "
	lf := "\n"
	tb := "  "
	sc := ";"
	if output_compress {
		sp = ""
		lf = ""
		tb = ""
		sc = ""
	}
	for name, sw := range palette.ExtractAwesome() {
		var bgcolor, fgcolor string

		if output_rgb {
			bgcolor = rgb(sw.Color.RGB())
			fgcolor = rgb(sw.Color.TitleTextColor().RGB())
		} else {
			bgcolor = sw.Color.RGBHex()
			fgcolor = sw.Color.TitleTextColor().RGBHex()
		}
		if output_lowercase {
			name = strings.ToLower(name)
			bgcolor = strings.ToLower(bgcolor)
			fgcolor = strings.ToLower(fgcolor)
		}
		if output_compress && !output_rgb {
			bgcolor = shorthex(bgcolor)
			fgcolor = shorthex(fgcolor)
		}
		fmt.Printf(".%s%s{%s", name, sp, lf)
		fmt.Printf("%sbackground-color:%s%s;%s", tb, sp, bgcolor, lf)
		fmt.Printf("%scolor:%s%s%s%s}%s", tb, sp, fgcolor, sc, lf, lf)
	}
}

func shorthex(hex string) string {
	if hex[1] == hex[2] && hex[3] == hex[4] && hex[5] == hex[6] {
		return "#" + string(hex[1]) + string(hex[3]) + string(hex[5])
	}
	return hex
}

func print_plain(palette vibrant.Palette) {
	for name, sw := range palette.ExtractAwesome() {
		var fmtstr, color string
		if output_rgb {
			fmtstr = "% 12s: %- 16s (population: %d)\n"
			color = rgb(sw.Color.RGB())
		} else {
			fmtstr = "% 12s: %- 6s (population: %d)\n"
			color = sw.Color.RGBHex()
		}
		if output_lowercase {
			name = strings.ToLower(name)
		}
		if output_compress && !output_rgb {
			color = shorthex(color)
		}

		fmt.Printf(fmtstr, name, color, sw.Population)
	}
}
