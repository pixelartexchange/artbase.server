package main

import (
	"fmt"
	// "image"
	// "image/png"
	"image/color"
	"log"
  // "strings"
	"strconv"
	// "bytes"
	"regexp"
	"context"
	"net/http"

	"./artbase"
	"./pixelart"
)



////////
// based / inspired by  Build a go router from scratch
//   https://codesalad.dev/blog/how-to-build-a-go-router-from-scratch-3
//
// and https://github.com/gorilla/mux/blob/master/mux.go
//     https://xujiajun.cn/gorouter/
//     https://www.alexedwards.net/blog/which-go-router-should-i-use
//
//     https://benhoyt.com/writings/go-routing/



type Route struct {
	Method   string
	 Pattern  *regexp.Regexp       //  change to regex/re/rx or path or ___?
	Handler  http.HandlerFunc
}


type Router struct {
	Routes []Route
}


func (r *Router) Add( route Route ) {
 r.Routes = append( r.Routes, route )
}


func (r *Router) GET( pattern string, handler http.HandlerFunc ) {
	fmt.Printf( "  adding route - GET %s\n", pattern )

	exactPattern := regexp.MustCompile( "^" + pattern + "$" )

	r.Add( Route{ Method:  http.MethodGet,    // "GET" ?????
								Pattern: exactPattern,
								Handler: handler } )
}



func (r *Router) ServeHTTP( w http.ResponseWriter, req *http.Request ) {

 // try to recover from server errors / panics
 defer func() {
	 if err := recover(); err != nil {
			 log.Println( "ERROR:", err ) // Log the error
			 http.Error( w, "Uh oh!", http.StatusInternalServerError )
	 }
 }()


 for _, route := range r.Routes {

	 if req.Method != route.Method {
		 continue  // Method mismatch
	 }

	 match := route.Pattern.FindStringSubmatch( req.URL.Path )
	 if match == nil {
		 continue  // No match found
	 }

	 // Create a map to store URL parameters in
	 params := make( map[string]string )
	 groupNames := route.Pattern.SubexpNames()
	 for i, group := range match {
			 params[groupNames[i]] = group
	 }

	 // Create new request with params stored in context
		 // ctx := context.WithValue(r.Context(), "params", params)
		 // e.HandlerFunc.ServeHTTP(w, r.WithContext(ctx))
		 // return
		 fmt.Printf( "==> Bingo! route %s %s matching:\n", req.Method, req.URL.Path )
		 fmt.Println( route )
		 fmt.Println( "  params:" )
		 fmt.Println( params )

		 // Create new request with params stored in context
		 ctx := context.WithValue( req.Context(), "params", params )
		 route.Handler.ServeHTTP( w, req.WithContext(ctx) )
		 return
 }

 // http.NotFound(w, r)
 fmt.Printf( "==> 404! no route %s %s matching\n", req.Method, req.URL.Path )
 http.NotFound( w, req )
}



// URLParam extracts a parameter from the URL by name
func URLParam( req *http.Request, name string ) string {
 ctx := req.Context()

 // ctx.Value returns an `interface{}` type, so we
 // also have to cast it to a map, which is the
 // type we'll be using to store our parameters.
 params := ctx.Value( "params" ).( map[string]string )
 return params[name]
}

func URLQuery( req *http.Request, name, default_value string ) string {
	value := req.URL.Query().Get( name )
	if value == "" {
     value = default_value
	}
  return value
}






const (
	ContentType         = "Content-Type"
	ContentTypeJSON     = "application/json"
	ContentTypeHTML     = "text/html; charset=utf-8"
	ContentTypeText     = "text/plain; charset=utf-8"
	ContentTypeImagePNG = "image/png"
	ContentTypePNG      = ContentTypeImagePNG
	ContentTypeImageSVG = "image/svg+xml"
	ContentTypeSVG      = ContentTypeImageSVG
)




func handleHome( collections []artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {
		 b := artbase.RenderHome( collections )

		 w.Header().Set( ContentType, ContentTypeHTML )
		 w.Write( b )

		 // was: ctx.Data( http.StatusOK, ContentTypeHTML, b )
		}
}

func handleCollection( col artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {
		 b := artbase.RenderCollection( &col )

		 w.Header().Set( ContentType, ContentTypeHTML )
		 w.Write( b )

		 // was: ctx.Data( http.StatusOK, ContentTypeHTML, b )
	}
}



