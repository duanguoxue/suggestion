package main

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "time"
    "fmt"
    "runtime"
)

func segContent( seg string )([]byte ,error) {
    resp, err := http.PostForm("http://127.0.0.1:8080/suggestion",url.Values{"content": {seg}})
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    dy, err := ioutil.ReadAll(resp.Body)
    return dy, nil
}

func main(){
    runtime.GOMAXPROCS(1)
    count := 10000
    cur := time.Now()
    for i := 0; i< count; i++ {
        segContent("我是")
    }
    fmt.Println(time.Since(cur))
    fmt.Println(count)
}
