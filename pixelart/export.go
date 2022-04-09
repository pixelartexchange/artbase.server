package pixelart


import (
	"fmt"
	"bytes"
  "image"
)


/////
// change Export name to something like to toTXT
//         or toArray or toString ??? - why? why not?

func (tile *Image) Export() string {
	var buf bytes.Buffer

  // get pix array
	// fmt.Print( tile.NRGBA )
	// Pixs -  []
	// Stride - 96    (24*4)  row lenght in bytes
	// Rect/Bounds (0,0)-(24,24)
	// fmt.Println()

  mem := tile.NRGBA   // "in-memory" image pointer

  width  := mem.Stride / 4
	height := len( mem.Pix ) / mem.Stride

	fmt.Printf( "  %d x %d (width x height)\n", width, height )

	for y:=0; y < height; y+= 1 {
		buf.WriteString( "  " )
		for x:=0; x < width; x+=1 {
			offset := y*mem.Stride+x*4
			pix := uint32(mem.Pix[offset+3]) |
			       uint32(mem.Pix[offset+2])<<8 |
						 uint32(mem.Pix[offset+1])<<16 |
						 uint32(mem.Pix[offset+0])<<24

		  // fmt.Printf( "[%d/%d] %x\n", x, y, pix)
			if pix == 0 {
			  buf.WriteString( "0, " )
			} else {
			  buf.WriteString( fmt.Sprintf( "0x%08x, ", pix ))
			}
		}
		buf.WriteString( "\n" )
	}

	return buf.String()
}



func MakeImage( pix []uint32, imageSize *image.Point ) *Image {
  // auto-build in-memory image from scratch
	var nrgba *image.NRGBA = &image.NRGBA{ Pix: ToByte( pix ),
		                                     Stride: imageSize.X*4,
                                         Rect: image.Rect( 0, 0, imageSize.X, imageSize.Y),
	                                      }

	return &Image{ nrgba }
}



// convert rgba "quadruples" to byte array []uint8/byte
//   note:  uint8 == byte
//  todo/check: keep function internal only - why? why not?

func ToByte( data []uint32 ) []byte {
  b := make( []byte, len(data)*4 )

	for i,num := range data {
	   values := [4]byte{
            			byte(num >> 24),
			            byte(num >> 16),
			            byte(num >> 8),
			            byte(num),
		               }
	   // binary.BigEndian.PutUint32(values[0:4], num )

		 // fmt.Printf( "[%d] %08x => %v\n", i, num, values )
		 b[i*4]   = values[0]
		 b[i*4+1] = values[1]
		 b[i*4+2] = values[2]
		 b[i*4+3] = values[3]
	}
	return b
}

