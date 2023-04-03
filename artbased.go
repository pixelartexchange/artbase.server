package main

import (
	"fmt"
	// "log"
	"os"
	"net/http"
	"strings"
	"strconv"
	"flag"


	"github.com/pixelartexchange/artbase.server/artbase"
	"github.com/pixelartexchange/artbase.server/serve"
)



func main() {

	//// note:
	// download "standard" collections for now,
	//   yes, you can - use / set-up your own collections

	var configValue string
	configDef := "https://github.com/pixelartexchange/artbase.server/raw/master/collections.csv"
    // configDef := "https://github.com/pixelartexchange/ordbase.server/raw/master/collections.csv"
	  configUsage := "path or url to collections dataset"

    flag.StringVar(&configValue, "config", configDef, configUsage)
    flag.StringVar(&configValue, "c",      configDef, configUsage+" (shorthand)")


    var portValue int
    portDef   := 8080
    portUsage := "port"

    // check for port in env settings - required by heroku et al
    portEnv := os.Getenv( "PORT" )
    if portEnv != "" {
       if i, err := strconv.Atoi( portEnv ); err == nil {
           portDef = i   // change default to env variable; if error ignore for now
       }
    }

    flag.IntVar(&portValue, "port", portDef, portUsage )
    flag.IntVar(&portValue, "p",    portDef, portUsage+" (shorthand)")


    flag.Parse()

		 var collections []*artbase.Collection

    if strings.HasPrefix( configValue, "http://" ) ||
       strings.HasPrefix( configValue, "https://" ) {
				collections = artbase.DownloadCollections( configValue )
    } else {
    		collections = artbase.ReadCollections( configValue )
    }


	  serve := serve.NewRouter( collections )


		// default addr to localhost:8080 for now
	  //    for windows include localhost to avoid firewall warning/popup
	  //       if binding to :8080 only  - why? why not?

		addr := "localhost:" + strconv.Itoa( portValue )

 	  http.ListenAndServe( addr, serve )

	  fmt.Println( "Bye!" )
}

