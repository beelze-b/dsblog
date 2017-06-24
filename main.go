package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	// "github.com/beelzebud/dsblog/article" //when using go build
	// "time"
	"article" // when using goapp serve
)

var agg = article.Aggregate()
var templates = template.Must(template.ParseFiles(
	"static/html/templatepost.html", "static/html/search_page.html", "static/html/mainpage.html"))

func ArticleTemplate(w http.ResponseWriter, r *http.Request, a article.Article) {
	err := templates.ExecuteTemplate(w, "template_post", a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/article/"):]
	article, err := article.LoadArticleTitle(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ArticleTemplate(w, r, article)
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
}

func SearchBarHandler(w http.ResponseWriter, r *http.Request) {
	searchTerms := r.FormValue("searchbar")
	searchResults := article.NewSearchResults(searchTerms)
	fmt.Println(searchTerms)
	fmt.Println(searchResults)
	err := templates.ExecuteTemplate(w, "search_page", searchResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func HomePageFunc(w http.ResponseWriter, r *http.Request) {
	// load main page and print it
	err := templates.ExecuteTemplate(w, "main_page", agg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AboutPageFunc(w http.ResponseWriter, r *http.Request) {
	AboutPage, err := ioutil.ReadFile("static/html/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", AboutPage)
}

func StaticRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://storage.googleapis.com/dsblog-158823.appspot.com" + r.URL.Path, 303)
}

func init() {

	/**
	Title := "Hello"
	Url := "hello.html"
	Author := "Nabeel"
	Tags := []string{"pew", "miracle"}
	Date := time.Now()
	LimitedContent := "Limited Content"
	Content := []byte("This is the content full")
	art := article.Article{Title, Url, Author, Date, Tags, Content, LimitedContent}
	article.SaveJSONArticle(art)
	*/


	http.HandleFunc("/", HomePageFunc)
	http.HandleFunc("/emailme", sendEmailHandler)
	http.HandleFunc("/search", SearchBarHandler)
	http.HandleFunc("/about.html", AboutPageFunc)
	http.HandleFunc("/article/", articleHandler)
	http.HandleFunc("/static/", StaticRedirectHandler)
	//	http.ListenAndServe(":8080", nil)
}
