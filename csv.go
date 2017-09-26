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

	type pkeytype struct {
		pkey, pmap  int
	}
	var pkey []pkeytype
	var tempmap map[string]int
	tempmap = make(map[string]int)
	items := reflect.ValueOf(get_csv_to).Elem()
	argtype :=items.Type()
	structtype:=argtype.Elem()

	csvFile, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
		}
	r := csv.NewReader(bufio.NewReader(csvFile))
	r.Comma = ';'
	lineheader, err := r.Read()
	if err != nil {
		log.Fatal(err)
		}
	for i, num:= range lineheader {
		tempmap[num]=i
	}
	fmt.Println(lineheader)
	fmt.Println(structtype)

	for i := 0; i < structtype.NumField(); i++ {
		if a,ok:=structtype.Field(i).Tag.Lookup("csv");ok {
			b,ok:=tempmap[a]
			if ok != true {
				log.Fatal(err)
			}
			pkey=append(pkey,pkeytype{i,b})
			fmt.Println(a)
		}
		//fmt.Println("===",v.Type().Field(j).Name, v.Field(j).Interface())
	}
	fmt.Println(pkey)
	fmt.Println(argtype)
	items.Set(reflect.MakeSlice(argtype,1,1000))
	var people []good
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		good1:= reflect.New(structtype)
		good2:=good1.Elem()
		for _,k:=range pkey {
		 good2.Field(k.pkey).SetString(line[k.pmap])
		}
		//people = append(people, good1)
		fmt.Println(line)
		items.Set(reflect.Append(items,good2))
	}

	fmt.Println(people)





	if items.Kind() == reflect.Slice {
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			if item.Kind() == reflect.Struct {
				v := reflect.Indirect(item)
				for j := 0; j < v.NumField(); j++ {
					//v.Field(j).SetString("1039")
					fmt.Println("===",v.Type().Field(j).Name, v.Field(j).Interface())
				}
			}
		}
	}
	items.Set(reflect.AppendSlice(items,reflect.ValueOf(people)))

}