package escape

import "bytes"

//
// фильтр символов, второй аргумент функция фильтрования
// rune это int32
//

func Go(s_in []byte, F func(r rune) rune) (s_out []byte) {
	return bytes.Map(F, s_in)
}

var Rus [][]rune = [][]rune{[]rune{'а','я'},[]rune{'А','Я'}}
var Digit [][]rune = [][]rune{[]rune{'0','9'}}
var Eng [][]rune = [][]rune{[]rune{'a','z'},[]rune{'A','Z'}}

func FnEng09(r rune) rune {
	if (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') {
		return r
	}
	return -1
}

func Fn09(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}

func FnNotSpecial(r rune) rune { // не рабочий исходный пример
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

func FnRus09(r rune) rune {
	if (r >= 'а' && r <= 'я') ||
		(r >= 'А' && r <= 'Я') ||
		(r >= '0' && r <= '9') {
		return r
	}
	return -1
}

func MakeFn(Allowed [][]rune) func(rune) rune {
	return func(r rune) (ret rune) {
		for _,k:=range Allowed{
			if len(k)==2&&r>=k[0]&&r<=k[1]{
				return r
			}
		}
		return -1
	}
}




