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

		if err != nil && (folder=="pre/"||folder=="200/"||folder=="300/"||folder=="400/"){
			mfsAsDir,ok:=mfs.FileSystem.(http.Dir)
			if !ok {return nil,os.ErrPermission}
			basedir:=string (mfsAsDir)
			imgDim,err8:=strconv.Atoi(folder[0:len(folder)-1]) // Интересная ошибка если здесь err, создается локальная переменная и глобальная err не переписывается, потом ошибка нет файла
			if err8!=nil{
				fmt.Println("a1",err8)
				imgDim=150
			}
			fmt.Println("aa",err8)
			Thumb(basedir+name+ext,basedir+folder+name+ext,imgDim,imgDim) // опасно если folder или name могут содержать что нибудь кроме текста

			f, err = mfs.FileSystem.Open(folder+name+ext)

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

	func Thumb(inFile, outFile string, maxWidth, maxHeight int ) error {
		start := time.Now()

		reader, err := os.Open(inFile)
		if err != nil {return err}
		defer reader.Close()
		inImage, tt, err := image.Decode(reader)
		fmt.Println(tt)
		if err != nil {return err}

		file, err := os.Create(outFile)
		if err != nil {return err}
		Height:=inImage.Bounds().Dy()*maxWidth/inImage.Bounds().Dx()
		Width:=inImage.Bounds().Dx()*maxHeight/inImage.Bounds().Dy()
		fmt.Println(inImage.Bounds().Dx(),inImage.Bounds().Dy(),Width,Height)
		if Height>maxHeight {Height=maxHeight}
		if Width>maxWidth {Width=maxWidth} // Получилось симметрично но странно. Одно лишнее вычисление
		fmt.Println(inImage.Bounds().Dx(),inImage.Bounds().Dy(),Width,Height,maxWidth, maxHeight)

		if Width>inImage.Bounds().Dx() || Height>inImage.Bounds().Dx() { // Не умеет повышать разрешение!
			err=jpeg.Encode(file, inImage,nil)
			if err != nil {return err}
			file.Close()
			return nil
		}

		outImage := image.NewRGBA64(image.Rect(0, 0, Width, Height))
		//draw.Draw(outImage, outImage.Bounds(), &image.Uniform{teal}, image.ZP, draw.Src) // Водяные знаки!

		var stepX,stepY float32
		stepX=float32(inImage.Bounds().Dx())/float32(outImage.Bounds().Dx())
		stepY=float32(inImage.Bounds().Dy())/float32(outImage.Bounds().Dy())
		for x:= outImage.Bounds().Min.X;x< outImage.Bounds().Max.X;x++{
			for y:= outImage.Bounds().Min.Y;y< outImage.Bounds().Max.Y;y++{
				var cnt uint32 =0
				var R,G,B,A uint32 =0,0,0,0
				for mx:= int(float32(x)*stepX);mx<int(float32(x+1)*stepX);mx++{
					for my:= int(float32(y)*stepY);my<int(float32(y+1)*stepY);my++{
						cnt++
						R1,G1,B1,A1:= inImage.At(mx,my).RGBA()
						R+=R1;G+=G1;B+=B1;A+=A1;
					}
				}
				outImage.SetRGBA64(x,y,color.RGBA64{uint16(R/cnt),uint16(G/cnt),uint16(B/cnt),uint16(A/cnt)})
			}
		}

		//png.Encode(file, outImage)
		err=jpeg.Encode(file, outImage,nil)
		if err != nil {return err}

		t := time.Now()
		elapsed := t.Sub(start)
		log.Println("timer ==", elapsed)
		log.Println(file.Close())
		return nil

	}
type Breadcrumbs_type []struct{
	Url,Name string
}
func Breadcrumbs() Breadcrumbs_type {
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
