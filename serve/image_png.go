package serve

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/learnpixelart/pixelart.go/pixelart"

	"github.com/pixelartexchange/artbase.server/artbase"
	"github.com/pixelartexchange/artbase.server/router"     // simple http router & helpers from scratch (no 3rd party deps) - replace with your own http libs/frameworks
)



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

	transparent, ok := router.QueryBool( req,  "transparent" )
  if transparent {
		opts.Transparent = true
	}


	circle, ok := router.QueryBool( req,  "circle" )
  if circle {
		opts.Circle = true
	}


	zoom, ok := router.QueryInt( req,  "zoom" )
	if !ok {
		zoom, ok = router.QueryInt( req,  "z"  )  // allow shortcut z too
	}

  if zoom > 1 {
		opts.Zoom = zoom
	}

	size, ok := router.QueryInt( req,  "size" )
	if !ok {
		size, ok = router.QueryInt( req,  "s"  )  // allow shortcut s too
	}

  if size > 0 {
		opts.Resize = size
	}



	save, ok := router.QueryBool( req,  "autosave" )
	if !ok {
		save, ok = router.QueryBool( req,  "save" )  // allow shortcut save too
	}

  if save {
    opts.Save = true
	}


	b := col.HandleTilePNG( id, opts )

	w.Header().Set( "Content-Type", "image/png" )
	w.Write( b )
  }
}



