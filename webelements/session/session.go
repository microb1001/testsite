package session

import (
	"net/http"
	"strconv"
	"../servo"
)

const COOKIES_NAME="PHPSESSID"

// Получает id сессии в виде uint64, если сессии нет, создает
func Get(w http.ResponseWriter, r *http.Request) uint64 {
	var id *http.Cookie
	var res uint64
	id, err := r.Cookie(COOKIES_NAME)
	if err == nil { // кука
		res, err = strconv.ParseUint(id.Value, 10, 64)
		if err == nil {
			return res
		}
	}
	res, err = strconv.ParseUint(r.FormValue(COOKIES_NAME), 10, 64)
	if err == nil { // GET параметр
		return res
	}
	res = servo.Rand64() // создать
	http.SetCookie(w, &http.Cookie{Name: COOKIES_NAME, Value: strconv.FormatUint(res, 10)})
	return res
}
