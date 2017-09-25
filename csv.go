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




	items := reflect.ValueOf(get_csv_to)
	if items.Kind() == reflect.Slice {
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			if item.Kind() == reflect.Struct {
				v := reflect.Indirect(item)
				for j := 0; j < v.NumField(); j++ {
					v.Field(j).SetString("103")
					fmt.Println(v.Type().Field(j).Name, v.Field(j).Interface())
				}
			}
		}
	}
}