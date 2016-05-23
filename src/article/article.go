package article

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strings"         //remove blank identifier to remove unused compiler error
	_ "text/template" //remove blank identifier to remove unused compiler error
	"time"
)

type Article struct {
	Title    string
	Url      string
	Author   string
	Date     time.Time
	Tags     []string
	Content  []byte //content should be template style: see documentation details on golang site
	Comments []string
}

func SaveJSONArticle(a Article) {
	b, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
		return
	}
	title := strings.ToLower(strings.Replace(a.Title, " ", "_", -1))
	filePath := "../static/articles/" + title
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
	filePath := "../static/articles/" + strings.ToLower(strings.Replace(title, " ", "_", -1))
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
	var entry = `<li class="comment">
                                <div class="clearfix">
                                    <h4 class="pull-left">` + author + `</h4>
                                    <p class="pull-right">` + date + `</p>
                                </div>
                                <p>
                                <em>` + comment + `</em>
                                </p>
                            </li>`
	a.Comments = append(a.Comments, entry)
}

/**
func main() {
	Title := "HEllo"
	Url := "/article/test.html"
	Author := "Nabeel"
	Tags := []string{"pew", "miracle"}
	Date := time.Now()
	Content := []byte("This is the content")
	Comments := []string{"comment 1", "comment 2"}
	article := Article{Title, Url, Author, Date, Tags, Content, Comments}
	SaveJSONArticle(article)
}
*/
