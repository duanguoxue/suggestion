package main

import (
	"fmt"
	"github.com/duanguoxue/suggestion"
)

func queryDict(s *suggestion.Search, querystr string, count int) {
	mdata := s.SearchSuggest(querystr)
	for i, data := range mdata {
		if i > count {
			break
		}
		fmt.Println(data.Term, data.Weight)
	}
}

func main() {
	s := &suggestion.Search{}
	s.Init("../data/pinyin-utf8.dat", "../data/dict.txt")
	queryDict(s, "有限", 10)
	queryDict(s, "youxian", 10)
	queryDict(s, "yx", 10)
	queryDict(s, "bzj", 10)
}
