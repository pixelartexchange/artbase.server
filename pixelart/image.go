package pixelart

import (
	"fmt"
	"image"
)




func ZoomImage(img image.Image, zoom int) (*image.RGBA, error) {

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()  // bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	fmt.Println( bounds, width, height )
	// e.g.   punk #0    (0,0)-(24,24)
  //             #561  (1464,120)-(1488,144)
  //             #3100 (0,744)-(24,768)
  //             #7804 (96,1872)-(120,1896)
  //             #8857 (1368,2112)-(1392,2136)

	new_img := image.NewRGBA( image.Rect(0,0, width*zoom, height*zoom) )

	for x:=0; x < width; x++ {
		for y:=0; y < height; y++ {
				pixel := img.At( bounds.Min.X+x, bounds.Min.Y+y )
        for n:=0; n < zoom; n++ {
					for m:=0; m < zoom; m++ {
						new_img.Set( n+zoom*x, m+zoom*y, pixel )
					}
				}
	 }
}

	return new_img, nil
}



func MirrorImage(img image.Image) (*image.RGBA, error) {

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()  // bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y


	fmt.Println( bounds, width, height )
	// e.g.   punk #0    (0,0)-(24,24)
  //             #561  (1464,120)-(1488,144)
  //             #3100 (0,744)-(24,768)
  //             #7804 (96,1872)-(120,1896)
  //             #8857 (1368,2112)-(1392,2136)

	new_img := image.NewRGBA( image.Rect(0,0, width, height) )

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At( bounds.Min.X+x, bounds.Min.Y+y )
			new_img.Set( (width-1)-x, y, pixel )
		}
	}

	return new_img, nil
}



