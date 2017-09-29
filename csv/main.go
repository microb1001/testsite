package mycsv

import (
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"log"
	"reflect"
	"fmt"
	"strings"
	"strconv"
)

func Load_csv (load_csv_to interface{}, fname, prefix string) {
// load_csv_to это обязательно слайс структуры с произвольными полями.
// поля, которые нужно заполнить из файла должны иметь тег "csv:имя столбца в csv"
// остальные получат нулевые значения

	type csv_field_type struct {
		fldNum, mapNum int
	}
	var csvField []csv_field_type

	var csvHeadersMap map[string]int = make(map[string]int)

	items := reflect.Indirect(reflect.ValueOf(load_csv_to)) // т.к. передаем ссылку
	struct_type :=items.Type().Elem()

	csvFile, err := os.Open(fname)
	if err != nil {	log.Fatal("Файл: ",fname," не найден. Ошибка: ",err)	}

	r := csv.NewReader(bufio.NewReader(csvFile))
	r.Comma = ';'
	csvHeaders, err := r.Read()
	if err != nil {log.Fatal("Нет заголовков. Файл: ",fname," Ошибка: ",err)}

	for i, num:= range csvHeaders {
		csvHeadersMap[num]=i
	}

	// Создание в csvField соответствия столбец файла - поле структуры
	for i := 0; i < struct_type.NumField(); i++ {
		if a,ok:= struct_type.Field(i).Tag.Lookup(prefix);ok { // для полей с нужным тегом
			b,ok:= csvHeadersMap[a]
			if ok != true {	log.Fatal("Нет столбца csv с именем: ",a," Ошибка: ",err)}

			csvField =append(csvField, csv_field_type{i,b})
		}
	}

	// Заполнение массива структур
	items.Set(reflect.MakeSlice(items.Type(),0,1000))

	count:=2
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {log.Fatal("Не удалось обработать строку: ", count," Ошибка: ",err)}

		newItem :=reflect.Indirect(reflect.New(struct_type))

		for _,k:=range csvField {
			// есть Convert(t Type) Value но реализован в языке недавно!
			switch Tp := newItem.Field(k.fldNum).Kind(); Tp {
			case reflect.String:
				newItem.Field(k.fldNum).SetString(line[k.mapNum])
			case reflect.Bool:
				s,err:=strconv.ParseBool(line[k.mapNum])
				if err != nil {log.Fatal("Не удалось конвертировать string в ", Tp, " Строка: ", count," Ошибка: ",err)}
				newItem.Field(k.fldNum).SetBool(s)
			case reflect.Int, reflect.Int32, reflect.Int64:
				s,err:=strconv.ParseInt(line[k.mapNum],10,32)
				if err != nil {log.Fatal("Не удалось конвертировать string в ", Tp, " Строка: ", count," Ошибка: ",err)}
				newItem.Field(k.fldNum).SetInt(s)
			case reflect.Float32,reflect.Float64:
				s,err:=strconv.ParseFloat(line[k.mapNum],32)
				if err != nil {log.Fatal("Не удалось конвертировать string в ", Tp, " Строка: ", count," Ошибка: ",err)}
				newItem.Field(k.fldNum).SetFloat(s)
			case reflect.Slice: // Баг: только []string, не []int например, иначе  panic
				newItem.Field(k.fldNum).Set(reflect.ValueOf(strings.Split(line[k.mapNum],"|")))
			default:
				log.Fatal("Не могу преобразовать строку в тип: ", Tp," Не реализовано.")
			}
		}
		items.Set(reflect.Append(items, newItem))
		count++
	}
}

// Чужая функция для понимания работы reflect
func Dump(datasets interface{}) {
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

// пригодится строка определения типа переменной
// Interface() показывает базовый тип а не value
//fmt.Printf("type of ms: %T\n", Item.Interface())

//c := intPtr2.Elem().Interface().(int) преобразование типов


// Пример парсинга csv
// Не используется удалить
/*
func parse(fname string,fields []string){

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
*/