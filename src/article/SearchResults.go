package article

import (
	"golang.org/x/text/language"
	"golang.org/x/text/search"
	"io/ioutil"
	"log"
	"path/filepath"
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
	var articleDirectory string = filepath.Dir("src/static/article/")
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
	for _, searchTerm := range searchTerms {
		for _, file := range files {
			validArticle := UseMatcher(matcher, searchTerm, file.Name())
			if validArticle {
				article, err := LoadArticleFilePath(file.Name())
				if err != nil {
					log.Panic(err)
				}
				articlesArray = append(articlesArray, article)
			}
		}
	}
	return SearchResults{articlesArray}
}
