package main

import (
	"fmt"
	// "image"
	// "image/png"
	"image/color"
	"log"
  "strings"
	"strconv"
	// "bytes"
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

		opts := artbase.ImageSVGOpts{}

mirrorParam := ctx.DefaultQuery( "mirror", "0" )
if mirrorParam == "0" || mirrorParam[0] == 'f' || mirrorParam[0] == 'n'  {
	mirrorParam = ctx.DefaultQuery( "m", "0" )  // allow shortcut m too
}

if mirrorParam == "1" || mirrorParam[0] =='t' || mirrorParam[0] =='y' {
	  opts.Mirror = true
}

	saveParam := ctx.DefaultQuery( "save", "0" )
	if saveParam == "0" || saveParam[0] == 'f' || saveParam[0] == 'n'  {
		saveParam = ctx.DefaultQuery( "s", "0" )  // allow shortcut s too
	}

  if saveParam == "1" || saveParam[0] =='t' || saveParam[0] =='y' {
    opts.Save = true
  }


	bytes := artbase.HandleCollectionImageSVG( &col, id, opts )

  ctx.Data( http.StatusOK, ContentTypeImageSVG, bytes )

	} else {   // assume "png" format


  var background color.Color = nil   // interface by default (zero-value) nil??
  var err error              = nil   // interface by default (zero-value) nil??


	opts := artbase.ImageOpts{}

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

		 opts.Background     = background
		 opts.BackgroundName = backgroundParam
	}

	mirrorParam := ctx.DefaultQuery( "mirror", "0" )
	if mirrorParam == "0" || mirrorParam[0] == 'f' || mirrorParam[0] == 'n'  {
		mirrorParam = ctx.DefaultQuery( "m", "0" )  // allow shortcut m too
	}

  if mirrorParam == "1" || mirrorParam[0] =='t' || mirrorParam[0] =='y' {
		opts.Mirror = true
	}


	zoomParam := ctx.DefaultQuery( "zoom", "0" )
	if zoomParam == "0" {
		zoomParam = ctx.DefaultQuery( "z", "1" )  // allow shortcut z too
	}

	zoom, _ := strconv.Atoi( zoomParam )

  if zoom > 1 {
		opts.Zoom = zoom
	}


	saveParam := ctx.DefaultQuery( "save", "0" )
	if saveParam == "0" || saveParam[0] == 'f' || saveParam[0] == 'n'  {
		saveParam = ctx.DefaultQuery( "s", "0" )  // allow shortcut s too
	}

  if saveParam == "1" || saveParam[0] =='t' || saveParam[0] =='y' {
    opts.Save = true
	}


	bytes := artbase.HandleCollectionImage( &col, id, opts )


	ctx.Data( http.StatusOK, ContentTypeImagePNG, bytes )
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

