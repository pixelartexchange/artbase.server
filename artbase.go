package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"image/draw"
	"os"
	"log"
	"errors"
  "strings"
	"strconv"
	"bytes"
	"net/http"
	"io"
	// "io/ioutil"
	"html/template"
	"github.com/gin-gonic/gin"
)



type collection struct {
	Name         string
	Width        int
	Height       int
	Path         string
	Url          string
	// note:  background==false (default) => transparent
	//        background==true            => images have backgrounds (NOT transparent)
	Background   bool
	Count        int
}


var collections = []collection{
  {Name: "punks",     Width: 24, Height: 24,
	 Path: "./punks.png",
	 Url: "https://github.com/cryptopunksnotdead/awesome-24px/raw/master/collection/punks.png" },

	 {Name: "morepunks",  Width: 24, Height: 24,
	 Path: "./morepunks.png",
	 Url: "https://github.com/cryptopunksnotdead/awesome-24px/raw/master/collection/morepunks.png" },

	 {Name: "readymadepunks",  Width: 24, Height: 24,
	 Path: "./readymadepunks.png",
	 Url: "https://github.com/cryptopunksnotdead/punks.readymade/raw/master/readymades.png" },


	 {Name: "bwpunks",  Width: 24, Height: 24,
	  Path: "./bwpunks.png",
	  Url: "https://github.com/pixelartexchange/collections/raw/master/bwpunks-24x24.png",
	  Background: true },

		{Name: "frontpunks",  Width: 24, Height: 24,
	  Path: "./frontpunks.png",
	  Url: "https://github.com/cryptopunksnotdead/awesome-24px/raw/master/collection/frontpunks.png" },

		{Name: "intlpunks",   Width: 24, Height: 24,
	  Path: "./intlpunks.png",
    Url: "https://github.com/cryptopunksnotdead/awesome-24px/raw/master/collection/intlpunks.png" },

		{Name: "boredapes",   Width: 28, Height: 28,
	  Path: "./boredapes.png",
    Url: "https://github.com/cryptopunksnotdead/awesome-24px/raw/master/collection/boredapes.png" },

		{Name: "apes",   Width: 35, Height: 35,
 	  Path: "./apes.png",
     Url: "https://github.com/pixelartexchange/collections/raw/master/apes/apes-35x35.png",
		 Background: true },

		 {Name: "basicboredapes",   Width: 50, Height: 50,
 	  Path: "./basicboredapes.png",
     Url:	"https://github.com/pixelartexchange/collections/raw/master/basicboredapeclub/basicboredapeclub-50x50.png",
		 Background: true },


	 {Name: "coolcats",  Width: 24, Height: 24,
	 Path: "./coolcats.png",
	 Url: "https://github.com/cryptopunksnotdead/awesome-24px/raw/master/collection/coolcats.png" },

	 {Name: "doge",  Width: 24, Height: 24,
	 Path: "./doge.png",
   Url: "https://github.com/cryptopunksnotdead/programming-cryptopunks/raw/master/i/doge.png" },

	 {Name: "dooggies",  Width: 32, Height: 32,
	 Path: "./dooggies.png",
   Url: "https://github.com/pixelartexchange/collections/raw/master/dooggies-32x32.png",
	 Background: true },


	 {Name: "blockydoge",  Width: 60, Height: 60,
	 Path: "./blockydoge.png",
   Url: "https://github.com/pixelartexchange/collections/raw/master/blockydoge/blockydoge-60x60.png",
	 Background: true },


	 {Name: "wiener",  Width: 32, Height: 32,
	 Path: "./wiener.png",
   Url: "https://github.com/pixelartexchange/collections/raw/master/wiener-32x32.png",
	 Background: true },

	 {Name: "rocks",  Width: 24, Height: 24,
	 Path: "./rocks.png",
   Url: "https://github.com/cryptopunksnotdead/awesome-24px/raw/master/collection/rocks.png" },

	 {Name: "punkrocks",  Width: 24, Height: 24,
	 Path: "./punkrocks.png",
   Url: "https://github.com/cryptopunksnotdead/awesome-24px/raw/master/collection/punkrocks.png" },



}


const templateHome = `
<h1>{{ len . }} Collections</h1>
<ul>
    {{range .}}
        <li>
				  <a href="/{{.Name}}">/{{.Name}}</a>
					({{ .Width}}×{{ .Height}})

          {{if .Background}}
					   - incl. backgrounds
          {{end}}

					<span style="font-size: 80%">
					  <a href="{{ .Url }}">(Download .png)</a>
					</span>
				</li>
    {{end}}
</ul>


<hr>

<p style="font-size: 80%">
  New to artbase? Find out more at
	<a href="https://github.com/pixelartexchange/artbase">/artbase »</a>
</p>
`


