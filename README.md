Suggestion
====

最近了解了许多suggestion的相关技术，为巩固新学的Go语言，基于Go实现了一个简单的版本，当前版本内存占用较多，未做过多优化。
支持：汉字，拼音，简拼提示 
启动过程加载字典较慢2-6s

安装更新
====
```
go get -u github.com/duanguoxue/suggestion
```
使用
====
先看一个例子（来自[example.go](/example/example.go)）
```go
package main

import (
	"fmt"
	"github.com/duanguoxue/suggestion/sug"
)

func queryDict(s *sug.Search, querystr string, count int) {
	mdata := s.SearchSuggest(querystr)
	for i, data := range mdata {
		if i > count {
			break
		}
		fmt.Println(data.Term, data.Weight)
	}
}

func main() {
	s := &sug.Search{}
	s.Init("./data/pinyin-utf8.dat", "./data/dict.txt")
	queryDict(s, "有限", 10)
	queryDict(s, "youxian", 10)
	queryDict(s, "yx", 10)
	queryDict(s, "bzj", 10)
}


```
