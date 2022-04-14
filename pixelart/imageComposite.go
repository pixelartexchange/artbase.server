package pixelart

import (
	"fmt"
	"image"
	"image/draw"
)


//
// for convenience add a type alias for image.Point
type Point = image.Point



// note: for now ALWAYS use
// image of type *image.NRGBA
//   (32bit RGBA colors, not pre-multiplied by alpha)
//  internally!!!
//
// type NRGBA struct {
//	Pix []uint8
//	Stride int
//	Rect Rectangle
// }

// note: use "composition" - for "nicer" api - why? why not?
//    type alias is not working e.g. type Image = image.NRGBA
//      results in
//       cannot define new methods on non-local type *image.NRGBA
//       and such -  can be fixed? why? why not?

type Image struct {
	*image.NRGBA    // use "composition" - for "nicer"  api (lets you use like image.Image) - why? why not?
}

/*
 see https://stackoverflow.com/questions/61964247/idiomatic-go-for-handling-any-image-type-with-any-color-type
     https://stackoverflow.com/questions/47535474/convert-image-from-image-ycbcr-to-image-rgba
*/
func ConvertToNRGBA( img image.Image ) *image.NRGBA {
    var ok bool
		var nrgba *image.NRGBA

		nrgba, ok = img.(*image.NRGBA)
    if ok {
        return nrgba
    }

    fmt.Printf( "  auto-converting image %v of type %T to *image.NRGBA\n", img.Bounds(), img )

    bounds := img.Bounds()
    nrgba = image.NewNRGBA( image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
    draw.Draw( nrgba, nrgba.Bounds(), img, bounds.Min, draw.Src )

    return nrgba
}


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

	return &ImageComposite{ Image: *img,
		                      TileWidth:  tileSize.X,
		                      TileHeight: tileSize.Y }
}


func NewImageComposite( cols int, rows int,
	                      tileSize *image.Point ) *ImageComposite {

  img := NewImage( cols*tileSize.X, rows*tileSize.Y)

  return &ImageComposite{ Image: *img,    // note: use * for dereference - NewImage returns pointer to Image!!
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


func (composite *ImageComposite) Tile( id int ) *Image {
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

  tile := NewImage( composite.TileWidth, composite.TileHeight )

	 // sp (starting point) in composite
	 sp  := image.Point{ x*composite.TileWidth, y*composite.TileHeight }
	 draw.Draw( tile, tile.Bounds(), composite, sp, draw.Over )

	return tile
}



func (composite *ImageComposite) Add( tile image.Image ) {
	 bounds := composite.Bounds()
	 width, height := bounds.Dx(), bounds.Dy()

	 cols, rows :=  width / composite.TileWidth,
	                height / composite.TileHeight
	 y, x :=  divmod( composite.Count, cols )

   fmt.Printf( "adding tile %d @ x/y %d/%d in (%dx%d)...\n", composite.Count, x, y, cols, rows )

	// sp (starting point) in composite
	sp  := image.Point{ x*composite.TileWidth, y*composite.TileHeight }
	 // fmt.Println( sp )

	rect := image.Rect( sp.X, sp.Y,
		                  sp.X+tile.Bounds().Dx(),
											sp.Y+tile.Bounds().Dy())
	draw.Draw( composite,
	           rect, tile, image.Point{0,0}, draw.Over )

	composite.Count += 1
}



func (composite *ImageComposite) Strip() *Image {
	count := 9     // make count into an optional parameter - why? why not?
	// note: if count is 9 use a 9x1 grid and so on
  strip := NewImageComposite( count, 1, &image.Point{composite.TileWidth,
		                                                 composite.TileHeight})

  // note:  check max count
	safe_count := min( count, composite.Max() )

  for i:=0; i<safe_count; i++ {
     strip.Add( composite.Tile( i ) )
	}

	return &strip.Image
}

