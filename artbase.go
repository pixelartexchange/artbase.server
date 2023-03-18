package main

import (
	"fmt"
	// "log"
	"os"
	"net/http"


	// "github.com/pixelartexchange/artbase.server/artbase"
	"github.com/pixelartexchange/artbase.server/artbase/collections"
	"github.com/pixelartexchange/artbase.server/serve"
)



func main() {

	//// note:
	// use built-in "standard" collections for now,
	//   yes, you can - use / set-up your own collections

	 collections := collections.Standard
	 // collections := collections.Ordinals

	serve := serve.NewRouter( collections )



	// default addr to localhost:8080 for now
	//    for windows include localhost to avoid firewall warning/popup
	//       if binding to :8080 only  - why? why not?

  addr := "localhost:8080"

	// check for port in env settings - required by heroku
	port := os.Getenv( "PORT" )
	if port != "" {
		addr = ":" + port
	}


 	http.ListenAndServe( addr, serve )

	fmt.Println( "Bye!")
}

