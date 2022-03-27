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


func (tile *ImageTile) Background( background_any interface{} ) *ImageTile {

   background := MakeColor( background_any )

	// todo/fix: change to newNRGBA (better match for png - why? why not?)
	width, height := tile.Bounds().Dx(), tile.Bounds().Dy()
	img := image.NewRGBA( image.Rect(0,0, width, height) )

	/// use Image.ZP for image.Point{0,0} - why? why not?
	draw.Draw( img, img.Bounds(), &image.Uniform{ background }, image.Point{0,0}, draw.Src )
	draw.Draw( img, img.Bounds(), tile, image.Point{0,0}, draw.Over )

	return &ImageTile{ Image: img }
}


// draw flag of ukraine -- glory to ukraine! fuck (vladimir) putin! stop the war!
func (tile *ImageTile) Ukraine() *ImageTile {

	blue   := color.RGBA{0x00, 0x57, 0xb7, 0xff}  // rgb( 0, 87, 183 )
	yellow := color.RGBA{0xff, 0xdd, 0x00, 0xff}  // rgb( 255, 221, 0)

 // todo/fix: change to newNRGBA (better match for png - why? why not?)
 width, height := tile.Bounds().Dx(), tile.Bounds().Dy()
 img := image.NewRGBA( image.Rect(0,0, width, height) )

 /// use Image.ZP for image.Point{0,0} - why? why not?
 draw.Draw( img, image.Rect( 0,0,width, height/2 ),
	               &image.Uniform{ blue }, image.Point{0,0}, draw.Src )
 draw.Draw( img, image.Rect( 0,height/2,width,height ),
								 &image.Uniform{ yellow }, image.Point{0,0}, draw.Src )

 draw.Draw( img, img.Bounds(), tile, image.Point{0,0}, draw.Over )

 return &ImageTile{ Image: img }
}



func (tile *ImageTile) Silhouette( foreground_any interface{} ) *ImageTile {

	foreground := MakeColor( foreground_any )

	transparent := color.NRGBA{ R: 0,
                           		G: 0,
		                          B: 0,
		                          A: 0 }

 // todo/fix: change to newNRGBA (better match for png - why? why not?)
 bounds        := tile.Bounds()
 width, height := bounds.Dx(), bounds.Dy()
 img := image.NewRGBA( image.Rect(0,0, width, height) )

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

 return &ImageTile{ Image: img }
}



func (tile *ImageTile) Transparent() *ImageTile {

	bounds      := tile.Bounds()
	background := tile.At( bounds.Min.X,
		                     bounds.Min.Y ) // 0,0

  transparent := color.NRGBA{ R: 0,
														  G: 0,
													    B: 0,
													    A: 0 }

 // todo/fix: change to newNRGBA (better match for png - why? why not?)
 width, height := bounds.Dx(), bounds.Dy()
 img := image.NewRGBA( image.Rect(0,0, width, height) )


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


 return &ImageTile{ Image: img }
}



func (tile *ImageTile) Circle() *ImageTile {

	 bounds := tile.Bounds()
   width, height := bounds.Dx(), bounds.Dy()

	 // for radius use min of width / height
	 r := min( width, height ) / 2

   center_x := width  / 2
	 center_y := height / 2

   ////////
	 //  try with 96x96
	 //    center_x:  96 / 2 = 48
	 //    center_y:  96 / 2 = 48
   //
	 //     r:    96 / 2 = 48


	 // use color.Alpha{0} - why? why not?
	 transparent := color.NRGBA{ R: 0,
                            	 G: 0,
	                             B: 0,
	                             A: 0 }

   // todo/fix: change to newNRGBA (better match for png - why? why not?)
   img := image.NewRGBA( image.Rect(0,0, width, height) )

	 for x := 0; x < width; x++ {
     for y := 0; y < height; y++ {
         pixel := tile.At( bounds.Min.X+x,
										       bounds.Min.Y+y )

		xx, yy, rr := float64( x - center_x )+0.5,
		              float64( y - center_y )+0.5,
									float64( r )

						if xx*xx+yy*yy < rr*rr {
							img.Set(x, y, pixel )
						} else {
							img.Set( x,y, transparent )
						}
					}
		}
  	return &ImageTile{ Image: img }
}



func (tile *ImageTile) Zoom( zoom int ) *ImageTile {
			bounds := tile.Bounds()
			width, height := bounds.Dx(), bounds.Dy()  // note: same as bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

			// fmt.Println( bounds, width, height )
			// e.g.   punk #0    (0,0)-(24,24)
			//             #561  (1464,120)-(1488,144)
			//             #3100 (0,744)-(24,768)
			//             #7804 (96,1872)-(120,1896)
			//             #8857 (1368,2112)-(1392,2136)

			img := image.NewRGBA( image.Rect(0,0, width*zoom, height*zoom) )

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

		return &ImageTile{ Image: img }
}


func (tile *ImageTile) Mirror() *ImageTile {
		bounds := tile.Bounds()
		width, height := bounds.Dx(), bounds.Dy()  // bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

		// fmt.Println( bounds, width, height )
		// e.g.   punk #0    (0,0)-(24,24)
		//             #561  (1464,120)-(1488,144)
		//             #3100 (0,744)-(24,768)
		//             #7804 (96,1872)-(120,1896)
		//             #8857 (1368,2112)-(1392,2136)

		img := image.NewRGBA( image.Rect(0,0, width, height) )

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				pixel := tile.At( bounds.Min.X+x, bounds.Min.Y+y )
				img.Set( (width-1)-x, y, pixel )
			}
		}

		return &ImageTile{ Image: img }
}



func (tile *ImageTile) Paste( img image.Image ) {
	// note - image.Image (is read-only - no Set() method)
	//          convert/typecast to draw.Image (that includes Set() method)
  //  see https://stackoverflow.com/questions/36573413/change-color-of-a-single-pixel-golang-image

	draw.Draw( tile.Image.(draw.Image),
	           img.Bounds(), img, image.Point{0,0}, draw.Over )
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




