package article

import(
	"time"
	_ "text/template" //remove blank identifier to remove unused compiler error
	"io/ioutil"
	_ "strings" //remove blank identifier to remove unused compiler error
	"net/url"
	"strconv"
)

type Article struct{
	Title string
	Author string
	Date time.Time
	Tags []string
	Content []byte //content should be template style: see documentation details on golang site
	Views int //this information might have to handled by a database
	UniqueViews int //this information might need to be handled by a database
}

/*
Creates a pointer to an article. Article format must have format above.
Parsing isn't set up, but it could be done using this format:
Title:
Author:
Date:
	date
tags:
	tag1, tag2, tag3
Content:
	Content
Detail to change as no articles are written yet.
*/
func LoadArticle(articleId int) *Article{
	/*Use mgo (MongoDB bindings for Golang) to load article*/
	content, err := ioutil.ReadFile(Title + ".txt")
	if err != nil {
		return new(Article)
	}
}

/*
Creates a URL version of an article's title. Signatures use date and title for uniqueness in URL.
The URL is not totally valid, but the handler will enable its use.
*/
func (a *Article) parseTitle() string{
	date := "/" + strconv.Itoa(a.Date.Year()) + "/" + strconv.Itoa(int(a.Date.Month())) + "/" + strconv.Itoa(a.Date.Day())
	return url.QueryEscape(date + "/" + a.Title + ".html")
}
