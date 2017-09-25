package main

import (
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"log"
	"fmt"
	"reflect"
)

func parse2(fname string,fields []string){
	//"list.csv"
	csvFile, _ := os.Open(fname)
	r := csv.NewReader(bufio.NewReader(csvFile))
	r.Comma = ';'
	var people []good
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		people = append(people, good{
			Articul: line[0],
			Info:  line[1],
			Image: line[2],
		})
		fmt.Println(line)
	}
	fmt.Println(people)


}

func dump(datasets interface{}, fname string) {
	items := reflect.ValueOf(datasets)
	if items.Kind() == reflect.Slice {
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			if item.Kind() == reflect.Struct {
				v := reflect.Indirect(item)
				for j := 0; j < v.NumField(); j++ {
					fmt.Println(v.Type().Field(j).Name, v.Field(j).Interface())
				}
			}
		}
	}
}

func parse(get_csv_to interface{}, fname string) {

	csvFile, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
		}
	r := csv.NewReader(bufio.NewReader(csvFile))
	r.Comma = ';'

	items := reflect.ValueOf(get_csv_to).Elem()
	itemstype :=items.Type()
	fmt.Println(itemstype)
	items.Set(reflect.MakeSlice(itemstype,1,1000))
	var people []good
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		people = append(people, good{
			Articul: line[0],
			Info:  line[1],
			Image: line[2],
		})
		//fmt.Println(line)
	}

	fmt.Println(people)





	if items.Kind() == reflect.Slice {
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			if item.Kind() == reflect.Struct {
				v := reflect.Indirect(item)
				for j := 0; j < v.NumField(); j++ {
					v.Field(j).SetString("1039")
					fmt.Println("===",v.Type().Field(j).Name, v.Field(j).Interface())
				}
			}
		}
	}
	items.Set(reflect.AppendSlice(items,reflect.ValueOf(people)))

}