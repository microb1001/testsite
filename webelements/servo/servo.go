package servo

import(
"crypto/rand"
	"fmt"
	"image/jpeg"
	"image"
	"image/color"
	"time"
	"log"
	"os"
)

// генерация истинно случайного числа из crypto
func Rand64 () uint64 {
	var res uint64
	rand64 := make([]byte, 8)
	rand.Read(rand64)
	for _,i:=range rand64 {res=res*256+uint64(i)}
	return res
}

// понятно из названия
func MinMax(index,min,max int) int{
	if index<min {
		return min
	}
	if index>max {
		return max
	}
	return index
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
