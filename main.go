package main

import (
    "github.com/NabeelSarwar/dsblog/src/article"
    "html/template"
    "net/http"
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

func main() {
  http.HandleFunc("/", HomePageFunc)
  fs := http.FileServer(http.Dir("src/static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))

}