const templateIndex = `
<p style="font-size: 80%">
  <a href="/">« Collections</a>
</p>

<h1>{{ .Name}}   ({{ .Width}}×{{ .Height}}) Collection</h1>

<p>
  To get images, use /{{ .Name}}/<em>:id</em>
</p>

<p>
	Example:
	<a href="/{{ .Name}}/0">/{{ .Name}}/0</a>,
	<a href="/{{ .Name}}/1">/{{ .Name}}/1</a>,
	<a href="/{{ .Name}}/2">/{{ .Name}}/2</a>, ...
</p>

<p>
	Note: The default image size is ({{ .Width}}×{{ .Height}}).
	 Use the z (zoom) parameter to upsize.
</p>
<p>
	 Try 2x:
	 <a href="/{{ .Name}}/0?z=2">/{{ .Name}}/0?z=2</a>,
	 <a href="/{{ .Name}}/1?z=2">/{{ .Name}}/1?z=2</a>,
	 <a href="/{{ .Name}}/2?z=2">/{{ .Name}}/2?z=2</a>, ...
</p>

<p>
	 Try 8x:
	 <a href="/{{ .Name}}/0?z=8">/{{ .Name}}/0?z=8</a>,
	 <a href="/{{ .Name}}/1?z=8">/{{ .Name}}/1?z=8</a>,
	 <a href="/{{ .Name}}/2?z=8">/{{ .Name}}/2?z=8</a>, ...

	 And so on.
</p>



<hr>

<p style="font-size: 80%">
  New to artbase? Find out more at
	<a href="https://github.com/pixelartexchange/artbase">/artbase »</a>
</p>

`



var templates map[string]*template.Template

func compileTemplates() {
	if templates == nil {
		templates = make( map[string]*template.Template )

    templates["home"]  = template.Must( template.New("home").Parse( templateHome ))
		templates["index"] = template.Must( template.New("index").Parse( templateIndex ))
	}
}


func renderHome( data []collection ) []byte {
	buf := new( bytes.Buffer )
	templates["home"].Execute( buf, data )
	return buf.Bytes()
}

func renderCollection( data *collection ) []byte {
	buf := new( bytes.Buffer )
	templates["index"].Execute( buf, data )
	return buf.Bytes()
}





// check if divod exists built-in - different name or such ??
func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient  = numerator / denominator   // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}




func downloadImage( url, outpath string ) {
  fmt.Printf( "==> Downloading %s...\n", url )

	resp, err := http.Get( url )
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// todo/fix:
	//  (double) check for content type - why? why not?

	fmt.Printf( "HTTP %v\n", resp.StatusCode )
  // dump  headers for debugging
  for name, headers := range resp.Header {
	  for _, h := range headers {
		  fmt.Printf( "  %v  :  %v\n", name, h )
	  }
  }

	 f, err := os.Create( outpath )
	 if err != nil {
		 log.Fatal(err)
	 }
	 defer f.Close()



   n, err := io.Copy( f, resp.Body )   // note:  Copy( writer, reader )
	 if err != nil {
		  fmt.Printf( "  !! ERROR - writing %d byte(s) to %s...\n", n, outpath )
		  log.Fatal(err)
	 }
	 fmt.Printf( "  writing %d byte(s) to %s...\n", n, outpath )
}






// Map contains named colors defined in the SVG 1.1 spec.
//   see https://raw.githubusercontent.com/golang/image/master/colornames/table.go

