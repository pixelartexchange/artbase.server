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


////////////////////////////////
// tile methods  "convenience helpers" for easy chaining

func (tile *ImageTile) Background( background color.Color ) *ImageTile {

	// todo/fix: change to newNRGBA (better match for png - why? why not?)
	width, height := tile.Bounds().Dx(), tile.Bounds().Dy()
	img := image.NewRGBA( image.Rect(0,0, width, height) )

	/// use Image.ZP for image.Point{0,0} - why? why not?
	draw.Draw( img, img.Bounds(), &image.Uniform{ background }, image.Point{0,0}, draw.Src )
	draw.Draw( img, img.Bounds(), tile, image.Point{0,0}, draw.Over )

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




