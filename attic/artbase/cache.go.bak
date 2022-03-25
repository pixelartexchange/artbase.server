

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


