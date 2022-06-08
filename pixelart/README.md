# pixelart Package

Let's You Read and Write Pixel Art Images in the Portable Network Graphics (PNG) Format Including Support for Composites and Special Effects (Mirror, Zoom, Flip, etc.) and More


Example:


``` go
package main


import (
  "fmt"
  "github.com/pixelartexchange/artbase.server/pixelart"
)


var dir = "../basic"


func main() {
  fmt.Printf( "Hello, Pixel Art v%s!\n", pixelart.Version )

  ///////////
  // read in f(emale) attributes
  female2        := pixelart.ReadImage( dir + "/female2.png" )
  earring        := pixelart.ReadImage( dir + "/f/earring.png" )
  blondebob      := pixelart.ReadImage( dir + "/f/blondebob.png" )
  greeneyeshadow := pixelart.ReadImage( dir + "/f/greeneyeshadow.png" )

  // test drive
  // generate punk #0
  punk := pixelart.NewImage( 24, 24 )
  punk.Paste( female2 )
  punk.Paste( earring )
  punk.Paste( blondebob )
  punk.Paste( greeneyeshadow )

  punk.Save( "./punk0.png" )
  punk.Zoom(20).Save( "./punk0@20x.png" )

  // (re)try with background
  punk = pixelart.NewImage( 24, 24 ).Background( "#60A4F7" )
  punk.Paste( female2 )
  punk.Paste( earring )
  punk.Paste( blondebob )
  punk.Paste( greeneyeshadow )

  punk.Save( "./bluepunk0.png" )
  punk.Zoom(20).Save( "./bluepunk0@20x.png" )


  ///////////
  // read in m(ale) attributes
  male1   := pixelart.ReadImage( dir + "/male1.png" )
  smile   := pixelart.ReadImage( dir + "/m/smile.png" )
  mohawk  := pixelart.ReadImage( dir + "/m/mohawk.png" )

  // generate punk #1
  punk = pixelart.NewImage( 24, 24 )
  punk.Paste( male1 )
  punk.Paste( smile )
  punk.Paste( mohawk )

  punk.Save( "./punk1.png" )
  punk.Zoom(20).Save( "./punk1@20x.png" )

  // (re)try with background
  punk = pixelart.NewImage( 24, 24 ).Background( "#60A4F7" )
  punk.Paste( male1 )
  punk.Paste( smile )
  punk.Paste( mohawk )

  punk.Save( "./bluepunk1.png" )
  punk.Zoom(20).Save( "./bluepunk1@20x.png" )

  fmt.Println( "Bye")
}
```



For more see [**Let's Go! Programming (Crypto) Pixel Punk Profile Pictures & (Generative) Art with Go - Step-by-Step Book / Guide »**](https://github.com/cryptopunksnotdead/lets-go-programming-cryptopunks)




## Questions? Comments?

Yes, you can. Post them on the [D.I.Y. Punk (Pixel) Art reddit](https://old.reddit.com/r/DIYPunkArt). Thanks.
