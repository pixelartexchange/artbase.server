package pixelart

import (
	"fmt"
	"image"
	"image/color"
	"bytes"
)




// todo/check: rename pixel_raw to pixel_uncast/generic/... or similar - why? why not?
func renderPixel( x, y int,
	                pixel color.NRGBA,
									colors map[string]*ColorInfo )  string {

	// convert color to rgba() hexstring
  hex := fmt.Sprintf( "#%02x%02x%02x%02x",
	                pixel.R, pixel.G, pixel.B, pixel.A )

	//return fmt.Sprintf( "<rect class=\"px\" x=\"%d\" y=\"%d\" width=\"1\" height=\"1\" fill=\"%s\"/>",
  //                          x, y, hex )

  colInfo, _ := colors[ hex ]

	return fmt.Sprintf( "<rect class=\"px c%d\" x=\"%d\" y=\"%d\" width=\"1\" height=\"1\"/>",
                            colInfo.Index+1, x, y )


}





// todo/check: find a better name - why? why not?
type ColorInfo struct {
	 Count int      // no. of pixels
	 Index int      // zero-based running index (0,1,2, etc.)
}

func imageColors( img image.Image ) map[string]*ColorInfo  {

  colors := make( map[string]*ColorInfo )


	transparent := color.NRGBA{ R: 0,
															 G: 0,
															 B: 0,
																A: 0 }

  bounds := img.Bounds()

	for y:=bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x:=bounds.Min.X; x < bounds.Max.X; x++ {
			// note: cast/convert interface image.color to type NRGBA (supporting image.color)
			// pixel, ok := img.At( bounds.Min.X+x, bounds.Min.Y+y ).(image.NRGBA)
			// r, g, b, a := pixel.R, pixel.G, pixel.B, pixel.A
			pixel := color.NRGBAModel.Convert( img.At( x, y )).(color.NRGBA)
			// if !ok {
			//	fmt.Printf( "!! ERROR - expected color.NRGBA pixel/color but got: %v %T\n",
			//																		 pixel, pixel )
			//	os.Exit(1)
			//}


			if pixel == transparent {
				 // do nothing; skip
			}  else {
				// convert color to rgba() hexstring
				hex := fmt.Sprintf( "#%02x%02x%02x%02x",
												pixel.R, pixel.G, pixel.B, pixel.A )

         colInfo, ok := colors[ hex ]
				 if !ok {
            colInfo = &ColorInfo{ Index: len( colors ) }
						colors[ hex ] = colInfo
				 }

				 colInfo.Count += 1
				 fmt.Printf( "%s - %v %v\n", hex, colInfo, *colInfo )
				}
	 }
	}

	return colors
}



func (tile *Image) ToSVG() string  {

		/////
		// type NRGBA
		//   type NRGBA struct {
	  //     R, G, B, A uint8
    //   }
    //  NRGBA represents a non-alpha-premultiplied 32-bit color.

		// The NRGBA struct type represents an 8-bit non-alpha-premultiplied color,
		//  as used by the PNG image format.
		// When manipulating an NRGBA’s fields directly,
		// the values are non-alpha-premultiplied,
		// but when calling the RGBA method, the return values are alpha-premultiplied.


    // transparent := color.NRGBA{ 0, 0, 0, 0 }
		transparent := color.NRGBA{ R: 0,
			G: 0,
			B: 0,
			 A: 0 }

    // todo/check: try adding
		//   height="24px/240px" width="24px/240px"
    //		viewBox="0 0 24 24"
		//   style="background: red;">


		colors := imageColors( tile )

		fmt.Println( "colors:" )
		fmt.Println( colors )


		bounds        := tile.Bounds()
		width, height := bounds.Dx(), bounds.Dy()


    var buf bytes.Buffer
    buf.WriteString( "<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.2\"\n" )
		buf.WriteString(
			fmt.Sprintf( "    width=\"%dpx\" height=\"%dpx\" viewBox=\"0 0 %d %d\">\n",
			   width, height, width, height ))
		buf.WriteString( "  <style>\n" )

		// add colors  - note: start with one (1) NOT zero (0)-based - why? why not?
		//    e.g rect.c1 { fill: #000000ff }
    for colHex,colInfo := range colors {
			buf.WriteString(
				fmt.Sprintf( "    rect.c%d {  fill: %s }  /* %3d pixel(s) */\n",
			         colInfo.Index+1, colHex, colInfo.Count ))
		}

		buf.WriteString( "    rect.px {  shape-rendering: crispEdges  }\n" )
		buf.WriteString( "  </style>\n" )



		for y:=0; y < height; y++ {
		  for x:=0; x < width; x++ {
				// note: cast/convert interface image.color to type NRGBA (supporting image.color)
				// pixel, ok := img.At( bounds.Min.X+x, bounds.Min.Y+y ).(image.NRGBA)
				// r, g, b, a := pixel.R, pixel.G, pixel.B, pixel.A
				pixel := color.NRGBAModel.Convert(
					            tile.At( bounds.Min.X+x, bounds.Min.Y+y )).(color.NRGBA)


					if pixel == transparent {
				}  else {
					   inst := renderPixel( x, y, pixel, colors )
					   line := fmt.Sprintf( "  %s\n", inst )

					buf.WriteString( line )
				}
		 }
		}

		buf.WriteString( "</svg>\n" )

		return buf.String()
}
