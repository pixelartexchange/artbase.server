package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pixelartexchange/artbase.server/artbase"
	"github.com/pixelartexchange/artbase.server/pixelart"

	"github.com/pixelartexchange/artbase.server/router"     // simple http router & helpers from scratch (no 3rd party deps) - replace with your own http libs/frameworks
)



func handleHome( collections []artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {
		 b := artbase.RenderHome( collections )

		 w.Header().Set( "Content-Type", "text/html; charset=utf-8" )
		 w.Write( b )
		}
}

func handleCollection( col artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {
		 b := artbase.RenderCollection( &col )

		 w.Header().Set( "Content-Type", "text/html; charset=utf-8" )
		 w.Write( b )
	}
}



func handleCollectionImagePNG( col artbase.Collection ) http.HandlerFunc  {
		return func( w http.ResponseWriter, req *http.Request ) {

	id, _ := router.ParamInt( req, "id" )

	opts := artbase.PNGOpts{}


	backgroundQuery, ok := router.Query( req, "background" )
	if !ok {
		backgroundQuery, ok = router.Query( req, "bg" )  // allow shortcut z too
	}

	if ok {
		 log.Printf( "=> parsing background color (in hex) >%s<...\n", backgroundQuery )

     background, err := pixelart.ParseColor( backgroundQuery )
     if err != nil {
			 log.Panic( err )
		 }

		 opts.Background     = background
		 opts.BackgroundName = backgroundQuery
	}


	silhouetteQuery, ok := router.Query( req, "silhouette" )
	if ok {
		 log.Printf( "=> parsing silhouette (forground) color (in hex) >%s<...\n", silhouetteQuery )

     silhouette, err := pixelart.ParseColor( silhouetteQuery )
     if err != nil {
			 log.Panic( err )
		 }

		 opts.Silhouette     = silhouette
		 opts.SilhouetteName = silhouetteQuery
	}

	flag, ok := router.Query( req, "flag" )
	if ok {
		 opts.Flag = flag
	}


	mirror, ok := router.QueryBool( req,  "mirror" )
	if !ok {
		mirror, ok = router.QueryBool( req,  "m" )  // allow shortcut m too
	}

  if mirror {
		opts.Mirror = true
	}


	zoom, ok := router.QueryInt( req,  "zoom" )
	if !ok {
		zoom, ok = router.QueryInt( req,  "z"  )  // allow shortcut z too
	}

  if zoom > 1 {
		opts.Zoom = zoom
	}


	save, ok := router.QueryBool( req,  "save" )
	if !ok {
		save, ok = router.QueryBool( req,  "s" )  // allow shortcut s too
	}

  if save {
    opts.Save = true
	}

	b := col.HandleTilePNG( id, opts )

	w.Header().Set( "Content-Type", "image/png" )
	w.Write( b )
  }
}


func handleCollectionImageSVG( col artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {

  id, _ := router.ParamInt( req, "id" )

		opts := artbase.SVGOpts{}

   mirror, ok := router.QueryBool( req, "mirror" )
	 if !ok {
	   mirror, ok = router.QueryBool( req, "m" )  // allow shortcut m too
   }

   if mirror {
	   opts.Mirror = true
   }

	 save, ok := router.QueryBool( req, "save" )
	 if !ok {
		 save, ok = router.QueryBool( req, "s" )  // allow shortcut s too
	 }

  if save {
    opts.Save = true
  }

	b := col.HandleTileSVG( id, opts )

	w.Header().Set( "Content-Type", "image/svg+xml" )
	w.Write( b )
  }
}



func main() {

  //////
	// for debugging and double check on module print version strings
	fmt.Println( "go package versions:" )
	fmt.Println( "  artbase:",  artbase.Version )
	fmt.Println( "  pixelart:", pixelart.Version )
	fmt.Println( "  router:",   router.Version )
  fmt.Println()


	//// note:
	// use built-in "standard" collections for now,
	//   yes, you can - use / set-up your own collections
	collections := artbase.Collections

	fmt.Printf( "%d collection(s):\n", len( collections ))
	fmt.Println( collections )



	serve := router.Router{}

	serve.GET( "/",  handleHome( collections ) )

	for i,c := range collections {
		fmt.Printf( "  [%d] %s  %dx%d - %s\n", i, c.Name, c.Width, c.Height, c.Path )

		serve.GET( "/" + c.Name,  handleCollection( c ) )

		// note - &c will NOT work - as c as reference gets
		//          all handlers pointing to last collection!!!!
		serve.GET( "/" + c.Name + `/(?P<id>[0-9]+)(\.png)?`, handleCollectionImagePNG( c ) )
		serve.GET( "/" + c.Name + `/(?P<id>[0-9]+)\.svg`,    handleCollectionImageSVG( c ) )
	}

	http.ListenAndServe( "localhost:8080", &serve )

	fmt.Println( "Bye!")
}

