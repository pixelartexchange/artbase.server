package pixelart

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"log"
	"path/filepath"
)



func min(a, b int) int {
	if a < b {
			return a
	} else {
	    return b
  }
}



////////////////////////////////
// tile methods  "convenience helpers" for easy chaining



// todo - find a better name - ensureColor/safeColor or ? - why? why not?
func MakeColor( a interface{} ) color.Color {
	var c color.Color
	switch a := a.(type) {
	  case color.Color:
		  c = a
		case string:
			c, _ = ParseColor( a )
	  default:
			log.Fatal( "[MakeColor] unexpected type %T: %v", a, a )
  }
  return c
}




// todo - move to ImageTitle.go file - why? why not?
func NewImage( width, height int ) *Image {
  img := image.NewNRGBA( image.Rect( 0,0, width, height ))
  return &Image{ img }
}


func (tile *Image) Background( background_any interface{} ) *Image {

   background := MakeColor( background_any )

	// todo/fix: change to newNRGBA (better match for png - why? why not?)
	width, height := tile.Bounds().Dx(), tile.Bounds().Dy()
	img := NewImage( width, height )

	/// use Image.ZP for image.Point{0,0} - why? why not?
	draw.Draw( img, img.Bounds(), &image.Uniform{ background }, image.Point{0,0}, draw.Src )
	draw.Draw( img, img.Bounds(), tile, image.Point{0,0}, draw.Over )

	return img
}


// draw flag of ukraine -- glory to ukraine! fuck (vladimir) putin! stop the war!
func (tile *Image) Ukraine() *Image {

	blue   := color.NRGBA{0x00, 0x57, 0xb7, 0xff}  // rgb( 0, 87, 183 )
	yellow := color.NRGBA{0xff, 0xdd, 0x00, 0xff}  // rgb( 255, 221, 0)

 // todo/fix: change to newNRGBA (better match for png - why? why not?)
 width, height := tile.Bounds().Dx(), tile.Bounds().Dy()
 img := NewImage( width, height )

 /// use Image.ZP for image.Point{0,0} - why? why not?
 draw.Draw( img, image.Rect( 0,0,width, height/2 ),
	               &image.Uniform{ blue }, image.Point{0,0}, draw.Src )
 draw.Draw( img, image.Rect( 0,height/2,width,height ),
								 &image.Uniform{ yellow }, image.Point{0,0}, draw.Src )

 draw.Draw( img, img.Bounds(), tile, image.Point{0,0}, draw.Over )

 return img
}



func (tile *Image) Silhouette( foreground_any interface{} ) *Image {

	foreground := MakeColor( foreground_any )

	transparent := color.NRGBA{ R: 0,
                           		G: 0,
		                          B: 0,
		                          A: 0 }

 bounds        := tile.Bounds()
 width, height := bounds.Dx(), bounds.Dy()
 img := NewImage( width, height )

 for y := 0; y < height; y++ {
	for x := 0; x < width; x++ {
		pixel := color.NRGBAModel.Convert( tile.At( bounds.Min.X+x,
			                                          bounds.Min.Y+y )).(color.NRGBA)

	  if pixel == transparent {
		   img.Set( bounds.Min.X+x,
				        bounds.Min.Y+y,
							  pixel )
		}  else {
		    img.Set( bounds.Min.X+x,
					       bounds.Min.Y+y,
								 foreground )
		}
	}
}

 return img
}



func (tile *Image) Transparent() *Image {

	bounds      := tile.Bounds()
	background := tile.At( bounds.Min.X,
		                     bounds.Min.Y ) // 0,0

  transparent := color.NRGBA{ R: 0,
														  G: 0,
													    B: 0,
													    A: 0 }

 width, height := bounds.Dx(), bounds.Dy()
 img := NewImage( width, height )

 ///
 //   todo/fix:
 //     change to "fill" algo for now!!!
 //       if algo hits a non-background color it stops

 for x := 0; x < width; x++ {
 	for y := 0; y < height; y++ {
				pixel := tile.At( bounds.Min.X+x,
			                  bounds.Min.Y+y )

	         if pixel == background {
		        img.Set( bounds.Min.X+x,
      				        bounds.Min.Y+y,
			      				  transparent )
					} else {
		        img.Set( bounds.Min.X+x,
						         bounds.Min.Y+y,
							       pixel )
				 }
		}
	}

 return img
}



