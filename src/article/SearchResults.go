package article

import (
	"bytes"
	"golang.org/x/text/language"
	"golang.org/x/text/search"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

type SearchResults struct {
	RelevantArticles []Article
}

/**
Takes a matcher and a search term and sees if that search term is in any of the matches
*/
func UseMatcher(matcher *search.Matcher, searchTerm string, fileName string) bool {
	fileContent, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Panic(err)
	}
	start, _ := matcher.Index(fileContent, []byte(searchTerm))
	return start != -1
}

func NewSearchResults(searchTermsString string) SearchResults {
	var articleDirectory string = filepath.Dir("src/static/articles/")
	articlesArray := make([]Article, 10)

	if articleDirectory == "." {
		log.Panic("Directory not found")
	}

	searchTerms := strings.Split(searchTermsString, " ")

	files, err := ioutil.ReadDir(articleDirectory)
	if err != nil {
		log.Panic(err)
	}

	matcher := search.New(language.English)
	for _, file := range files {
		for _, searchTerm := range searchTerms {
			validArticle := UseMatcher(matcher, searchTerm, file.Name())
			if validArticle {
				article, err := LoadArticleFilePath(file.Name())
				if err != nil {
					log.Panic(err)
				}
				articlesArray = append(articlesArray, article)
			}
			// do not want to add an article twice
			break
		}
	}
	return SearchResults{articlesArray}
}

/**
Very similar to the use in aggregator.go
*/
func DisplaySearchResult(a Article) template.HTML {
	var url = a.Url
	var display = template.HTML(`<article><h2> <a href="/article/` + url + `">` + a.Title + `</a> </h4>
                        <div class="row">
                            <div class="group1 col-sm-6 col-md-6">
                                <span class="glyphicon glyphicon-folder-open"></span>  <a href="#">Signs</a>
                                <span class="glyphicon glyphicon-bookmark"></span> <a href="#">Aries</a>,
                                <a href="#">Fire</a>, <a href="#">Mars</a>
                            </div>
                            <div class="group2 col-sm-6 col-md-6">
                                <span class="glyphicon glyphicon-pencil"></span> <a href="/article/` + url + `#comments">` +
		strconv.Itoa(len(a.Comments)) + ` Comments</a>  
								<span class="glyphicon glyphicon-time"></span>` + a.Date.String() + `
                            </div>
                        </div>
                        <hr>

                        <br />
                        <p>` + a.LimitedContent + `</p>
                        <p class="text-right">
                        <a href="/article/` + url + `"class="text-right">
                            continue reading...
                        </a>
                        </p>
                        </hr></article>`)
	return display

}

func (results SearchResults) DisplaySearchResults() template.HTML {
	var buffer bytes.Buffer
	for _, value := range results.RelevantArticles {
		buffer.WriteString(string(DisplaySearchResult(value)))
	}
	return template.HTML(buffer.String())
}