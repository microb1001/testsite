package myfs

import (
	"net/http"
	"strings"
	"os"
	"strconv"
	"fmt"
	"../servo"
)

//
// позволяет отраничить FileServer по типам отдаваемых файлов
// пример управляемой FileSystem почитать здесь https://golang.org/src/net/http/fs.go
type FileSystem struct {
	http.FileSystem
}

func (mfs FileSystem) Open(fname string) (http.File, error) {
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

	if err != nil && (folder=="pre/"||folder=="200/"||folder=="300/"||folder=="400/"||folder=="100/"){
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

