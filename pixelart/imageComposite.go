package pixelart

import (
	"fmt"
	"image"
	"image/draw"
)



type Image struct {
	image.Image    // use "composition" - for "nicer" api - why? why not?
}

// check -
//   change ImageTile to just Image
//  and add a type ImageTile = Image  alias - why? why not?
type ImageTile = Image


type ImageComposite struct {
	Image     // use "composition" - note: does NOT use image.Image but our own!!!
	TileWidth, TileHeight int
	Count int     // optional - not all tiles (full cap(acity) might be used)
}



// check if divod exists built-in - different name or such ??
func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient  = numerator / denominator   // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}


func ReadImageComposite( path string, tileSize *image.Point ) *ImageComposite {
	img := ReadImage( path )

	return &ImageComposite{ Image: Image{ img },
		                      TileWidth:  tileSize.X,
		                      TileHeight: tileSize.Y }
}


func (composite *ImageComposite) Max() int {
  // change to TileCountMax or TileCap(acity) or such - why? why not?
	//   note: not all tiles might be "filled-up / in-use / painted"
	//    passed in .Count   should be the real / actual count

	bounds := composite.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	cols, rows :=  width / composite.TileWidth,  height / composite.TileHeight

	tileCount := cols * rows
	return tileCount
}


func (composite *ImageComposite) Tile( id int ) *ImageTile {
	bounds := composite.Bounds()
	// fmt.Println( bounds )
	// e.g.   punks.png  (0,0)-(2400,2400)
	width, height := bounds.Dx(), bounds.Dy()

	cols, rows :=  width / composite.TileWidth,  height / composite.TileHeight

	tileCount := cols * rows


	fmt.Printf( "composite %dx%d (cols x rows) - %d tiles - %dx%d (width x height) \n",
									 cols, rows, tileCount, composite.TileWidth, composite.TileHeight )

	y, x := divmod( id, cols )
	fmt.Printf( "  #%d - tile @ x/y %d/%d\n", id, x, y )

	//
	// todo/fix: change to newNRGBA (better match for png - why? why not?)
	tile := image.NewRGBA( image.Rect(0,0, composite.TileWidth, composite.TileHeight) )

	 // sp (starting point) in composite
	 sp  := image.Point{ x*composite.TileWidth, y*composite.TileHeight }
	 draw.Draw( tile, tile.Bounds(), composite, sp, draw.Over )

	return &ImageTile{ Image: tile }
}