func (tile *Image) Circle() *Image {

	bounds := tile.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// for radius use min of width / height
	//   if overflow_x  - center (add padding/overflow left & right)
	//   if overflow_y  - anchor to bottom
	min_square := min( width, height )

  overflow_x := width  - min_square
  overflow_y := height - min_square

  fmt.Printf( "   circle %v to %d +%dpx x %d +%dpx overflow\n",
	                  bounds, min_square, overflow_x,
										        min_square, overflow_y )

	r      := min_square / 2
	center := min_square / 2


	// use color.Alpha{0} - why? why not?
	transparent := color.NRGBA{ R: 0,
															G: 0,
															B: 0,
															A: 0 }

	img := NewImage( width, height )

	for x := 0; x < min_square; x++ {
		for y := 0; y < min_square; y++ {
				pixel := tile.At( bounds.Min.X+x+(overflow_x/2),
													bounds.Min.Y+y+overflow_y )

	 xx, yy, rr := float64( x - center )+0.5,
								 float64( y - center )+0.5,
								 float64( r )

					 if xx*xx+yy*yy < rr*rr {
						 img.Set( x+(overflow_x/2),
							        y+overflow_y,
											pixel )
					 } else {
						 img.Set( x+(overflow_x/2),
							        y+overflow_y,
											transparent )
					 }
				 }
	 }
	 return img
}




func (tile *Image) Zoom( zoom int ) *Image {
			bounds := tile.Bounds()
			width, height := bounds.Dx(), bounds.Dy()  // note: same as bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

			// fmt.Println( bounds, width, height )
			// e.g.   punk #0    (0,0)-(24,24)
			//             #561  (1464,120)-(1488,144)
			//             #3100 (0,744)-(24,768)
			//             #7804 (96,1872)-(120,1896)
			//             #8857 (1368,2112)-(1392,2136)

			img := NewImage( width*zoom, height*zoom )

			for x:=0; x < width; x++ {
				for y:=0; y < height; y++ {
						pixel := tile.At( bounds.Min.X+x, bounds.Min.Y+y )
						for n:=0; n < zoom; n++ {
							for m:=0; m < zoom; m++ {
								img.Set( n+zoom*x, m+zoom*y, pixel )
							}
						}
			 }
		}

		return img
}


func (tile *Image) Mirror() *Image {
		bounds := tile.Bounds()
		width, height := bounds.Dx(), bounds.Dy()  // bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

		// fmt.Println( bounds, width, height )
		// e.g.   punk #0    (0,0)-(24,24)
		//             #561  (1464,120)-(1488,144)
		//             #3100 (0,744)-(24,768)
		//             #7804 (96,1872)-(120,1896)
		//             #8857 (1368,2112)-(1392,2136)

		img := NewImage( width, height )

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				pixel := tile.At( bounds.Min.X+x, bounds.Min.Y+y )
				img.Set( (width-1)-x, y, pixel )
			}
		}

		return img
}



func (tile *Image) Paste( img image.Image ) {
	// note - image.Image (is read-only - no Set() method)
	//          convert/typecast to draw.Image (that includes Set() method)
  //  see https://stackoverflow.com/questions/36573413/change-color-of-a-single-pixel-golang-image

	draw.Draw( tile,
	           img.Bounds(), img, image.Point{0,0}, draw.Over )
}


func (tile *Image) Resize( width int ) *Image {
  // note: for now resize only width
	//   note: for now assumes width is always greater than actual width
	//              todo -assert width is greater (otherwise report error)
	bounds := tile.Bounds()
	image_width, image_height := bounds.Dx(), bounds.Dy()  // bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

  zoom_x, overflow_x := divmod( width, image_width )

	// resize height to same proportion as width
	height := image_height * width / image_width
  fmt.Printf( "   resize %v to %d x %d\n",   bounds, width, height )

  zoom_y, overflow_y := divmod( height, image_height )
	fmt.Printf( "     using zoom (x) %dx +%dpx, (y) %dx +%dpx overflow\n",
	                        zoom_x, overflow_x,
												  zoom_y, overflow_y)

  base := tile

  if zoom_x > 1 {
		base = tile.Zoom( zoom_x )
	}

  img := NewImage( width, height )

	startPoint  := image.Point{}
	center_x, _ := divmod( overflow_x, 2 )
	startPoint.X = center_x      // note: center (add padding/overflow left & right)
	startPoint.Y = overflow_y    // note: anchor to bottom (if overflow)

  rect := image.Rect( startPoint.X, startPoint.Y,
		                    startPoint.X+img.Bounds().Dx(),
		                    startPoint.Y+img.Bounds().Dy())

	draw.Draw( img, rect, base, image.Point{0,0}, draw.Over )

  return img
}



func mkdirs( path string ) {
	// fmt.Println( "dirs: " + path )
  if path == "." {
	  return     // skip mkdir for current dir (e.g. .)
  }

  err := os.MkdirAll( path, 0755 )
  if err != nil {
		log.Fatalf( "failed to create directories >%s<: %#v", path, err )
	}
}

func (tile *Image) Save( path string ) {
	fmt.Printf( "  saving image to >%s<...\n", path )

  // note: auto-create directories in path - why? why not?
	mkdirs( filepath.Dir( path ) )

	fout, err := os.Create( path )
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	png.Encode( fout, tile )
}




