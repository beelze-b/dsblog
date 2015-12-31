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
)

func ArticleTemplate(w http.ResponseWriter, r *http.Request, a article.Article) {
	t, _ := template.ParseFiles("src/static/html/templatepost.html")
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
	http.HandleFunc("/emailme", sendEmailHandler)
	fs := http.FileServer(http.Dir("src/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
