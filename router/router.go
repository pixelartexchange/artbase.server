package router


import (
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


/// not working - why?
//   resulting in:
//     serve.GET undefined (type func() *router.Router has no field or method GET)
//
// func New() *Router {
//	return &Router{}
//}



func (r *Router) Add( route Route ) {
 r.Routes = append( r.Routes, route )
}


func (r *Router) GET( pattern string, handler http.HandlerFunc ) {
	log.Printf( "  adding route - GET %s\n", pattern )

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
		 log.Printf( "==> Bingo! route %s %s matching:\n", req.Method, req.URL.Path )
		 log.Println( route )
		 log.Println( "  params:" )
		 log.Println( params )

		 // Create new request with params stored in context
		 ctx := context.WithValue( req.Context(), "params", params )
		 route.Handler.ServeHTTP( w, req.WithContext(ctx) )
		 return
 }

 // http.NotFound(w, r)
 log.Printf( "==> 404! no route %s %s matching\n", req.Method, req.URL.Path )
 http.NotFound( w, req )
}




//////////////////////////////////
// url param & query helpers
//


// URLParam extracts a parameter from the URL by name
func Param( req *http.Request, name string ) (string,bool) {
 ctx := req.Context()

 // ctx.Value returns an `interface{}` type, so we
 // also have to cast it to a map, which is the
 // type we'll be using to store our parameters.
 params := ctx.Value( "params" ).( map[string]string )
 value, ok := params[name]

 return value, ok
}

// todo/check:  change to ParamUint - why? why not?
func ParamInt( req *http.Request, name string ) (int,bool) {
	value := 0    // default default_value to zero 0 for now
	var err error

	param, ok := Param( req, name )
	if ok {
		value, err = strconv.Atoi( param )
		// value, err = strconv.ParseUint( param, 10, 0 )  // note: only allow natural/positive (0,1,2) numbers (not negative)
    if err != nil {
			log.Panic( err )
		}
	}
  return value, ok
}


func Query( req *http.Request, name string ) (string,bool) {
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

// todo/check:  change to QueryUint - why? why not?
func QueryInt( req *http.Request, name string ) (int,bool) {
  value := 0  // default default_value to 0 for now - why? why not?
  var err error

	query, ok := Query( req, name )
  if ok {
		value, err = strconv.Atoi( query )
		// value, err = strconv.ParseUint( query, 10, 0 )  // note: only allow natural/positive (0,1,2) numbers (not negative)
    if err != nil {
			log.Panic( err )
		}
	}
	return value, ok
}

func QueryBool( req *http.Request, name string ) (bool,bool) {
	value := false   // default default_value to false for now - why? why not?
  var err error

	query, ok := Query( req, name )
	if ok {
		value, err = strconv.ParseBool( query )
		// note: supports  0/1, f/t, F/T, False/True, FALSE/TRUE
    if err != nil {
			log.Panic( err )
		}
	}
  return value, ok
}


