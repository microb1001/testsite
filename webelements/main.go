package webelements

import (
	"strconv"
	"net/http"
	"os"
	"strings"
	"fmt"
	"log"
	"image"
	_ "image/jpeg"
	//"image/png"
	"image/color"
	"time"
	"image/jpeg"

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
	type MyFs struct {
		http.FileSystem
	}

	func (mfs MyFs) Open(fname string) (http.File, error) {
		var name, folder, ext string
		var f http.File

		Splits :=strings.SplitAfter(fname,".")
		if len(Splits)==2 && (Splits[1]=="jpg" || Splits[1]=="csv" || Splits[1]=="js" || Splits[1]=="png"){
			ext= Splits[1]
		} else {return nil,os.ErrPermission}

		Splits = strings.SplitAfter(Splits[0][1:],"/")
		if len(Splits)==1 {name=Splits[0]; folder=""
		} else
		if len(Splits)==2 {name=Splits[1]; folder=Splits[0]
		} else {return nil,os.ErrPermission}

		nm:= "/"+folder+name+ext
		f, err := mfs.FileSystem.Open(nm)

		if err != nil && folder=="pre/"{
			maindir,ok:=mfs.FileSystem.(http.Dir)
			if !ok {return nil,os.ErrPermission}
			basedir:=string (maindir)
			fmt.Println(basedir)
			thumb(basedir+name+ext,basedir+"pre/"+name+ext,300,300)
			f, err = mfs.FileSystem.Open("pre/"+name+ext)
		}

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

	func thumb (ifile, ofile string, maxWidth, maxHeight int ) error {
		start := time.Now()

		reader, err1 := os.Open(ifile)
		if err1 != nil {return err1}
		defer reader.Close()
		//var teal color.Color = color.RGBA{0, 200, 200, 255}
		//	var red  color.Color = color.RGBA{200, 30, 30, 255}
		var m image.Image
		m, _, err1 = image.Decode(reader)
		if err1 != nil {return err1}
		file, err := os.Create(ofile)
		if err != nil {	return err1	}

		img := image.NewRGBA64(image.Rect(0, 0, maxWidth, maxHeight))
		//draw.Draw(img, img.Bounds(), &image.Uniform{teal}, image.ZP, draw.Src)

		var stepX,stepY float32
		stepX=float32(m.Bounds().Dx())/float32(img.Bounds().Dx())
		stepY=float32(m.Bounds().Dy())/float32(img.Bounds().Dy())
		for x:= img.Bounds().Min.X;x<img.Bounds().Max.X;x++{
			for y:= img.Bounds().Min.Y;y<img.Bounds().Max.Y;y++{
				var cnt uint32 =0
				var R,G,B,A uint32 =0,0,0,0
				for mx:= int(float32(x)*stepX);mx<int(float32(x+1)*stepX);mx++{
					for my:= int(float32(y)*stepY);my<int(float32(y+1)*stepY);my++{
						cnt++
						R1,G1,B1,A1:=m.At(mx,my).RGBA()
						R+=R1;G+=G1;B+=B1;A+=A1;
					}
				}
				img.SetRGBA64(x,y,color.RGBA64{uint16(R/cnt),uint16(G/cnt),uint16(B/cnt),uint16(A/cnt)})
			}
		}

		//png.Encode(file, img)
		err=jpeg.Encode(file, img,nil)
		if err != nil {return err1}
		file.Close()
		t := time.Now()
		elapsed := t.Sub(start)
		log.Println("timer ==", elapsed)
		return nil










	}


/*
	func main() {
		fs := justFiles Filesystem{http.Dir("/tmp/")}
		http.ListenAndServe(":8080", http.FileServer(fs))
	}


func
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
