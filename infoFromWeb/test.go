package main

import (
	"fmt"
	"reflect"

)

type AppConfig struct {
	Pg    string `cli:"pg"    env:"PG"    default:"host=host.local dbname=db user=user password=password" description:"Connection to PostgreSQL"`
	Redis string `cli:"redis" env:"REDIS" default:"host.local"                                            description:"Redis server"`
}
var t AppConfig = AppConfig{Pg:"143", Redis:"Red"}
func main2() {
	GetConfig(&t)
//fmt.Println(&config)
}

func GetConfig(config interface{}) {

	ref1 := reflect.TypeOf(config)
	fmt.Println(ref1)

	value := reflect.ValueOf(config).Elem()
	fmt.Println(value)

	ref := value.Type()
	fmt.Println(ref)

	for i := 0; i < value.NumField(); i++ {
		field := ref.Field(i)

		def:= field.Tag.Get("default")
		fmt.Println(def)
		value.Field(i).SetString("236")
		fmt.Println(value.Field(i).Type())


	}

	return
}
// This works too, but "/static2/" fragment remains and need to be striped manually
//http.HandleFunc("/static2/", func(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, r.URL.Path[1:])
// })
//NOTE Serving up a filesystem naively is a potentially dangerous thing (there are potentially ways to break out of the rooted tree) hence I recommend that unless you really know what you're doing, use http.FileServer and http.Dir as they include checks to make sure people can't break out of the FS, which ServeFile doesn't.
/*  пример управляемой FileServer
package main

	import (
		"net/http"
	"os"
	)

	type justFilesFilesystem struct {
		fs http.FileSystem
	}

	func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
	return nil, err
	}
	return neuteredReaddirFile{f}, nil
	}

	type neuteredReaddirFile struct {
		http.File
	}

	func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
	}

	func main() {
		fs := justFilesFilesystem{http.Dir("/tmp/")}
		http.ListenAndServe(":8080", http.FileServer(fs))
	} */