var colorMap = map[string]color.RGBA{
	"aliceblue":            color.RGBA{0xf0, 0xf8, 0xff, 0xff}, // rgb(240, 248, 255)
	"antiquewhite":         color.RGBA{0xfa, 0xeb, 0xd7, 0xff}, // rgb(250, 235, 215)
	"aqua":                 color.RGBA{0x00, 0xff, 0xff, 0xff}, // rgb(0, 255, 255)
	"aquamarine":           color.RGBA{0x7f, 0xff, 0xd4, 0xff}, // rgb(127, 255, 212)
	"azure":                color.RGBA{0xf0, 0xff, 0xff, 0xff}, // rgb(240, 255, 255)
	"beige":                color.RGBA{0xf5, 0xf5, 0xdc, 0xff}, // rgb(245, 245, 220)
	"bisque":               color.RGBA{0xff, 0xe4, 0xc4, 0xff}, // rgb(255, 228, 196)
	"black":                color.RGBA{0x00, 0x00, 0x00, 0xff}, // rgb(0, 0, 0)
	"blanchedalmond":       color.RGBA{0xff, 0xeb, 0xcd, 0xff}, // rgb(255, 235, 205)
	"blue":                 color.RGBA{0x00, 0x00, 0xff, 0xff}, // rgb(0, 0, 255)
	"blueviolet":           color.RGBA{0x8a, 0x2b, 0xe2, 0xff}, // rgb(138, 43, 226)
	"brown":                color.RGBA{0xa5, 0x2a, 0x2a, 0xff}, // rgb(165, 42, 42)
	"burlywood":            color.RGBA{0xde, 0xb8, 0x87, 0xff}, // rgb(222, 184, 135)
	"cadetblue":            color.RGBA{0x5f, 0x9e, 0xa0, 0xff}, // rgb(95, 158, 160)
	"chartreuse":           color.RGBA{0x7f, 0xff, 0x00, 0xff}, // rgb(127, 255, 0)
	"chocolate":            color.RGBA{0xd2, 0x69, 0x1e, 0xff}, // rgb(210, 105, 30)
	"coral":                color.RGBA{0xff, 0x7f, 0x50, 0xff}, // rgb(255, 127, 80)
	"cornflowerblue":       color.RGBA{0x64, 0x95, 0xed, 0xff}, // rgb(100, 149, 237)
	"cornsilk":             color.RGBA{0xff, 0xf8, 0xdc, 0xff}, // rgb(255, 248, 220)
	"crimson":              color.RGBA{0xdc, 0x14, 0x3c, 0xff}, // rgb(220, 20, 60)
	"cyan":                 color.RGBA{0x00, 0xff, 0xff, 0xff}, // rgb(0, 255, 255)
	"darkblue":             color.RGBA{0x00, 0x00, 0x8b, 0xff}, // rgb(0, 0, 139)
	"darkcyan":             color.RGBA{0x00, 0x8b, 0x8b, 0xff}, // rgb(0, 139, 139)
	"darkgoldenrod":        color.RGBA{0xb8, 0x86, 0x0b, 0xff}, // rgb(184, 134, 11)
	"darkgray":             color.RGBA{0xa9, 0xa9, 0xa9, 0xff}, // rgb(169, 169, 169)
	"darkgreen":            color.RGBA{0x00, 0x64, 0x00, 0xff}, // rgb(0, 100, 0)
	"darkgrey":             color.RGBA{0xa9, 0xa9, 0xa9, 0xff}, // rgb(169, 169, 169)
	"darkkhaki":            color.RGBA{0xbd, 0xb7, 0x6b, 0xff}, // rgb(189, 183, 107)
	"darkmagenta":          color.RGBA{0x8b, 0x00, 0x8b, 0xff}, // rgb(139, 0, 139)
	"darkolivegreen":       color.RGBA{0x55, 0x6b, 0x2f, 0xff}, // rgb(85, 107, 47)
	"darkorange":           color.RGBA{0xff, 0x8c, 0x00, 0xff}, // rgb(255, 140, 0)
	"darkorchid":           color.RGBA{0x99, 0x32, 0xcc, 0xff}, // rgb(153, 50, 204)
	"darkred":              color.RGBA{0x8b, 0x00, 0x00, 0xff}, // rgb(139, 0, 0)
	"darksalmon":           color.RGBA{0xe9, 0x96, 0x7a, 0xff}, // rgb(233, 150, 122)
	"darkseagreen":         color.RGBA{0x8f, 0xbc, 0x8f, 0xff}, // rgb(143, 188, 143)
	"darkslateblue":        color.RGBA{0x48, 0x3d, 0x8b, 0xff}, // rgb(72, 61, 139)
	"darkslategray":        color.RGBA{0x2f, 0x4f, 0x4f, 0xff}, // rgb(47, 79, 79)
	"darkslategrey":        color.RGBA{0x2f, 0x4f, 0x4f, 0xff}, // rgb(47, 79, 79)
	"darkturquoise":        color.RGBA{0x00, 0xce, 0xd1, 0xff}, // rgb(0, 206, 209)
	"darkviolet":           color.RGBA{0x94, 0x00, 0xd3, 0xff}, // rgb(148, 0, 211)
	"deeppink":             color.RGBA{0xff, 0x14, 0x93, 0xff}, // rgb(255, 20, 147)
	"deepskyblue":          color.RGBA{0x00, 0xbf, 0xff, 0xff}, // rgb(0, 191, 255)
	"dimgray":              color.RGBA{0x69, 0x69, 0x69, 0xff}, // rgb(105, 105, 105)
	"dimgrey":              color.RGBA{0x69, 0x69, 0x69, 0xff}, // rgb(105, 105, 105)
	"dodgerblue":           color.RGBA{0x1e, 0x90, 0xff, 0xff}, // rgb(30, 144, 255)
	"firebrick":            color.RGBA{0xb2, 0x22, 0x22, 0xff}, // rgb(178, 34, 34)
	"floralwhite":          color.RGBA{0xff, 0xfa, 0xf0, 0xff}, // rgb(255, 250, 240)
	"forestgreen":          color.RGBA{0x22, 0x8b, 0x22, 0xff}, // rgb(34, 139, 34)
	"fuchsia":              color.RGBA{0xff, 0x00, 0xff, 0xff}, // rgb(255, 0, 255)
	"gainsboro":            color.RGBA{0xdc, 0xdc, 0xdc, 0xff}, // rgb(220, 220, 220)
	"ghostwhite":           color.RGBA{0xf8, 0xf8, 0xff, 0xff}, // rgb(248, 248, 255)
	"gold":                 color.RGBA{0xff, 0xd7, 0x00, 0xff}, // rgb(255, 215, 0)
	"goldenrod":            color.RGBA{0xda, 0xa5, 0x20, 0xff}, // rgb(218, 165, 32)
	"gray":                 color.RGBA{0x80, 0x80, 0x80, 0xff}, // rgb(128, 128, 128)
	"green":                color.RGBA{0x00, 0x80, 0x00, 0xff}, // rgb(0, 128, 0)
	"greenyellow":          color.RGBA{0xad, 0xff, 0x2f, 0xff}, // rgb(173, 255, 47)
	"grey":                 color.RGBA{0x80, 0x80, 0x80, 0xff}, // rgb(128, 128, 128)
	"honeydew":             color.RGBA{0xf0, 0xff, 0xf0, 0xff}, // rgb(240, 255, 240)
	"hotpink":              color.RGBA{0xff, 0x69, 0xb4, 0xff}, // rgb(255, 105, 180)
	"indianred":            color.RGBA{0xcd, 0x5c, 0x5c, 0xff}, // rgb(205, 92, 92)
	"indigo":               color.RGBA{0x4b, 0x00, 0x82, 0xff}, // rgb(75, 0, 130)
	"ivory":                color.RGBA{0xff, 0xff, 0xf0, 0xff}, // rgb(255, 255, 240)
	"khaki":                color.RGBA{0xf0, 0xe6, 0x8c, 0xff}, // rgb(240, 230, 140)
	"lavender":             color.RGBA{0xe6, 0xe6, 0xfa, 0xff}, // rgb(230, 230, 250)
	"lavenderblush":        color.RGBA{0xff, 0xf0, 0xf5, 0xff}, // rgb(255, 240, 245)
	"lawngreen":            color.RGBA{0x7c, 0xfc, 0x00, 0xff}, // rgb(124, 252, 0)
	"lemonchiffon":         color.RGBA{0xff, 0xfa, 0xcd, 0xff}, // rgb(255, 250, 205)
	"lightblue":            color.RGBA{0xad, 0xd8, 0xe6, 0xff}, // rgb(173, 216, 230)
	"lightcoral":           color.RGBA{0xf0, 0x80, 0x80, 0xff}, // rgb(240, 128, 128)
	"lightcyan":            color.RGBA{0xe0, 0xff, 0xff, 0xff}, // rgb(224, 255, 255)
	"lightgoldenrodyellow": color.RGBA{0xfa, 0xfa, 0xd2, 0xff}, // rgb(250, 250, 210)
	"lightgray":            color.RGBA{0xd3, 0xd3, 0xd3, 0xff}, // rgb(211, 211, 211)
	"lightgreen":           color.RGBA{0x90, 0xee, 0x90, 0xff}, // rgb(144, 238, 144)
	"lightgrey":            color.RGBA{0xd3, 0xd3, 0xd3, 0xff}, // rgb(211, 211, 211)
	"lightpink":            color.RGBA{0xff, 0xb6, 0xc1, 0xff}, // rgb(255, 182, 193)
	"lightsalmon":          color.RGBA{0xff, 0xa0, 0x7a, 0xff}, // rgb(255, 160, 122)
	"lightseagreen":        color.RGBA{0x20, 0xb2, 0xaa, 0xff}, // rgb(32, 178, 170)
	"lightskyblue":         color.RGBA{0x87, 0xce, 0xfa, 0xff}, // rgb(135, 206, 250)
	"lightslategray":       color.RGBA{0x77, 0x88, 0x99, 0xff}, // rgb(119, 136, 153)
	"lightslategrey":       color.RGBA{0x77, 0x88, 0x99, 0xff}, // rgb(119, 136, 153)
	"lightsteelblue":       color.RGBA{0xb0, 0xc4, 0xde, 0xff}, // rgb(176, 196, 222)
	"lightyellow":          color.RGBA{0xff, 0xff, 0xe0, 0xff}, // rgb(255, 255, 224)
	"lime":                 color.RGBA{0x00, 0xff, 0x00, 0xff}, // rgb(0, 255, 0)
	"limegreen":            color.RGBA{0x32, 0xcd, 0x32, 0xff}, // rgb(50, 205, 50)
	"linen":                color.RGBA{0xfa, 0xf0, 0xe6, 0xff}, // rgb(250, 240, 230)
	"magenta":              color.RGBA{0xff, 0x00, 0xff, 0xff}, // rgb(255, 0, 255)
	"maroon":               color.RGBA{0x80, 0x00, 0x00, 0xff}, // rgb(128, 0, 0)
	"mediumaquamarine":     color.RGBA{0x66, 0xcd, 0xaa, 0xff}, // rgb(102, 205, 170)
	"mediumblue":           color.RGBA{0x00, 0x00, 0xcd, 0xff}, // rgb(0, 0, 205)
	"mediumorchid":         color.RGBA{0xba, 0x55, 0xd3, 0xff}, // rgb(186, 85, 211)
	"mediumpurple":         color.RGBA{0x93, 0x70, 0xdb, 0xff}, // rgb(147, 112, 219)
	"mediumseagreen":       color.RGBA{0x3c, 0xb3, 0x71, 0xff}, // rgb(60, 179, 113)
	"mediumslateblue":      color.RGBA{0x7b, 0x68, 0xee, 0xff}, // rgb(123, 104, 238)
	"mediumspringgreen":    color.RGBA{0x00, 0xfa, 0x9a, 0xff}, // rgb(0, 250, 154)
	"mediumturquoise":      color.RGBA{0x48, 0xd1, 0xcc, 0xff}, // rgb(72, 209, 204)
	"mediumvioletred":      color.RGBA{0xc7, 0x15, 0x85, 0xff}, // rgb(199, 21, 133)
	"midnightblue":         color.RGBA{0x19, 0x19, 0x70, 0xff}, // rgb(25, 25, 112)
	"mintcream":            color.RGBA{0xf5, 0xff, 0xfa, 0xff}, // rgb(245, 255, 250)
	"mistyrose":            color.RGBA{0xff, 0xe4, 0xe1, 0xff}, // rgb(255, 228, 225)
	"moccasin":             color.RGBA{0xff, 0xe4, 0xb5, 0xff}, // rgb(255, 228, 181)
	"navajowhite":          color.RGBA{0xff, 0xde, 0xad, 0xff}, // rgb(255, 222, 173)
	"navy":                 color.RGBA{0x00, 0x00, 0x80, 0xff}, // rgb(0, 0, 128)
	"oldlace":              color.RGBA{0xfd, 0xf5, 0xe6, 0xff}, // rgb(253, 245, 230)
	"olive":                color.RGBA{0x80, 0x80, 0x00, 0xff}, // rgb(128, 128, 0)
	"olivedrab":            color.RGBA{0x6b, 0x8e, 0x23, 0xff}, // rgb(107, 142, 35)
	"orange":               color.RGBA{0xff, 0xa5, 0x00, 0xff}, // rgb(255, 165, 0)
	"orangered":            color.RGBA{0xff, 0x45, 0x00, 0xff}, // rgb(255, 69, 0)
	"orchid":               color.RGBA{0xda, 0x70, 0xd6, 0xff}, // rgb(218, 112, 214)
	"palegoldenrod":        color.RGBA{0xee, 0xe8, 0xaa, 0xff}, // rgb(238, 232, 170)
	"palegreen":            color.RGBA{0x98, 0xfb, 0x98, 0xff}, // rgb(152, 251, 152)
	"paleturquoise":        color.RGBA{0xaf, 0xee, 0xee, 0xff}, // rgb(175, 238, 238)
	"palevioletred":        color.RGBA{0xdb, 0x70, 0x93, 0xff}, // rgb(219, 112, 147)
	"papayawhip":           color.RGBA{0xff, 0xef, 0xd5, 0xff}, // rgb(255, 239, 213)
	"peachpuff":            color.RGBA{0xff, 0xda, 0xb9, 0xff}, // rgb(255, 218, 185)
	"peru":                 color.RGBA{0xcd, 0x85, 0x3f, 0xff}, // rgb(205, 133, 63)
	"pink":                 color.RGBA{0xff, 0xc0, 0xcb, 0xff}, // rgb(255, 192, 203)
	"plum":                 color.RGBA{0xdd, 0xa0, 0xdd, 0xff}, // rgb(221, 160, 221)
	"powderblue":           color.RGBA{0xb0, 0xe0, 0xe6, 0xff}, // rgb(176, 224, 230)
	"purple":               color.RGBA{0x80, 0x00, 0x80, 0xff}, // rgb(128, 0, 128)
	"red":                  color.RGBA{0xff, 0x00, 0x00, 0xff}, // rgb(255, 0, 0)
	"rosybrown":            color.RGBA{0xbc, 0x8f, 0x8f, 0xff}, // rgb(188, 143, 143)
	"royalblue":            color.RGBA{0x41, 0x69, 0xe1, 0xff}, // rgb(65, 105, 225)
	"saddlebrown":          color.RGBA{0x8b, 0x45, 0x13, 0xff}, // rgb(139, 69, 19)
	"salmon":               color.RGBA{0xfa, 0x80, 0x72, 0xff}, // rgb(250, 128, 114)
	"sandybrown":           color.RGBA{0xf4, 0xa4, 0x60, 0xff}, // rgb(244, 164, 96)
	"seagreen":             color.RGBA{0x2e, 0x8b, 0x57, 0xff}, // rgb(46, 139, 87)
	"seashell":             color.RGBA{0xff, 0xf5, 0xee, 0xff}, // rgb(255, 245, 238)
	"sienna":               color.RGBA{0xa0, 0x52, 0x2d, 0xff}, // rgb(160, 82, 45)
	"silver":               color.RGBA{0xc0, 0xc0, 0xc0, 0xff}, // rgb(192, 192, 192)
	"skyblue":              color.RGBA{0x87, 0xce, 0xeb, 0xff}, // rgb(135, 206, 235)
	"slateblue":            color.RGBA{0x6a, 0x5a, 0xcd, 0xff}, // rgb(106, 90, 205)
	"slategray":            color.RGBA{0x70, 0x80, 0x90, 0xff}, // rgb(112, 128, 144)
	"slategrey":            color.RGBA{0x70, 0x80, 0x90, 0xff}, // rgb(112, 128, 144)
	"snow":                 color.RGBA{0xff, 0xfa, 0xfa, 0xff}, // rgb(255, 250, 250)
	"springgreen":          color.RGBA{0x00, 0xff, 0x7f, 0xff}, // rgb(0, 255, 127)
	"steelblue":            color.RGBA{0x46, 0x82, 0xb4, 0xff}, // rgb(70, 130, 180)
	"tan":                  color.RGBA{0xd2, 0xb4, 0x8c, 0xff}, // rgb(210, 180, 140)
	"teal":                 color.RGBA{0x00, 0x80, 0x80, 0xff}, // rgb(0, 128, 128)
	"thistle":              color.RGBA{0xd8, 0xbf, 0xd8, 0xff}, // rgb(216, 191, 216)
	"tomato":               color.RGBA{0xff, 0x63, 0x47, 0xff}, // rgb(255, 99, 71)
	"turquoise":            color.RGBA{0x40, 0xe0, 0xd0, 0xff}, // rgb(64, 224, 208)
	"violet":               color.RGBA{0xee, 0x82, 0xee, 0xff}, // rgb(238, 130, 238)
	"wheat":                color.RGBA{0xf5, 0xde, 0xb3, 0xff}, // rgb(245, 222, 179)
	"white":                color.RGBA{0xff, 0xff, 0xff, 0xff}, // rgb(255, 255, 255)
	"whitesmoke":           color.RGBA{0xf5, 0xf5, 0xf5, 0xff}, // rgb(245, 245, 245)
	"yellow":               color.RGBA{0xff, 0xff, 0x00, 0xff}, // rgb(255, 255, 0)
	"yellowgreen":          color.RGBA{0x9a, 0xcd, 0x32, 0xff}, // rgb(154, 205, 50)

	// our own custom colors
  // todo - use add alias punksv2, punksv3 - why? why not?
	"v2":              color.RGBA{0x63, 0x85, 0x96, 0xff},   // larva gray #638596?
	"v3":              color.RGBA{0x60, 0xa4, 0xf7, 0xff},   // baby blue  #60A4F7
	// todo: add v1,v4 - any others?
}


