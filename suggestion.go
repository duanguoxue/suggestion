//suggestion 功能
//通过查询词频率词典,实现汉字,前缀拼音,简拼的 suggestion
//同时可以使用汉字转拼音,对同音字进行简单纠正 spellcheck
//内部使用go map实现 内存开销比较大
package suggestion

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"sort"
	"strconv"
	"unicode/utf8"
)

type PinyinSet map[string]bool

//implement Word sort
type WordSorter []*Word
type ByWeight struct{ WordSorter }

type Word struct {
	Term   string
	Weight int
}

func (s WordSorter) Len() int      { return len(s) }
func (s WordSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

//func (s ByWeight) Less(i, j int) bool   { return s.WordSorter[i].Weight > s.WordSorter[j].Weight }
func (s WordSorter) Less(i, j int) bool { return s[i].Weight > s[j].Weight }
func (s WordSorter) Sort()              { sort.Sort(s) }

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Lexicon struct {
	Lex map[rune][]string
}

//汉字转拼音 字典
func (lex *Lexicon) Load(path string) {
	if len(lex.Lex) > 0 {
		return
	}
	lex.Lex = make(map[rune][]string)
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	bio := bufio.NewReader(f)
	for {
		line, err := bio.ReadBytes('\n')
		if err != nil {
			break
		}
		line = bytes.Trim(line, "\n")
		lines := bytes.Split(line, []byte{0x20})
		ru, _ := utf8.DecodeRune(lines[0])
		lines = lines[1:]
		slen := len(lines)
		pinyinset := make(PinyinSet)
		for i := 0; i < slen; i++ {
			ilen := len(lines[i])
			if lines[i][ilen-1] >= '0' && lines[i][ilen-1] <= '9' {
				lines[i] = lines[i][0 : ilen-1]
				pinyinset[string(lines[i])] = true
			}
		}
		pys := make([]string, len(pinyinset))
		i := 0
		for k, _ := range pinyinset {
			pys[i] = k
			i++
		}
		lex.Lex[ru] = pys
	}
	log.Printf("lexicon: %d", len(lex.Lex))
}

//拼音切片 flag = 0 前缀全拼, 1 前缀简拼
func (lex *Lexicon) GetPinyinKey(s string, flag int) PinyinSet {
	runes := []rune(s)
	allpinyin := make(PinyinSet)
	tpy := make(PinyinSet)
	tpy[""] = true
	for _, v := range runes {
        pinyins, ok := lex.Lex[v]
        if ok !=true {//其他数字或者字符
            pinyins = append(pinyins, string(v))
        }
        tmpapy := tpy
        tpylocal := make(PinyinSet)
        for _, pk := range pinyins {
            for k, _ := range tmpapy {
                if flag == 0 {
                    k = k + pk
                }
                if flag == 1 {
                    k = k + string([]byte(pk)[0])
                }
                allpinyin[k] = true
                tpylocal[k] = true
            }
        }
        if len(allpinyin) > 40 {
            break
        }
        tpy = tpylocal
	}
	//log.Printf("%c  %s\r\n", runes, allpinyin)
	return allpinyin
}

//汉字转拼音 检测是否为拼音和数字
func (lex *Lexicon) ConvertPinyin(s string) string {
	runes := []rune(s)
	var pinyin string
	for _, v := range runes {
		pinyins, ok := lex.Lex[v]
		if ok == false {
			if ('a' <= v && v <= 'z') || ('0' <= v && v <= '9') || ('A' <= v && v <= 'Z') {
				pinyins = append(pinyins, string(v))
			} else {
				continue
			}
		}
		for _, pk := range pinyins {
			pinyin += pk
			break
		}
	}
	return pinyin
}

//load 短语-查询频率 词典
func LoadQueryDict(path string, dict map[string]int) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	bio := bufio.NewReader(f)
	for {
		line, err := bio.ReadBytes('\n')
		if err != nil {
			break
		}
		line = bytes.Trim(line, "\n")
		lines := bytes.Split(line, []byte{'#', '#', '#'})
		ivalue, _ := strconv.Atoi(string(lines[1]))
		if len(lines[0]) != 0 && len(lines[1]) != 0 {
			dict[string(lines[0])] = ivalue
		}
	}
	log.Printf("query key statistics: %d", len(dict))
}

type Search struct {
	Dict map[string]WordSorter
	Lex  *Lexicon
}

//前缀 汉字切片
func getSpliceZhKey(s string) []string {
	runes := []rune(s)
	splicestrs := make([]string, len(runes))
	for i, _ := range runes {
		splicestrs[i] = string(runes[0:i+1])
	}
	return splicestrs
}

//处理词频词典到suggestion 字典
func (s *Search) procQueryFreqKeys(freqdict map[string]int, lexicon *Lexicon) {
	for ss, v := range freqdict {
		w := Word{ss, v}
		//前缀拼音 quan pin yin
        pysp := lexicon.GetPinyinKey(ss, 0)
		//前缀简拼 jian pin yin
		pysimsp := lexicon.GetPinyinKey(ss, 1)
        //前缀汉字 han zi dict
		spzh := getSpliceZhKey(ss)

		for _, key := range spzh {
            pysp[key] = true
		}
		for key, _ := range pysimsp {
			//判断是否与前缀拼音 有相同的key,丢弃重复key 例如:阿等单音字
			if _, ok := pysp[key]; ok == false {
                pysp[key] = true
			}
		}

		for key, _ := range pysp {
			s.Dict[key] = append(s.Dict[key], &w)
		}
	}
	//sort by weight
	for _, words := range s.Dict {
		//sort.Sort(ByWeight{words})
		//sort.Sort(words) //error why  || go type
		words.Sort()
	}
    //key := "101ye"
    //log.Println(key,s.Dict[key][0].Term)
}

func (s *Search) Init(pypath, querypath string) {
	if len(s.Dict) != 0 {
		return
	}
	s.Dict = make(map[string]WordSorter)
	resch := make(chan bool, 1)
	lexicon := new(Lexicon)
	go func() {
		lexicon.Load(pypath)
		resch <- true
	}()
	querydict := make(map[string]int)
	go func() {
		LoadQueryDict(querypath, querydict)
		resch <- true
	}()
	<-resch
	<-resch
	s.procQueryFreqKeys(querydict, lexicon)
	s.Lex = lexicon
}

func (s *Search) SearchSuggest(term string) []*Word {
	return s.Dict[term]
}

func (s *Search) SearchSpell(term string) []*Word {
	key := s.Lex.ConvertPinyin(term)
	if key == "" {
		key = term
	}
	return s.Dict[key]
}
