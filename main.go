package main

import (
    "github.com/NabeelSarwar/dsblog/src/article"
    "html/template"
    "net/http"
    "io/ioutil"
    "fmt"
)


func ArticleTemplate(w http.ResponseWriter, r *http.Request, a article.Article) {
  t, _ := template.ParseFiles("src/static/html/articleTemplate.html")
  t.Execute(w, a)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/article/"):]
    article, err := article.LoadJSONArticle(title)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    ArticleTemplate(w, r, article)
}

func HomePageFunc(w http.ResponseWriter, r *http.Request) {
  // load main page and print it
  mainPage, err := ioutil.ReadFile("src/static/html/mainpage.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  fmt.Fprintf(w, "%s", mainPage)
}

func main() {
  http.HandleFunc("/", HomePageFunc)
  fs := http.FileServer(http.Dir("src/static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))
  http.ListenAndServe(":8080", nil)
}
