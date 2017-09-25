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
