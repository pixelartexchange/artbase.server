package artbase

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"errors"

	"../pixelart"   // todo/check if relative to "root" or package ???
)



func fileExist(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
			return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
			return false, nil
	}
	return false, err
}


/////
//  todo/check -  make (local) cache public - why? why not?
//                              or just keep as "internal" detail
//  note: remember map always requires make or map literal to init/setup
var Cache = make( map[string]image.Image )


///////////////////////////////
//  get image - if not present in (local) cache - (auto-)download first!!!

func (col *Collection) Image() image.Image  {

	path := col.Path

	if exist, _ := fileExist( path ); !exist  {
    fmt.Println( "    [artbase-cache] getting composite / download to (local) cache..." )

		pixelart.Download( col.Url, path )
	}

	name := col.Name

	var img image.Image

  if Cache[ name ] != nil {
	   fmt.Println( "    [artbase-cache] bingo! (re)using in-memory composite image from cache...\n" )
	   img = Cache[ name ]
  } else {
	   fmt.Println( "    [artbase-cache] adding composite image to in-memory cache...\n" )
	   img = pixelart.ReadImagePNG( path )
	   Cache[ name ] = img
  }

	return img
}




// check if divod exists built-in - different name or such ??
func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient  = numerator / denominator   // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}


/// todo/check:
///   move tile to pixelart - why? why not?
func (col *Collection) Tile( id int, background color.Color ) *image.RGBA  {

	composite := col.Image()   // get image via (built-in) cache

	bounds := composite.Bounds()
	fmt.Println( bounds )
	// e.g.   punks.png  (0,0)-(2400,2400)
	width, height := bounds.Max.X, bounds.Max.Y  // todo/check: use bounds.Dx(), bounds.Dy() ???


	tile_width, tile_height := col.Width, col.Height   // in px

	cols, rows  :=   width / tile_width,  height / tile_height

	tile_count := cols * rows


	fmt.Printf( "composite %dx%d (cols x rows) - %d tiles - %dx%d (width x height) \n",
									 cols, rows, tile_count, tile_width, tile_height )
	fmt.Println()


	y, x := divmod( id, cols )
	fmt.Printf( "  #%d - tile @ x/y %d/%d... ", id, x, y )

	//
	// todo/fix: change to newNRGBA (better match for png - why? why not?)
	tile := image.NewRGBA( image.Rect(0,0, tile_width, tile_height) )

	if background != nil {
	  /// use Image.ZP for image.Point{0,0} - why? why not?
		draw.Draw( tile, tile.Bounds(), &image.Uniform{ background }, image.Point{0,0}, draw.Src )
	}

	 // sp (starting point) in composite
	 sp    := image.Point{ x*tile_width, y*tile_height }
	 draw.Draw( tile, tile.Bounds(), composite, sp, draw.Over )

	return tile
}






