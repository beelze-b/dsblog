package article

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"strings" //remove blank identifier to remove unused compiler error
	"time"
)

type Comment struct {
	Author  string
	Date    string
	Content string //content should be template style: see documentation details on golang site
}
type Article struct {
	Title          string
	Url            string
	Author         string
	Date           time.Time
	Tags           []string
	Content        template.HTML //content should be template style: see documentation details on golang site
	LimitedContent string        // does not need to be type template.HTML because SearchResults and Aggregator Unescape
	Comments       []Comment
}

func SaveJSONArticle(a Article) {
	b, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
		return
	}
	title := strings.ToLower(strings.Replace(a.Title, " ", "_", -1)) + ".html"
	filePath := "src/static/articles/" + title
	// 0644 means overwrite
	ioutil.WriteFile(filePath, b, 0644)
	return
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
func LoadArticleTitle(title string) (Article, error) {
	// must first parse title to find file and then create article out of it
	filePath := "src/static/articles/" + strings.ToLower(strings.Replace(title, " ", "_", -1))
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Article{}, err
	}
	var article Article
	err = json.Unmarshal(b, &article)

	if err == nil {
		return article, nil
	} else {
		return Article{}, errors.New("Article not found")
	}
}

func LoadArticleFilePath(filePath string) (Article, error) {
	// must first parse title to find file and then create article out of it
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Article{}, err
	}
	var article Article
	err = json.Unmarshal(b, &article)

	if err == nil {
		return article, nil
	} else {
		return Article{}, errors.New("Article not found")
	}
}
func (a *Article) AddComment(author string, date string, comment string) {
	var entry = Comment{author, date, comment}
	a.Comments = append(a.Comments, entry)
}

/**
func main() {
	Title := "Hello"
	Url := "/article/hello.html"
	Author := "Nabeel"
	Tags := []string{"pew", "miracle"}
	Date := time.Now()
	Content := "This is the content"
	Comments := []string{"comment 1", "comment 2"}
	article := Article{Title, Url, Author, Date, Tags, Content, Comments}
	SaveJSONArticle(article)
}
*/
