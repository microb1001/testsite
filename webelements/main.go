package webelements

import (
	"strconv"
	"net/http"
	"os"
	"strings"
)

type PagerElemType struct{
	Page int
	Class string
	Url string
	Current bool
}

type PagerType struct{
	Elem [] PagerElemType
	Next string
	Prev string
}

func MinMax(index,min,max int) int{
	if index<min {
		return min
	}
	if index>max {
		return max
	}
	return index
}

func Pager (Page, items_per_page, itemsCnt int, urlPart string) (newP PagerType,i,j int) {
	maxPage :=(itemsCnt-1)/ items_per_page // начинается с нуля
	if Page > 0 {newP.Prev= urlPart +"p="+strconv.Itoa(Page-1)}
	if Page < maxPage {newP.Next= urlPart +"p="+strconv.Itoa(Page+1)}
	for ii:=MinMax(Page-2,0, maxPage);ii<=MinMax(Page+2,0, maxPage);ii++{
		newP.Elem=append(newP.Elem,PagerElemType{ii+1,"", urlPart +"p="+strconv.Itoa(ii),ii == Page} )
	}
	i=MinMax((Page)*items_per_page,0, itemsCnt)
	j=MinMax((Page+1)*items_per_page,0, itemsCnt)
	return
}


/*  пример управляемой FileServer
почитать здесь https://golang.org/src/net/http/fs.go
	) */

	type MyWebFilesystem struct {
		fs http.FileSystem
	}

	func (mfs MyWebFilesystem) Open(name string) (http.File, error) {
		name1:=strings.Split(name,".")
		if (len(name1)!=2) || (name1[1]!="jpg" && name1[1]!="csv" && name1[1]!="js" && name1[1]!="png") {
			return nil,os.ErrPermission
		}
		name2:=strings.Split(name1[0],"/")
		if (len(name2)>2) {
			return nil,os.ErrPermission
		}
	f, err := mfs.fs.Open(name2[0]+"."+name1[1])
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}
	return neuteredReaddirFile{f}, nil
	}

	type neuteredReaddirFile struct {
		http.File
	}

	func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
	}
/*
	func main() {
		fs := justFiles Filesystem{http.Dir("/tmp/")}
		http.ListenAndServe(":8080", http.FileServer(fs))
	}


I've settled on the following, which has the added benefit of returning 404 for directories.


Corrections appreciated.



type justFilesFilesystem struct {
   Fs http.FileSystem

}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
   f, err := fs.Fs.Open(name)

   if err != nil {
      return nil, err
   }

   stat, err := f.Stat()
   if stat.IsDir() {
      return nil, os.ErrNotExist
   }

   return f, nil
}

// A File is returned by a FileSystem's Open method and can be

// served by the FileServer implementation.

//

// The methods should behave the same as those on an *os.File.

type File interface {

	io.Closer

	io.Reader

	io.Seeker

	Readdir(count int) ([]os.FileInfo, error)

	Stat() (os.FileInfo, error)

}
// you need the image package, and a format package for encoding/decoding
import (
    "bytes"
    "image"
    "image/jpeg"

    // if you don't need to use jpeg.Encode, import like so:
    // _ "image/jpeg"
)

// Decoding gives you an Image.
// If you have an io.Reader already, you can give that to Decode
// without reading it into a []byte.
image, _, err := image.Decode(bytes.NewReader(data))
// check err

newImage := resize.Resize(160, 0, original_image, resize.Lanczos3)

// Encode uses a Writer, use a Buffer if you need the raw []byte
err = jpeg.Encode(someWriter, newImage, nil)
// check err

*/
