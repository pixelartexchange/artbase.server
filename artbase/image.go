package artbase


import (
	"log"
	"os"
	"fmt"
	"image/color"
	"image/png"
	"bytes"

	"github.com/pixelartexchange/artbase.server/pixelart"   // todo/check if relative to "root" or package ???
)


///
// check - change ImageOpts & HandleImage to
//                TileOpts & HandleTile  - why? why not?


type PNGOpts struct {
	 Background color.Color  // default: nil
	 BackgroundName string   // default: ""
	 Mirror bool             // default: false
	 Zoom int                // default: FIX??? use 1 NOT 0 - how?
	 Save bool                // default: false
}


func (col *Collection) HandleTilePNG( id int,
										                opts PNGOpts )  []byte  {


  tile := col.Image().Tile( id )

	if opts.Background != nil {
    tile = tile.Background( opts.Background )
	}

	if opts.Mirror {
		tile = tile.Mirror()
	}

  if opts.Zoom > 1 {
		fmt.Printf( " %dx zooming...\n", opts.Zoom )
		tile = tile.Zoom( opts.Zoom )
	}


	name := col.Name

	if opts.Save {
	  basename := fmt.Sprintf( "%s-%06d", name, id )

		if opts.Zoom > 1 {
      basename = fmt.Sprintf( "%s@%dx", basename, opts.Zoom )
		}
    if opts.Mirror {
			basename = fmt.Sprintf( "%s_mirror", basename )
		}
		if opts.Background != nil {
			 basename = fmt.Sprintf( "%s_(%s)", basename, opts.BackgroundName )
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


	// note: change default (0) to 1
  zoom := opts.Zoom
	if opts.Zoom == 0 {
		zoom = 1
	}


	bytesTile := buf.Bytes()
  fmt.Printf( "%s-%d@%dx png %dx%d image - %d byte(s)\n",
	               name, id,
								 zoom,
	               col.Width, col.Height,
	               len( bytesTile ))

	return bytesTile
}




type SVGOpts struct  {
	Background color.Color  // default: nil
	BackgroundName string   // default: ""
	Mirror bool             // default: false
	Save bool                // default: false
}

func (col *Collection) HandleTileSVG( id int,
														          opts SVGOpts )  []byte  {

  tile := col.Image().Tile( id )

	if opts.Mirror {
		tile = tile.Mirror()
	}

  buf :=  pixelart.ImageToSVG( tile )

	name := col.Name

	if opts.Save {
	  basename := fmt.Sprintf( "%s-%06d", name, id )

    if opts.Mirror {
			basename = fmt.Sprintf( "%s_mirror", basename )
		}
		if opts.Background != nil {
			 basename = fmt.Sprintf( "%s_(%s)", basename, opts.BackgroundName )
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


	   fmt.Printf( "%s-%d.svg %dx%d image - %d byte(s)\n",
		           name, id,
							 col.Width, col.Height,
							 len( buf ))

     return []byte( buf )
}

