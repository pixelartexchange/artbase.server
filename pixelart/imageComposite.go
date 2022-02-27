package pixelart

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"log"
)



type ImageComposite struct {
	image.Image     // use "composition"
	TileWidth, TileHeight int
	Count int     // optional - not all tiles (full cap(acity) might be used)
}

type ImageTile struct {
	image.Image    // use "composition" - for "nicer" api - why? why not?
}



// check if divod exists built-in - different name or such ??
func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient  = numerator / denominator   // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}



func (composite *ImageComposite) Tile( id int, background color.Color ) *ImageTile {
	bounds := composite.Bounds()
	fmt.Println( bounds )
	// e.g.   punks.png  (0,0)-(2400,2400)
	width, height := bounds.Max.X, bounds.Max.Y  // todo/check: use bounds.Dx(), bounds.Dy() ???


	cols, rows :=  width / composite.TileWidth,  height / composite.TileHeight

	tile_count := cols * rows


	fmt.Printf( "composite %dx%d (cols x rows) - %d tiles - %dx%d (width x height) \n",
									 cols, rows, tile_count, composite.TileWidth, composite.TileHeight )
	fmt.Println()


	y, x := divmod( id, cols )
	fmt.Printf( "  #%d - tile @ x/y %d/%d... ", id, x, y )

	//
	// todo/fix: change to newNRGBA (better match for png - why? why not?)
	tile := image.NewRGBA( image.Rect(0,0, composite.TileWidth, composite.TileHeight) )

	if background != nil {
	  /// use Image.ZP for image.Point{0,0} - why? why not?
		draw.Draw( tile, tile.Bounds(), &image.Uniform{ background }, image.Point{0,0}, draw.Src )
	}

	 // sp (starting point) in composite
	 sp    := image.Point{ x*composite.TileWidth, y*composite.TileHeight }
	 draw.Draw( tile, tile.Bounds(), composite, sp, draw.Over )

	return &ImageTile{ Image: tile }
}



////////////////////////////////
// tile methods  "convenience helpers" for easy chaining

func (tile *ImageTile) Zoom( zoom int ) *ImageTile {
    img, _ := ZoomImage( tile.Image, zoom )
		return &ImageTile{ Image: img }
}

func (tile *ImageTile) Mirror() *ImageTile {
	img, _ := MirrorImage( tile.Image )
	return &ImageTile{ Image: img }
}


func (tile *ImageTile) Save( path string ) {

	fmt.Printf( "  saving image to >%s<...\n", path )

  // todo/check - auto-create directories in path - why? why not?
	fout, err := os.Create( path )
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	png.Encode( fout, tile )
}

