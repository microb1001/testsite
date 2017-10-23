package sphinx

import (
	"index/suffixarray"
	"bytes"
	"sort"
)

type Sphinx struct{
Data  []byte
index *suffixarray.Index
key []int
}

func (s *Sphinx) Init(){
//s.key=make(map[int]int,0)
}

type AddToSphinx interface {
ToByte(int) []byte
Len() int
}

func (s *Sphinx) Add(w AddToSphinx){
for i:=0;i<w.Len();i++{
s.Data =append (s.Data, 0)
s.Data =append (s.Data, bytes.ToLower(w.ToByte(i))...)
s.key=append(s.key,len(s.Data))
}
s.Data =append (s.Data, 0)
s.index = suffixarray.New(s.Data)
}

func (s *Sphinx) Find(str string) (ret []int) {
tempRet := s.index.Lookup(bytes.ToLower([]byte(str)), -1)
sort.Ints(tempRet)
var i, j int = 0, 0
var out bool = true // Повторы убрать
for i = 0; (i < len(tempRet)) && (j < len(s.key)); {
if tempRet[i] < s.key[j] {
if out {ret = append(ret, j)}
out=false
i++
} else {
out=true
j++

}
}
return ret
}