package suggestion

import (
	"sort"
	"testing"
)

func TestWordSorter(t *testing.T) {
	wdsort := make(WordSorter, 4)
	wdsort[0] = &Word{"中国人", 100}
	wdsort[1] = &Word{"中国", 15}
	wdsort[2] = &Word{"人", 133}
	wdsort[3] = &Word{"美国人", 1999}
	sort.Sort(wdsort)
	if wdsort[0].Weight != 1999 || wdsort[1].Weight != 133 || wdsort[2].Weight != 100 || wdsort[3].Weight != 15 {
		t.Error("Self define Sort() failed. Got", wdsort, "Expected 1999 133 100 15")
	}
}

func TestLoadPinyin(t *testing.T) {
	lexicon := new(Lexicon)
	lexicon.Load("./data/pinyin-utf8.dat")
	word := []rune("我")
	pinyin, ok := lexicon.Lex[word[0]]
	if ok != true {
		t.Error("pinyin dict key get failed. Got", ok, "Expected true")
	}
	for _, val := range pinyin {
		if val == "wo" {
			return
		}
	}
	t.Error("pinyin value (我=wo) get failed. Got", ok, "Expected true")
}

func TestLoadQueryDict(t *testing.T) {
	querydict := make(map[string]int)
	LoadQueryDict("./data/dict.txt", querydict)
	querycount, ok := querydict["卬头阔步"]
	if ok != true {
		t.Error("pinyin dict key get failed. Got", ok, "Expected true")
	}
	if querycount != 765605927 {
		t.Error("pinyin value (卬头阔步=2377) get failed. Got", querycount, "Expected 2377")
	}
}

func queryDict(s *Search, querystr, except string, t *testing.T, flag int ) {
    var mdata []*Word
    if flag == 0 {
	    mdata = s.SearchSuggest(querystr)
    } else {
	    mdata = s.SearchSpell(querystr)
    }
	if len(mdata) == 0 {
		t.Error("pinyin dict key  get failed. got", querystr, true, "expected true")
	}
	for _, data := range mdata {
		if data.Term == except {
			return
		}
	}
	t.Error("except suggestion get failed. Got", mdata, "Expected ", except)
}

func TestSuggestSpell(t *testing.T) {
	s := &Search{}
	s.Init("./data/pinyin-utf8.dat", "./data/dict.txt")
    //suggestion
	queryDict(s, "卬头阔", "卬头阔步", t, 0)
	queryDict(s, "youxian", "有限公司", t, 0)
	queryDict(s, "yx", "有限公司", t, 0)
	queryDict(s, "bzj", "兵在精而不在多", t, 0)
    //spell 
	queryDict(s, "并在精而不在多", "兵在精而不在多", t, 1)
	queryDict(s, "bingzaijingerbuzaiduo", "兵在精而不在多", t, 1)
}