func parseColor(s string) (color.RGBA, error) {
	// try color map with known color names first
	if color, ok := colorMap[ s ]; ok {
    return color, nil
	} else {
		// assume hexstring
		return parseHexColor( s )
	}
}


func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	switch len(s) {
	case 6:
		_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 3:
		_, err = fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		//    0xf (15) * 0x11 (17) = 0xff (255)
		//    0xb (11) * 0x11 (17) = 0xbb (187)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}





func zoomImage(img image.Image, zoom int) (*image.RGBA, error) {

	bounds := img.Bounds()
	width, height := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	fmt.Println( bounds, width, height )
	// e.g.   punk #0    (0,0)-(24,24)
  //             #561  (1464,120)-(1488,144)
  //             #3100 (0,744)-(24,768)
  //             #7804 (96,1872)-(120,1896)
  //             #8857 (1368,2112)-(1392,2136)

	new_img := image.NewRGBA( image.Rect(0,0, width*zoom, height*zoom) )

	for x:=0; x < width; x++ {
		for y:=0; y < height; y++ {
				pixel := img.At( bounds.Min.X+x, bounds.Min.Y+y )
        for n:=0; n < zoom; n++ {
					for m:=0; m < zoom; m++ {
						new_img.Set( n+zoom*x, m+zoom*y, pixel )
					}
				}
	 }
}

	return new_img, nil
}


func mirrorImage(img image.Image) (*image.RGBA, error) {

	bounds := img.Bounds()
	width, height := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	fmt.Println( bounds, width, height )
	// e.g.   punk #0    (0,0)-(24,24)
  //             #561  (1464,120)-(1488,144)
  //             #3100 (0,744)-(24,768)
  //             #7804 (96,1872)-(120,1896)
  //             #8857 (1368,2112)-(1392,2136)

	new_img := image.NewRGBA( image.Rect(0,0, width, height) )

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At( bounds.Min.X+x, bounds.Min.Y+y )
			new_img.Set( (width-1)-x, y, pixel )
		}
	}

	return new_img, nil
}




