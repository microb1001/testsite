package webelements

import (
	"strconv"
	"net/http"
	"os"
	"strings"
	"fmt"
	_ "image/jpeg"
	//"image/png"
	"bytes"
	"../webelements/servo"
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
//
// Возвращает структуру для создания страниц в template
// i,j какие элементы списка сейчас будут отображаться
//
func Pager(Page, items_per_page, itemsCnt int, urlPart string) (newP PagerType, i, j int) {
	const PAGERWIDTH  = 2
	maxPage := (itemsCnt - 1) / items_per_page // начинается с нуля
	if Page > 0 {
		newP.Prev = urlPart + "p=" + strconv.Itoa(Page-1)
	}
	if Page < maxPage {
		newP.Next = urlPart + "p=" + strconv.Itoa(Page+1)
	}
	for ii := servo.MinMax(Page-PAGERWIDTH, 0, maxPage); ii <= servo.MinMax(Page+PAGERWIDTH, 0, maxPage); ii++ {
		newP.Elem = append(newP.Elem, PagerElemType{ii + 1, "", urlPart + "p=" + strconv.Itoa(ii), ii == Page})
	}
	i = servo.MinMax((Page)*items_per_page, 0, itemsCnt)
	j = servo.MinMax((Page+1)*items_per_page, 0, itemsCnt)
	return
}
//
// позволяет отраничить FileServer по типам отдаваемых файлов
// пример управляемой FileSystem почитать здесь https://golang.org/src/net/http/fs.go
//

type MyFs struct {
	http.FileSystem
}

func (mfs MyFs) Open(fname string) (http.File, error) {
		var name, folder, ext string
		var f http.File

		Splits :=strings.SplitAfter(fname,".")
		if len(Splits)==2 && (Splits[1]=="jpg" || Splits[1]=="css" || Splits[1]=="js" || Splits[1]=="png"){
			ext= Splits[1]
		} else {return nil,os.ErrPermission}

		Splits = strings.SplitAfter(Splits[0][1:],"/")
		if len(Splits)==1 {name=Splits[0]; folder=""
		} else
		if len(Splits)==2 {name=Splits[1]; folder=Splits[0]
		} else {return nil,os.ErrPermission}

		nm:= "/"+folder+name+ext // надо поствить проверку только английские буквы и цифры
		f, err := mfs.FileSystem.Open(nm)

		if err != nil && (folder=="pre/"||folder=="200/"||folder=="300/"||folder=="400/"){
			// нет файла - содать
			mfsAsDir,ok:=mfs.FileSystem.(http.Dir)
			if !ok {return nil,os.ErrPermission}
			basedir:=string (mfsAsDir)
			imgDim,err8:=strconv.Atoi(folder[0:len(folder)-1]) // Интересная ошибка если здесь err, создается локальная переменная и глобальная err не переписывается, потом ошибка нет файла
			if err8!=nil{
				fmt.Println("a1",err8)
				imgDim=150
			}
			servo.Thumb(basedir+name+ext,basedir+folder+name+ext,imgDim,imgDim) // опасно если folder или name могут содержать что нибудь кроме текста

			f, err = mfs.FileSystem.Open(folder+name+ext)

		}
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}


type Breadcrumbs_type []struct{
	Url,Name string
}
func Breadcrumbs() Breadcrumbs_type {
return nil
}


//
// фильтр символов, второй аргумент функция фильтрования
// rune это int32
//
func Only_EngAndDigit(r rune) rune {
	if (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') {
		return r
	}
	return -1
}
func Only_Digit(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}
func Only_NotSpecial(r rune) rune { // не рабочий исходный пример
	switch {
	case r == 'a':
		return -1
	case r >= 'A' && r <= 'Z':
		return 'A' + (r-'A'+13)%26
	case r >= 'a' && r <= 'z':
		return -1
	}
	return -1
}
func Only_RusAndDigit(r rune) rune {
	if (r >= 'а' && r <= 'я') ||
		(r >= 'А' && r <= 'Я') ||
		(r >= '0' && r <= '9') {
		return r
	}
	return -1
}
type Only_type struct{
	Min, Max rune
}
func Only_MakeFn(Allowed [][]rune) func(rune) rune {
	return func(r rune) (ret rune) {
		for _,k:=range Allowed{
			if len(k)==2&&r>=k[0]&&r<=k[1]{
				return r
			}
		}
		return -1
	}
}

func OnlyS(s_in []byte, F func(r rune) rune) (s_out []byte) {
	return bytes.Map(F, s_in)
}

