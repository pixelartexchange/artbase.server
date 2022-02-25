package main

import (
	"fmt"
	// "image"
	// "image/png"
	// "image/color"
	"log"
  // "strings"
	// "bytes"
	"net/http"
	// "strconv"

	"./artbase"
	"./pixelart"
)



func handleHome( collections []artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {
		 b := artbase.RenderHome( collections )

		 w.Header().Set( ContentType, ContentTypeHTML )
		 w.Write( b )
		}
}

func handleCollection( col artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {
		 b := artbase.RenderCollection( &col )

		 w.Header().Set( ContentType, ContentTypeHTML )
		 w.Write( b )
	}
}



func handleCollectionImagePNG( col artbase.Collection ) http.HandlerFunc  {
		return func( w http.ResponseWriter, req *http.Request ) {

	id, _ := URLParamInt( req, "id" )

	opts := artbase.PNGOpts{}



	backgroundQuery, ok := URLQuery( req, "background" )
	if !ok {
		backgroundQuery, ok = URLQuery( req, "bg" )  // allow shortcut z too
	}

	if ok {
		 fmt.Printf( "=> parsing background color (in hex) >%s<...\n", backgroundQuery )

     background, err := pixelart.ParseColor( backgroundQuery )
     if err != nil {
			 log.Panic( err )
		 }

		 opts.Background     = background
		 opts.BackgroundName = backgroundQuery
	}

	mirror, ok := URLQueryBool( req,  "mirror" )
	if !ok {
		mirror, ok = URLQueryBool( req,  "m" )  // allow shortcut m too
	}

  if mirror {
		opts.Mirror = true
	}


	zoom, ok := URLQueryInt( req,  "zoom" )
	if !ok {
		zoom, ok = URLQueryInt( req,  "z"  )  // allow shortcut z too
	}


  if zoom > 1 {
		opts.Zoom = zoom
	}


	save, ok := URLQueryBool( req,  "save" )
	if !ok {
		save, ok = URLQueryBool( req,  "s" )  // allow shortcut s too
	}

  if save {
    opts.Save = true
	}

	b := col.HandleTilePNG( id, opts )

	w.Header().Set( ContentType, ContentTypeImagePNG )
	w.Write( b )
  }
}


func handleCollectionImageSVG( col artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {

  id, _ := URLParamInt( req, "id" )

		opts := artbase.SVGOpts{}

   mirror, ok := URLQueryBool( req, "mirror" )
   if !ok {
	   mirror, ok = URLQueryBool( req, "m" )  // allow shortcut m too
   }

   if mirror {
	   opts.Mirror = true
   }

	 save, ok := URLQueryBool( req, "save" )
	 if !ok {
		 save, ok = URLQueryBool( req, "s" )  // allow shortcut s too
	 }

  if save {
    opts.Save = true
  }

	b := col.HandleTileSVG( id, opts )

	w.Header().Set( ContentType, ContentTypeImageSVG )
	w.Write( b )
  }
}



func main() {

	//// note:
	// use built-in "standard" collections for now,
	//   yes, you can - use / set-up your own collections
	collections := artbase.Collections

	fmt.Printf( "%d collection(s):\n", len( collections ))
	fmt.Println( collections )


	var router Router

	router.GET( "/",  handleHome( collections ) )

	for i,c := range collections {
		fmt.Printf( "  [%d] %s  %dx%d - %s\n", i, c.Name, c.Width, c.Height, c.Path )

		router.GET( "/" + c.Name,  handleCollection( c ) )

		// note - &c will NOT work - as c as reference gets
		//          all handlers pointing to last collection!!!!
		router.GET( "/" + c.Name + `/(?P<id>[0-9]+)(\.png)?`, handleCollectionImagePNG( c ) )
		router.GET( "/" + c.Name + `/(?P<id>[0-9]+)\.svg`,    handleCollectionImageSVG( c ) )
	}

	http.ListenAndServe( "localhost:8080", &router )

	fmt.Println( "Bye!")
}