// todo/check: rename pixel_raw to pixel_uncast/generic/... or similar - why? why not?
func renderPixel( x, y int,
	                pixel color.NRGBA,
									colors map[string]*ColorInfo )  string {

	// convert color to rgba() hexstring
  hex := fmt.Sprintf( "#%02x%02x%02x%02x",
	                pixel.R, pixel.G, pixel.B, pixel.A )

	//return fmt.Sprintf( "<rect class=\"px\" x=\"%d\" y=\"%d\" width=\"1\" height=\"1\" fill=\"%s\"/>",
  //                          x, y, hex )

  colInfo, _ := colors[ hex ]

	return fmt.Sprintf( "<rect class=\"px c%d\" x=\"%d\" y=\"%d\" width=\"1\" height=\"1\"/>",
                            colInfo.Index+1, x, y )


}





// todo/check: find a better name - why? why not?
type ColorInfo struct {
	 Count int      // no. of pixels
	 Index int      // zero-based running index (0,1,2, etc.)
}

func imageColors( img image.Image ) map[string]*ColorInfo  {

  colors := make( map[string]*ColorInfo )


	transparent := color.NRGBA{ R: 0,
															 G: 0,
															 B: 0,
																A: 0 }

  bounds := img.Bounds()

	for y:=bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x:=bounds.Min.X; x < bounds.Max.X; x++ {
			// note: cast/convert interface image.color to type NRGBA (supporting image.color)
			// pixel, ok := img.At( bounds.Min.X+x, bounds.Min.Y+y ).(image.NRGBA)
			// r, g, b, a := pixel.R, pixel.G, pixel.B, pixel.A
			pixel := color.NRGBAModel.Convert( img.At( x, y )).(color.NRGBA)
			// if !ok {
			//	fmt.Printf( "!! ERROR - expected color.NRGBA pixel/color but got: %v %T\n",
			//																		 pixel, pixel )
			//	os.Exit(1)
			//}


			if pixel == transparent {
				 // do nothing; skip
			}  else {
				// convert color to rgba() hexstring
				hex := fmt.Sprintf( "#%02x%02x%02x%02x",
												pixel.R, pixel.G, pixel.B, pixel.A )

         colInfo, ok := colors[ hex ]
				 if !ok {
            colInfo = &ColorInfo{ Index: len( colors ) }
						colors[ hex ] = colInfo
				 }

				 colInfo.Count += 1
				 fmt.Printf( "%s - %v %v\n", hex, colInfo, *colInfo )
				}
	 }
	}

	return colors
}



