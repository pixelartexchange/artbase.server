package main


import (
	"fmt"
	"log"
	"strconv"
	"regexp"
	"context"
	"net/http"
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



//////////////////////////////////
// url param & query helpers
//


// URLParam extracts a parameter from the URL by name
func URLParam( req *http.Request, name string ) (string,bool) {
 ctx := req.Context()

 // ctx.Value returns an `interface{}` type, so we
 // also have to cast it to a map, which is the
 // type we'll be using to store our parameters.
 params := ctx.Value( "params" ).( map[string]string )
 value, ok := params[name]

 return value, ok
}

func URLParamInt( req *http.Request, name string ) (int,bool) {
	value := 0    // default default_value to zero 0 for now
	var err error

	param, ok := URLParam( req, name )
	if ok {
		value, err = strconv.Atoi( param )
    if err != nil {
			log.Panic( err )
		}
	}
  return value, ok
}


func URLQuery( req *http.Request, name string ) (string,bool) {
	value := ""   // default default_value to empty string for now

	q := req.URL.Query()
	values, ok := q[name]
	if ok {
		/// note: query returns values (string array)
		//   for now care only about the first entry
	  value = values[0]
	}
  return value, ok
}

func URLQueryInt( req *http.Request, name string ) (int,bool) {
  value := 0  // default default_value to 0 for now - why? why not?
  var err error

	query, ok := URLQuery( req, name )
  if ok {
		value, err = strconv.Atoi( query )
    if err != nil {
			log.Panic( err )
		}
	}
	return value, ok
}

func URLQueryBool( req *http.Request, name string ) (bool,bool) {
	value := false   // default default_value to false for now - why? why not?

	query, ok := URLQuery( req, name )
	if ok {
		if query == "1" || query[0] == 't' || query[0] == 'y'  {
			value  = true
		}
	  // note: for now assume everything else is false
		//    e.g. 0, f(alse), n(o), etc.
	}
  return value, ok
}


