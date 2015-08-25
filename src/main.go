package main

import (
    "flag"
    "html/template"
    "log"
    "net/http"
)


func ArticleTemplate(w http.ResponseWriter, r *http.Request, a Article) {
  t, _ := template.ParseFiles("/static/html/articleTemplate.html")
  t.Execute(w, a)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/article/"):]
    article, err := LoadJSONArticle(title)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    ArticleTemplate(w, r, article)

}

func main() {
  http.handleFunc('/', homePageFunc)
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))

}