func imageToSVG( img image.Image ) string  {

		/////
		// type NRGBA
		//   type NRGBA struct {
	  //     R, G, B, A uint8
    //   }
    //  NRGBA represents a non-alpha-premultiplied 32-bit color.

		// The NRGBA struct type represents an 8-bit non-alpha-premultiplied color,
		//  as used by the PNG image format.
		// When manipulating an NRGBA’s fields directly,
		// the values are non-alpha-premultiplied,
		// but when calling the RGBA method, the return values are alpha-premultiplied.


    // transparent := color.NRGBA{ 0, 0, 0, 0 }
		transparent := color.NRGBA{ R: 0,
			G: 0,
			B: 0,
			 A: 0 }

    // todo/check: try adding
		//   height="24px/240px" width="24px/240px"
    //		viewBox="0 0 24 24"
		//   style="background: red;">


		colors := imageColors( img )

		fmt.Println( "colors:" )
		fmt.Println( colors )


		bounds        := img.Bounds()
		width, height := bounds.Dx(), bounds.Dy()


    var buf bytes.Buffer
    buf.WriteString( "<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.2\"\n" )
		buf.WriteString(
			fmt.Sprintf( "    width=\"%dpx\" height=\"%dpx\" viewBox=\"0 0 %d %d\">\n",
			   width, height, width, height ))
		buf.WriteString( "  <style>\n" )

		// add colors  - note: start with one (1) NOT zero (0)-based - why? why not?
		//    e.g rect.c1 { fill: #000000ff }
    for colHex,colInfo := range colors {
			buf.WriteString(
				fmt.Sprintf( "    rect.c%d {  fill: %s }  /* %3d pixel(s) */\n",
			         colInfo.Index+1, colHex, colInfo.Count ))
		}

		buf.WriteString( "    rect.px {  shape-rendering: crispEdges  }\n" )
		buf.WriteString( "  </style>\n" )



		for y:=0; y < height; y++ {
		  for x:=0; x < width; x++ {
				// note: cast/convert interface image.color to type NRGBA (supporting image.color)
				// pixel, ok := img.At( bounds.Min.X+x, bounds.Min.Y+y ).(image.NRGBA)
				// r, g, b, a := pixel.R, pixel.G, pixel.B, pixel.A
				pixel := color.NRGBAModel.Convert(
					            img.At( bounds.Min.X+x, bounds.Min.Y+y )).(color.NRGBA)


					if pixel == transparent {
				}  else {
					   inst := renderPixel( x, y, pixel, colors )
					   line := fmt.Sprintf( "  %s\n", inst )

					buf.WriteString( line )
				}
		 }
		}

		buf.WriteString( "</svg>\n" )

		return buf.String()
}





