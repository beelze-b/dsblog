package main

import (
	"fmt"
	"github.com/NabeelSarwar/dsblog/src/article"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"time"
)

var agg = article.Aggregate()

func ArticleTemplate(w http.ResponseWriter, r *http.Request, a article.Article) {
	t, err := template.ParseFiles("src/static/html/templatepost.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, a)
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
	sender := r.FormValue("name")
	email := r.FormValue("email")
	subject := r.FormValue("subject")
	host := "aspmx.l.google.com"
	body := r.FormValue("message")

	header := make(map[string]string)
	header["From"] = (&mail.Address{sender, email}).String()
	header["To"] = (&mail.Address{"Nabeel Sarwar", "nabeelsarwar200@gmail.com"}).String()
	header["Subject"] = subject

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth("", "nabeelsarwar200@gmail.com", "password", host)
	err := smtp.SendMail(host+":25", auth, email, []string{"nabeelsarwar200@gmail.com"}, []byte(message))
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/contact.html", http.StatusFound)
}

func SearchBarHandler(w http.ResponseWriter, r *http.Request) {
	searchTerms := r.FormValue("searchbar")
	searchResults := article.SearchResults(searchTerms)
	t, err := template.ParseFiles("src/static/html/search_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, searchResults)

}

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	articleUrl := "/fakeurl.html"
	author := r.FormValue("name")
	date := time.Now().String()
	comment := r.FormValue("comment")
	// fix this loading
	art, err := article.LoadArticleTitle(articleUrl[len("/article/") : len(articleUrl)-len(".html")])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	(&art).AddComment(author, date, comment)
	http.Redirect(w, r, articleUrl, http.StatusFound)

}

func HomePageFunc(w http.ResponseWriter, r *http.Request) {
	// load main page and print it
	t, err := template.ParseFiles("src/static/html/mainpage.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, agg)
}

func AboutPageFunc(w http.ResponseWriter, r *http.Request) {
	AboutPage, err := ioutil.ReadFile("src/static/html/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", AboutPage)
}

func ContactPageFunc(w http.ResponseWriter, r *http.Request) {
	ContactPage, err := ioutil.ReadFile("src/static/html/contact.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", ContactPage)
}

func main() {

	Title := "Hello"
	Url := "hello.html"
	Author := "Nabeel"
	Tags := []string{"pew", "miracle"}
	Date := time.Now()
	LimitedContent := "Limited Content"
	Content := "This is the content full"
	Comments := []article.Comment{article.Comment{"John", "October 15, 2016", "Good stuff"},
		article.Comment{"Michael", "November 15, 2016", "Bad stuff"}}
	art := article.Article{Title, Url, Author, Date, Tags, Content, LimitedContent, Comments}
	article.SaveJSONArticle(art)

	http.HandleFunc("/", HomePageFunc)
	http.HandleFunc("/emailme", sendEmailHandler)
	http.HandleFunc("/search", SearchBarHandler)
	http.HandleFunc("/about.html", AboutPageFunc)
	http.HandleFunc("/contact.html", ContactPageFunc)
	http.HandleFunc("/article/", articleHandler)
	fs := http.FileServer(http.Dir("src/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
