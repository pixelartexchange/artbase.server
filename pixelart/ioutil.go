package pixelart


import (
	"fmt"
	"log"
	"io"
	"os"
  "net/http"
	"image"
	"image/png"
)



func ReadImage( path string ) image.Image  {
	fmt.Printf( "==> reading %s...\n", path )

	f, err := os.Open( path )
	if err != nil {
			log.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode( f )
	if err != nil {
			log.Fatal(err)
	}

	bounds := img.Bounds()
	fmt.Println( bounds )
	// e.g.   punks.png  (0,0)-(2400,2400)

	return img
}




func Download( url, outpath string ) {
  fmt.Printf( "==> Downloading %s...\n", url )

	resp, err := http.Get( url )
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// todo/fix:
	//  (double) check for content type - why? why not?

	fmt.Printf( "HTTP %v\n", resp.StatusCode )
  // dump  headers for debugging
  for name, headers := range resp.Header {
	  for _, h := range headers {
		  fmt.Printf( "  %v  :  %v\n", name, h )
	  }
  }

	 f, err := os.Create( outpath )
	 if err != nil {
		 log.Fatal(err)
	 }
	 defer f.Close()



   n, err := io.Copy( f, resp.Body )   // note:  Copy( writer, reader )
	 if err != nil {
		  fmt.Printf( "  !! ERROR - writing %d byte(s) to %s...\n", n, outpath )
		  log.Fatal(err)
	 }
	 fmt.Printf( "  writing %d byte(s) to %s...\n", n, outpath )
}



