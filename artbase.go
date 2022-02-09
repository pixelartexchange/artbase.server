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

	 {Name: "bwpunks",  Width: 24, Height: 24,
	  Path: "./bwpunks.png",
	  Url: "https://github.com/pixelartexchange/collections/raw/master/bwpunks-24x24.png",
	  Background: true },


	 {Name: "coolcats",  Width: 24, Height: 24,
	 Path: "./coolcats.png",
	 Url: "https://github.com/cryptopunksnotdead/awesome-24px/raw/master/collection/coolcats.png" },


	 {Name: "blockydoge",  Width: 60, Height: 60,
	 Path: "./blockydoge.png",
   Url: "https://github.com/pixelartexchange/collections/raw/master/blockydoge/blockydoge-60x60.png",
	 Background: true },

	 {Name: "dooggies",  Width: 32, Height: 32,
	 Path: "./dooggies.png",
   Url: "https://github.com/pixelartexchange/collections/raw/master/dooggies-32x32.png",
	 Background: true },

	 {Name: "wiener",  Width: 32, Height: 32,
	 Path: "./wiener.png",
   Url: "https://github.com/pixelartexchange/collections/raw/master/wiener-32x32.png",
	 Background: true },
}


const templateHome = `
<h1>{{ len . }} Collections</h1>
<ul>
    {{range .}}
        <li><a href="/{{.Name}}">{{.Name}} ({{ .Width}}x{{ .Height}})</a></li>
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

<h1>{{ .Name}}   ({{ .Width}}x{{ .Height}}) Collection</h1>

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
	Note: The default image size is ({{ .Width}}x{{ .Height}}).
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
		  fmt.Printf( "  %v: %v\n", name, h )
	  }
  }

	 f, err := os.Create( outpath )
	 if err != nil {
		 log.Fatal(err)
	 }
	 defer f.Close()


   // todo/check: use io.Copy( f, resp.Body ) for streaming/ saving - why? why not?
   //
	 // ==> Downloading https://github.com/pixelartexchange/collections/blob/master/dooggies-32x32.png...
   //   writing 173838 byte(s) to ./dooggies.png...
   //   will stop/fial - should be +1.5 MBs

   n, err := io.Copy( f, resp.Body )   // Copy( writer, reader )
	 if err != nil {
		  fmt.Printf( "  !! ERROR - writing %d byte(s) to %s...\n", n, outpath )
		  log.Fatal(err)
	 }
	 fmt.Printf( "  writing %d byte(s) to %s...\n", n, outpath )

	  // b, err := ioutil.ReadAll( resp.Body )
		// if err != nil {
		//	log.Fatal(err)
		// }
 	  // n, err := f.Write( b )

}



func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if strings.HasPrefix(s, "#") {
     s = s[1:]   /// cut-off leading (prefix)
	}

	switch len(s) {
	case 6:
		_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 3:
		_, err = fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
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


  id, _ := strconv.Atoi( ctx.Param( "id" ))


	zoomParam := ctx.DefaultQuery( "zoom", "0" )
	if zoomParam == "0" {
		zoomParam = ctx.DefaultQuery( "z", "1" )  // allow shortcut z too
	}

	zoom, _ := strconv.Atoi( zoomParam )





	y, x := divmod( id, cols)
	fmt.Printf( "  #%d - tile @ x/y %d/%d... ", id, x, y )


	tile := image.NewRGBA( image.Rect(0,0, tile_width, tile_height) )



	backgroundParam := ctx.DefaultQuery( "background", "" )
	if backgroundParam == "" {
		backgroundParam = ctx.DefaultQuery( "bg", "" )  // allow shortcut z too
	}

	if backgroundParam != "" {
		 fmt.Printf( "=> parsing background color (in hex) >%s<...\n", backgroundParam )

     background, err := parseHexColor( backgroundParam )
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

  if zoom > 1 {
		fmt.Printf( " %dx zooming...\n", zoom )
		tile, _ = zoomImage( tile, zoom )
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

