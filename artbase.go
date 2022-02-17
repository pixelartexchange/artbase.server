package main

import (
	"fmt"
	"image"
	"image/png"
	"image/draw"
	"os"
	"log"
	"errors"
  "strings"
	"strconv"
	"bytes"
	"net/http"
	"github.com/gin-gonic/gin"

	"./artbase"
	"./pixelart"
)





// check if divod exists built-in - different name or such ??
func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient  = numerator / denominator   // integer division, decimals are truncated
	remainder = numerator % denominator
	return
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




const (
	ContentTypeJSON     = "application/json"
	ContentTypeHTML     = "text/html; charset=utf-8"
	ContentTypeText     = "text/plain; charset=utf-8"
	ContentTypeImagePNG = "image/png"
	ContentTypePNG      = ContentTypeImagePNG
	ContentTypeImageSVG = "image/svg+xml"
	ContentTypeSVG      = ContentTypeImageSVG
)



var cache map[string]image.Image


func handleHome( collections []artbase.Collection ) gin.HandlerFunc  {
  return func( ctx *gin.Context ) {
		 b := renderHome( collections )
		 ctx.Data( http.StatusOK, ContentTypeHTML, b )
	}
}

func handleCollection( col artbase.Collection ) gin.HandlerFunc  {
  return func( ctx *gin.Context ) {
		 b := renderCollection( &col )
		 ctx.Data( http.StatusOK, ContentTypeHTML, b )
	}
}



func handleCollectionImage( col artbase.Collection ) gin.HandlerFunc  {
	// check if collection is in cache
  return func( ctx *gin.Context ) {

		path := col.Path

	if exist, _ := fileExist( path ); !exist  {
    fmt.Println( "  getting composite / download to (local) cache..." )

		url := col.Url
		pixelart.Download( url, path )
	}


	name                    := col.Name
	tile_width, tile_height := col.Width, col.Height   // in px


  var img image.Image


	if cache[ name ] != nil {
		 fmt.Println( "    bingo! (re)using in-memory composite image from cache...\n" )
     img = cache[ name ]
	} else {
		 fmt.Println( "   adding composite image to in-memory cache...\n" )
	   img = pixelart.ReadImagePNG( path )
		 cache[ name ] = img
	}

  composite := img


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
	tile, _ = pixelart.MirrorImage( tile )
}

  buf :=  pixelart.ImageToSVG( tile )


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

     background, err := pixelart.ParseColor( backgroundParam )
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
		tile, _ = pixelart.MirrorImage( tile )
	}




	zoomParam := ctx.DefaultQuery( "zoom", "0" )
	if zoomParam == "0" {
		zoomParam = ctx.DefaultQuery( "z", "1" )  // allow shortcut z too
	}

	zoom, _ := strconv.Atoi( zoomParam )

  if zoom > 1 {
		fmt.Printf( " %dx zooming...\n", zoom )
		tile, _ = pixelart.ZoomImage( tile, zoom )
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
	cache = make( map[string]image.Image )

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