func fileExist(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
			return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
			return false, nil
	}
	return false, err
}


func read_image( path string ) *image.Image  {
	fmt.Printf( "==> reading %s...\n", path )

	f, err := os.Open( path )
	if err != nil {
			log.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode( f )
	if err != nil {
			log.Fatal(err)
	}

	bounds := img.Bounds()
	fmt.Println( bounds )
	// e.g.   punks.png  (0,0)-(2400,2400)

	return &img;
}



const (
	ContentTypeJSON     = "application/json"
	ContentTypeHTML     = "text/html; charset=utf-8"
	ContentTypeText     = "text/plain; charset=utf-8"
	ContentTypeImagePNG = "image/png"
	ContentTypePNG      = ContentTypeImagePNG
	ContentTypeImageSVG = "image/svg+xml"
	ContentTypeSVG      = ContentTypeImageSVG
)



var cache map[string]*image.Image


func handleHome( collections []collection ) gin.HandlerFunc  {
  return func( ctx *gin.Context ) {
		 b := renderHome( collections )
		 ctx.Data( http.StatusOK, ContentTypeHTML, b )
	}
}

func handleCollection( col collection ) gin.HandlerFunc  {
  return func( ctx *gin.Context ) {
		 b := renderCollection( &col )
		 ctx.Data( http.StatusOK, ContentTypeHTML, b )
	}
}



func handleCollectionImage( col collection ) gin.HandlerFunc  {
	// check if collection is in cache
  return func( ctx *gin.Context ) {

		path := col.Path

	if exist, _ := fileExist( path ); !exist  {
    fmt.Println( "  getting composite / download to (local) cache..." )

		url := col.Url
		downloadImage( url, path )
	}

	name                    := col.Name
	tile_width, tile_height := col.Width, col.Height   // in px


  var imgptr *image.Image


	if cache[ name ] != nil {
		 fmt.Println( "    bingo! (re)using in-memory composite image from cache...\n" )
     imgptr = cache[ name ]
	} else {
		 fmt.Println( "   adding composite image to in-memory cache...\n" )
	   imgptr = read_image( path )
		 cache[ name ] = imgptr
	}

  composite := *imgptr


	bounds := composite.Bounds()
	fmt.Println( bounds )
	// e.g.   punks.png  (0,0)-(2400,2400)

	width, height := bounds.Max.X, bounds.Max.Y

	cols, rows  :=   width / tile_width,  height / tile_height

	tile_count := cols * rows


	fmt.Printf( "composite %dx%d (cols x rows) - %d tiles - %dx%d (width x height) \n",
									 cols, rows, tile_count, tile_width, tile_height )
	fmt.Println()


	////////////////////////////////////////////////////
	// note: for now id might optionally include an extension
	//                e.g. .png, .svg etc.

	parts := strings.Split( ctx.Param( "id" ), "." )

  id, _ := strconv.Atoi( parts[0] )

	format := "png"     // default format to png
	if len(parts) > 1 {
      format = parts[1]
	}





	y, x := divmod( id, cols)
	fmt.Printf( "  #%d - tile @ x/y %d/%d... ", id, x, y )


	//
	// todo/fix: change to newNRGBA (better match for png - why? why not?)
	tile := image.NewRGBA( image.Rect(0,0, tile_width, tile_height) )


  if format == "svg" {

// sp (starting point) in composite
sp    := image.Point{ x*tile_width, y*tile_height}
draw.Draw( tile, tile.Bounds(), composite, sp, draw.Over )
//  draw.Src )   // draw.Over )
// note: was draw.Src


mirrorParam := ctx.DefaultQuery( "mirror", "0" )
if mirrorParam == "0" || mirrorParam[0] == 'f' || mirrorParam[0] == 'n'  {
	mirrorParam = ctx.DefaultQuery( "m", "0" )  // allow shortcut m too
}

if mirrorParam == "1" || mirrorParam[0] =='t' || mirrorParam[0] =='y' {
	tile, _ = mirrorImage( tile )
}

  buf :=  imageToSVG( tile )


	saveParam := ctx.DefaultQuery( "save", "0" )
	if saveParam == "0" || saveParam[0] == 'f' || saveParam[0] == 'n'  {
		saveParam = ctx.DefaultQuery( "s", "0" )  // allow shortcut s too
	}


  if saveParam == "1" || saveParam[0] =='t' || saveParam[0] =='y' {

    basename := fmt.Sprintf( "%s-%06d", name, id )

	  if mirrorParam == "1" || mirrorParam[0] == 't' || mirrorParam[0] == 'y' {
			basename = fmt.Sprintf( "%s_mirror", basename )
		}

		outpath := fmt.Sprintf( "./%s.svg", basename )

		fmt.Printf( "  saving image to >%s<...\n", outpath )

		fout, err := os.Create( outpath )
		if err != nil {
			log.Fatal(err)
		}
		defer fout.Close()

		fout.WriteString( buf )
	}


   fmt.Printf( "%s-%d.svg %dx%d image - %d byte(s)\n", name, id,
							 tile_width, tile_height,
							 len( buf ))

ctx.Data( http.StatusOK, ContentTypeImageSVG,  []byte( buf ) )






	} else {   // assume "png" format

	backgroundParam := ctx.DefaultQuery( "background", "" )
	if backgroundParam == "" {
		backgroundParam = ctx.DefaultQuery( "bg", "" )  // allow shortcut z too
	}

	if backgroundParam != "" {
		 fmt.Printf( "=> parsing background color (in hex) >%s<...\n", backgroundParam )

     background, err := parseColor( backgroundParam )
     if err != nil {
			 // todo/fix:  only report parse color error and continue? why? why not?
			 log.Panic( err )
		 }

	  /// use Image.ZP for image.Point{0,0} - why? why not?
	   draw.Draw( tile, tile.Bounds(), &image.Uniform{ background}, image.Point{0,0}, draw.Src )
	}

	// sp (starting point) in composite
	sp    := image.Point{ x*tile_width, y*tile_height}
	draw.Draw( tile, tile.Bounds(), composite, sp, draw.Over )
	//  draw.Src )   // draw.Over )
	// note: was draw.Src



	mirrorParam := ctx.DefaultQuery( "mirror", "0" )
	if mirrorParam == "0" || mirrorParam[0] == 'f' || mirrorParam[0] == 'n'  {
		mirrorParam = ctx.DefaultQuery( "m", "0" )  // allow shortcut m too
	}

  if mirrorParam == "1" || mirrorParam[0] =='t' || mirrorParam[0] =='y' {
		tile, _ = mirrorImage( tile )
	}




	zoomParam := ctx.DefaultQuery( "zoom", "0" )
	if zoomParam == "0" {
		zoomParam = ctx.DefaultQuery( "z", "1" )  // allow shortcut z too
	}

	zoom, _ := strconv.Atoi( zoomParam )

  if zoom > 1 {
		fmt.Printf( " %dx zooming...\n", zoom )
		tile, _ = zoomImage( tile, zoom )
	}


	saveParam := ctx.DefaultQuery( "save", "0" )
	if saveParam == "0" || saveParam[0] == 'f' || saveParam[0] == 'n'  {
		saveParam = ctx.DefaultQuery( "s", "0" )  // allow shortcut s too
	}

  if saveParam == "1" || saveParam[0] =='t' || saveParam[0] =='y' {

    basename := fmt.Sprintf( "%s-%06d", name, id )

		if zoom > 1 {
      basename = fmt.Sprintf( "%s@%dx", basename, zoom )
		}
    if mirrorParam == "1" || mirrorParam[0] == 't' || mirrorParam[0] == 'y' {
			basename = fmt.Sprintf( "%s_mirror", basename )
		}
		if backgroundParam != "" {
			basename = fmt.Sprintf( "%s_(%s)", basename, backgroundParam )
		}

		outpath := fmt.Sprintf( "./%s.png", basename )

		fmt.Printf( "  saving image to >%s<...\n", outpath )

		fout, err := os.Create( outpath )
		if err != nil {
			log.Fatal(err)
		}
		defer fout.Close()

		png.Encode( fout, tile )
	}



	buf := new( bytes.Buffer )
	_ = png.Encode( buf, tile )

	bytesTile := buf.Bytes()
  fmt.Printf( "%s-%d@%dx png %dx%d image - %d byte(s)\n", name, id, zoom,
	               tile_width, tile_height,
	               len( bytesTile ))

	ctx.Data( http.StatusOK, ContentTypeImagePNG, bytesTile )
  }
  }
}



func main() {
	fmt.Printf( "%d collection(s):\n", len( collections ))
	fmt.Println( collections )

  fmt.Println( "cache:" )
  fmt.Println( cache )

	// check if make is required for setup to avoid crash / panic!!!
	cache = make( map[string]*image.Image )

  fmt.Println( "cache:" )
	fmt.Println( cache )


	compileTemplates()


	router := gin.Default()

	router.GET( "/",  handleHome( collections ) )

	for i,c := range collections {
		fmt.Printf( "  [%d] %s  %dx%d - %s\n", i, c.Name, c.Width, c.Height, c.Path )

		router.GET( "/" + c.Name,  handleCollection( c ) )

		// note - &c will NOT work - as c as reference gets
		//          all handlers pointing to last collection!!!!
		router.GET( "/" + c.Name + "/:id", handleCollectionImage( c ) )
	}

	router.Run( "localhost:8080" )


	fmt.Println( "Bye!")
}