func handleCollectionImagePNG( col artbase.Collection ) http.HandlerFunc  {
		return func( w http.ResponseWriter, req *http.Request ) {

	id, _ := strconv.Atoi( URLParam( req, "id" ) )

	opts := artbase.PNGOpts{}


  var background color.Color = nil   // interface by default (zero-value) nil??
  var err error              = nil   // interface by default (zero-value) nil??

	backgroundParam := URLQuery( req, "background", "" )
	if backgroundParam == "" {
		backgroundParam = URLQuery( req, "bg", "" )  // allow shortcut z too
	}

	if backgroundParam != "" {
		 fmt.Printf( "=> parsing background color (in hex) >%s<...\n", backgroundParam )

     background, err = pixelart.ParseColor( backgroundParam )
     if err != nil {
			 // todo/fix:  only report parse color error and continue? why? why not?
			 log.Panic( err )
		 }

		 opts.Background     = background
		 opts.BackgroundName = backgroundParam
	}

	mirrorParam := URLQuery( req,  "mirror", "0" )
	if mirrorParam == "0" || mirrorParam[0] == 'f' || mirrorParam[0] == 'n'  {
		mirrorParam = URLQuery( req,  "m", "0" )  // allow shortcut m too
	}

  if mirrorParam == "1" || mirrorParam[0] =='t' || mirrorParam[0] =='y' {
		opts.Mirror = true
	}


	zoomParam := URLQuery( req,  "zoom", "0" )
	if zoomParam == "0" {
		zoomParam = URLQuery( req,  "z", "1" )  // allow shortcut z too
	}

	zoom, _ := strconv.Atoi( zoomParam )

  if zoom > 1 {
		opts.Zoom = zoom
	}


	saveParam := URLQuery( req,  "save", "0" )
	if saveParam == "0" || saveParam[0] == 'f' || saveParam[0] == 'n'  {
		saveParam = URLQuery( req,  "s", "0" )  // allow shortcut s too
	}

  if saveParam == "1" || saveParam[0] =='t' || saveParam[0] =='y' {
    opts.Save = true
	}


	b := col.HandleTilePNG( id, opts )

	w.Header().Set( ContentType, ContentTypeImagePNG )
	w.Write( b )
	// was: ctx.Data( http.StatusOK, ContentTypeImagePNG, bytes )
  }
}


func handleCollectionImageSVG( col artbase.Collection ) http.HandlerFunc  {
  return func( w http.ResponseWriter, req *http.Request ) {

  id, _ := strconv.Atoi( URLParam( req, "id" ) )

		opts := artbase.SVGOpts{}

mirrorParam := URLQuery( req, "mirror", "0" )
if mirrorParam == "0" || mirrorParam[0] == 'f' || mirrorParam[0] == 'n'  {
	mirrorParam = URLQuery( req, "m", "0" )  // allow shortcut m too
}

if mirrorParam == "1" || mirrorParam[0] =='t' || mirrorParam[0] =='y' {
	  opts.Mirror = true
}

	saveParam := URLQuery( req, "save", "0" )
	if saveParam == "0" || saveParam[0] == 'f' || saveParam[0] == 'n'  {
		saveParam = URLQuery( req, "s", "0" )  // allow shortcut s too
	}

  if saveParam == "1" || saveParam[0] =='t' || saveParam[0] =='y' {
    opts.Save = true
  }


	b := col.HandleTileSVG( id, opts )

	w.Header().Set( ContentType, ContentTypeImageSVG )
	w.Write( b )
  // was: ctx.Data( http.StatusOK, ContentTypeImageSVG, bytes )
  }
}


func main() {


	fmt.Printf( "%d collection(s):\n", len( artbase.Collections ))
	fmt.Println( artbase.Collections )

  fmt.Println( "cache:" )
  fmt.Println( artbase.Cache )


  //// note:
	// use built-in "standard" collections for now,
	//   yes, you can - use / set-up your own collections
	collections := artbase.Collections


	var router Router

	router.GET( "/",  handleHome( collections ) )

	for i,c := range collections {
		fmt.Printf( "  [%d] %s  %dx%d - %s\n", i, c.Name, c.Width, c.Height, c.Path )

		router.GET( "/" + c.Name,  handleCollection( c ) )

		// note - &c will NOT work - as c as reference gets
		//          all handlers pointing to last collection!!!!
		router.GET( "/" + c.Name + `/(?P<id>[0-9]+)(\.png)?`, handleCollectionImagePNG( c ) )
		router.GET( "/" + c.Name + `/(?P<id>[0-9]+)\.svg`,    handleCollectionImageSVG( c ) )
	}

	http.ListenAndServe( "localhost:8080", &router )

	fmt.Println( "Bye!")
}

