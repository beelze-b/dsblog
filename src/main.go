package main

import (
    "flag"
    "html/template"
    "log"
    "net/http"
)




func main() {
  http.handleFunc('/', homePageFunc)
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))
}
