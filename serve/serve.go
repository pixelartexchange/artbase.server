package serve

import (
	"fmt"
	// "log"
	// "os"
	"net/http"

	"github.com/learnpixelart/pixelart.go/pixelart"

	"github.com/pixelartexchange/artbase.server/artbase"
	"github.com/pixelartexchange/artbase.server/router"     // simple http router & helpers from scratch (no 3rd party deps) - replace with your own http libs/frameworks
)


func handleHome( collections []artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {
		 b := renderHome( collections )

		 w.Header().Set( "Content-Type", "text/html; charset=utf-8" )
		 w.Write( b )
		}
}

func handleCollection( col artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {
		 b := renderCollection( &col )

		 w.Header().Set( "Content-Type", "text/html; charset=utf-8" )
		 w.Write( b )
	}
}

func handleCollectionStripPNG( col artbase.Collection ) http.HandlerFunc  {
	return func( w http.ResponseWriter, req *http.Request ) {
		b := col.HandleStripPNG()

		w.Header().Set( "Content-Type", "image/png" )
		w.Write( b )
	}
}




func NewRouter( collections []artbase.Collection ) *router.Router {

  //////
	// for debugging and double check on module print version strings
	fmt.Println( "go package versions:" )
	fmt.Println( "  artbase:",  artbase.Version )
	fmt.Println( "  pixelart:", pixelart.Version )
	fmt.Println( "  router:",   router.Version )
  fmt.Println()

	fmt.Printf( "%d collection(s):\n", len( collections ))
	fmt.Println( collections )



	serve := router.Router{}

	serve.GET( "/",  handleHome( collections ) )

	for i, c := range collections {
		fmt.Printf( "  [%d] %s  %dx%d - %s\n", i, c.Name, c.Width, c.Height, c.Path )

		serve.GET( "/" + c.Name,  handleCollection( c ) )
		serve.GET( "/" + c.Name + "-strip.png", handleCollectionStripPNG( c ) )

		// note - &c will NOT work - as c as reference gets
		//          all handlers pointing to last collection!!!!
		serve.GET( "/" + c.Name + `/(?P<id>[0-9]+)(\.png)?`, handleCollectionImagePNG( c ) )
		serve.GET( "/" + c.Name + `/(?P<id>[0-9]+)\.svg`,    handleCollectionImageSVG( c ) )
	}

	fmt.Println( "Bye!" )

	return &serve
}


