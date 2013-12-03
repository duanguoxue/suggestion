//实现简单网络调用接口
//code from <<Go 语言编程>> photoweb.go

package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"os/signal"
	"path"
	"runtime"
	"runtime/debug"
	//"runtime/pprof"
	"github.com/duanguoxue/suggestion"
	"time"
)

const (
	PINYIN_DICT_DIR = "../data/pinyin-utf8.dat"
	QUERY_DIR       = "../data/dict.txt"
	TEMPLATE_DIR    = "../views"
)

var templates (map[string]*template.Template)
var search suggestion.Search

//simple web interface
func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				//user-defined error 505 function
				//w.WriteHeader(http.StatusInternalServerError)
				//renderHtml(w, "error", e)
				log.Println("WARN: panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}

func renderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) (err error) {
	err = templates[tmpl].Execute(w, locals)
	return err
}

func suggestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := renderHtml(w, "suggestion.html", nil)
		check(err)
	}
	if r.Method == "POST" {
		//Must ParseForm 否则取不到Form数据
		r.ParseForm()
		data := r.Form["content"]
		log.Println("value:", r.Form["content"])
		spdata := search.SearchSuggest(data[0])
		result := ""
		for _, val := range spdata {
			result += val.Term + "/ "
			if val.Weight < 0 {
				break
			}
			log.Println(val.Term, val.Weight)
		}
		fmt.Fprintf(w, result)
	}
}

func spellHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := renderHtml(w, "spellcheck.html", nil)
		check(err)
	}
	if r.Method == "POST" {
		//Must ParseForm 否则取不到Form数据
		r.ParseForm()
		data := r.Form["content"]
		spdata := search.SearchSpell(data[0])
		result := ""
		for _, val := range spdata {
			if val.Weight < 1000 {
				break
			}
			log.Println(val.Term, val.Weight)
			result += val.Term + "/ "
			//break
		}
		//log.Println(result)
		fmt.Fprintf(w, result)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	if err != nil {
		panic(err)
	}
	templates = make(map[string]*template.Template)
	for _, fileInfo := range fileInfoArr {
		tempname := fileInfo.Name()
		if ext := path.Ext(tempname); ext != ".html" {
			continue
		}
		log.Println(TEMPLATE_DIR, tempname)
		t := template.Must(template.ParseFiles(TEMPLATE_DIR + "/" + tempname))
		templates[fileInfo.Name()] = t
	}
}

func getLogFileName() string {
	t := time.Now()
	os.Mkdir("log", 0644)
	logfilepath := "log/sug_" + t.Format(time.RFC3339) + ".log"
	return logfilepath
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logf, err := os.OpenFile(getLogFileName(), os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Fatalln(err)
	}
	defer logf.Close()
	log.SetOutput(io.MultiWriter(logf, os.Stdout))

	/*
		    // using parse performence
			profile := true
			if profile {
				f, _ := os.Create("profile.cpu")
				pprof.StartCPUProfile(f)
			}

			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)
				<-c

				if profile {
					pprof.StopCPUProfile()
				}

				fmt.Println("Caught interrupt.. shutting down.")
				os.Exit(0)
			}()
	*/
	search = suggestion.Search{}
	search.Init(PINYIN_DICT_DIR, QUERY_DIR)

	http.HandleFunc("/spellcheck", safeHandler(spellHandler))
	http.HandleFunc("/suggestion", safeHandler(suggestHandler))
	log.Println("provide suggestion and spellcheck base functions")
	log.Println("web port is 8080,visit: http://ip:8080/suggestion")
	log.Println("web port is 8080,visit: http://ip:8080/spellcheck")
	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err.Error())
	}
}
