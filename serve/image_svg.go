package serve

import (
	// "fmt"
	// "log"
	"net/http"

	// "github.com/learnpixelart/pixelart.go/pixelart"

	"github.com/pixelartexchange/artbase.server/artbase"
	"github.com/pixelartexchange/artbase.server/router"     // simple http router & helpers from scratch (no 3rd party deps) - replace with your own http libs/frameworks
)




func handleCollectionImageSVG( col *artbase.Collection ) http.HandlerFunc  {
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

	 save, ok := router.QueryBool( req, "autosave" )
	 if !ok {
		 save, ok = router.QueryBool( req, "save" )  // allow shortcut save too
	 }

  if save {
    opts.Save = true
  }

	b := col.HandleTileSVG( id, opts )

	w.Header().Set( "Content-Type", "image/svg+xml" )
	w.Write( b )
  }
}


