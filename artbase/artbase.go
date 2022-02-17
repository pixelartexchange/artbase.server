package artbase


type Collection struct {
	Name         string
	Width        int
	Height       int
	Path         string
	Url          string
	// note:  background==false (default) => transparent
	//        background==true            => images have backgrounds (NOT transparent)
	Background   bool
	Count        int
}



