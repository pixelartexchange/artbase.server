package artbase

import (
  "fmt"
  "io"
  "os"
	"net/http"
  "log"
  "strings"
  "strconv"
  "regexp"

  "encoding/csv"
)



var slug_re = regexp.MustCompile( "[ \t_-]" )

func slugify( str string ) string {
  str = strings.ToLower( str )
  str = slug_re.ReplaceAllString( str, "" )
  return str
}

func normify( str string ) string {
  str = strings.TrimSpace( str )
  return str
}



func parseBool( str string ) bool {
  str = strings.ToLower( str )
  str = strings.TrimSpace( str )

  value := false
  switch str {
    case "", "-", "0", "f", "false", "off", "no":
      value = false
    case "1", "t", "true", "on", "yes":
      value = true
    default:
      fmt.Printf( "WARN: ignoring unknown boolean value >%s<; defaulting to false\n", str )
  }
  return value
}



var space_re = regexp.MustCompile( "[ \t]" )

func parseNumber( str string ) (int) {
  str = space_re.ReplaceAllString( str, "" )
  num, _  := strconv.Atoi( str )
  return num
}


// convert dimension (width x height) "24x24" or "24 x 24" to  [24,24]
func parseDimension( str string ) (int,int) {
  str = strings.ToLower( str )
  str = space_re.ReplaceAllString( str, "" )
  dims := strings.Split( str, "x" )

  if len(dims) != 2 {
     fmt.Println( "dimension (size):" )
     fmt.Println( str, dims )
     log.Fatal( "expected dimension/size in the format 99x99; sorry" )
  }

  width,  _  := strconv.Atoi( dims[0] )
  height, _  := strconv.Atoi( dims[1] )

  return width, height
}



func ParseCollections( txt string ) []*Collection {

     r := csv.NewReader( strings.NewReader( txt ) )

    headers, err :=  r.Read()
    if err != nil {
      log.Fatal(err)
    }

    rows, err := r.ReadAll()
    if err != nil {
      log.Fatal(err)
    }

    var collections []*Collection  //:= make( []*Collection, size( rows ) )

    //////////////////////
    // headers to index
    fmt.Println(headers)

    name_idx       := -1
    size_idx       := -1
    url_idx        := -1
    background_idx := -1   // optional
    count_idx      := -1   // optional


    for index, header := range headers {
      header =  slugify( header )
      fmt.Printf("At index %d value is >%s<\n", index, header )
      switch header {
        case "name":
           name_idx = index
        case "size", "dim", "format":
           size_idx = index
        case "url", "src", "source", "href":
            url_idx = index
        case "background", "bg":
            background_idx = index
        case "count", "max":
            count_idx = index
        default:
          fmt.Printf( "WARN: skipping unknown column / header >%s<\n", header )
      }
    }

    if name_idx == -1 ||
       size_idx == -1 ||
       url_idx  == -1 {
        fmt.Println( "headers:" )
        fmt.Println( headers )
        log.Fatal( "one or more of required column(s) / header(s) missing - name/size/url; sorry")
    }


    /////
    // convert rows to records
    fmt.Println(rows)

    for index, row := range rows {

      name := normify( row[ name_idx ] )
      width, height := parseDimension( row[ size_idx ] )
      background := false
      if background_idx != -1 {
        background = parseBool( row[ background_idx ] )
      }
      url := normify( row[ url_idx ] )
      count := 0   // note: use 0 as zero (default) value :-) NOT -1 or such
      if count_idx != -1 {
        count = parseNumber( row[ count_idx] )
      }

      path := fmt.Sprintf( "./%s-%dx%d.png", name, width, height )

      fmt.Printf( "==> %d\n", index )
      fmt.Println( "  name: ", name )
      fmt.Printf( "  size: %d x %d px\n", width, height )
      fmt.Println( "  background: ", background )
      fmt.Println( "  url: ", url )
      fmt.Println( "  count: ", count )
      fmt.Println( "  path: ", path )


      col := &Collection{ Name: name,
                          Path: path,
                          Width: width,
                          Height: height,
                          Url: url,
                          Background: background,
                          Count: count }

      fmt.Printf( "  col: %t\n", col )

      collections = append( collections, col )
    }

    fmt.Println( "== collections:\n" )
    fmt.Printf( "  %t\n", collections )


    return collections
}



func ReadCollections( path string ) []*Collection {
   bytes, err := os.ReadFile( path )
   if err != nil {
      log.Fatal( err )
   }
   return ParseCollections( string(bytes) )
}


func DownloadCollections( url string ) []*Collection {
  fmt.Printf( "==> Downloading %s...\n", url )

  resp, err := http.Get( url )
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()

  // todo/fix:
  //  (double) check for content type - why? why not?

  fmt.Printf( "HTTP %v\n", resp.StatusCode )
  // dump  headers for debugging
  for name, headers := range resp.Header {
    for _, h := range headers {
      fmt.Printf( "  %v  :  %v\n", name, h )
    }
  }

   bytes, err := io.ReadAll( resp.Body )
   if err != nil {
      log.Fatal( err )
   }

   return ParseCollections( string(bytes) )
}

