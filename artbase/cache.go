package artbase

import (
	"fmt"
	"image"
	"os"
	"errors"
	"sync"

	"github.com/pixelartexchange/artbase.server/pixelart"   // todo/check if relative to "root" or package ???
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


///
//  note: for now use one lock
//             per collection id (that is, name)
var mutex = make( map[string]sync.Mutex )



func (col *Collection) Image() *pixelart.ImageComposite  {
	name := col.Name

	lock := mutex[ name ]
	lock.Lock()
	defer lock.Unlock()

	path := col.Path

	if exist, _ := fileExist( path ); !exist {
  	fmt.Println( "    [artbase-cache] getting composite / download to (local) cache..." )
	  pixelart.Download( col.Url, path )
	}

	var img image.Image

  if Cache[ name ] != nil {
	   fmt.Println( "    [artbase-cache] bingo! (re)using in-memory composite image from cache...\n" )
	   img = Cache[ name ]
  } else {
	   fmt.Println( "    [artbase-cache] adding composite image to in-memory cache...\n" )
	   img = pixelart.ReadImage( path )
	   Cache[ name ] = img
  }

	return &pixelart.ImageComposite{ Image: img,
	                                 TileWidth:  col.Width,
												           TileHeight: col.Height }
}



