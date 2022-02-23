package main

import (
	"fmt"
	// "image"
	"image/png"
	"image/color"
	"os"
	"log"
  "strings"
	"strconv"
	"bytes"
	"net/http"
	"github.com/gin-gonic/gin"

	"./artbase"
	"./pixelart"
)





const (
	ContentTypeJSON     = "application/json"
	ContentTypeHTML     = "text/html; charset=utf-8"
	ContentTypeText     = "text/plain; charset=utf-8"
	ContentTypeImagePNG = "image/png"
	ContentTypePNG      = ContentTypeImagePNG
	ContentTypeImageSVG = "image/svg+xml"
	ContentTypeSVG      = ContentTypeImageSVG
)





func handleHome( collections []artbase.Collection ) gin.HandlerFunc  {
  return func( ctx *gin.Context ) {
		 b := artbase.RenderHome( collections )
		 ctx.Data( http.StatusOK, ContentTypeHTML, b )
	}
}

func handleCollection( col artbase.Collection ) gin.HandlerFunc  {
  return func( ctx *gin.Context ) {
		 b := artbase.RenderCollection( &col )
		 ctx.Data( http.StatusOK, ContentTypeHTML, b )
	}
}



func handleCollectionImage( col artbase.Collection ) gin.HandlerFunc  {
  return func( ctx *gin.Context ) {


		name := col.Name

	////////////////////////////////////////////////////
	// note: for now id might optionally include an extension
	//                e.g. .png, .svg etc.

	parts := strings.Split( ctx.Param( "id" ), "." )

  id, _ := strconv.Atoi( parts[0] )

	format := "png"     // default format to png
	if len(parts) > 1 {
      format = parts[1]
	}



  if format == "svg" {

		tile := col.Tile( id, nil )    // no background (color) - use nil

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
							 col.Width, col.Height,
							 len( buf ))

ctx.Data( http.StatusOK, ContentTypeImageSVG,  []byte( buf ) )


	} else {   // assume "png" format



  var background color.Color = nil   // interface by default (zero-value) nil??
  var err error              = nil   // interface by default (zero-value) nil??

	backgroundParam := ctx.DefaultQuery( "background", "" )
	if backgroundParam == "" {
		backgroundParam = ctx.DefaultQuery( "bg", "" )  // allow shortcut z too
	}

	if backgroundParam != "" {
		 fmt.Printf( "=> parsing background color (in hex) >%s<...\n", backgroundParam )

     background, err = pixelart.ParseColor( backgroundParam )
     if err != nil {
			 // todo/fix:  only report parse color error and continue? why? why not?
			 log.Panic( err )
		 }
	}

  tile := col.Tile( id, background )


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
	               col.Width, col.Height,
	               len( bytesTile ))

	ctx.Data( http.StatusOK, ContentTypeImagePNG, bytesTile )
  }
  }
}



func main() {


	fmt.Printf( "%d collection(s):\n", len( artbase.Collections ))
	fmt.Println( artbase.Collections )

  fmt.Println( "cache:" )
  fmt.Println( artbase.Cache )


  //// note:
	// use built-in "standard" collections for now,
	//   yes, you can - use / set-up your own collections
	collections := artbase.Collections



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

