package article

import (
	"golang.org/x/text/search"
	"io/ioutil"
	"log"
	"path/filepath"
)

type SearchResults struct {
	RelevantArticles []Article
}

func SearchResults(searchTerms string) SearchResults {
	var articleDirectory string = filepath.Dir("src/static/article/")

	if articleDirectory == "." {
		log.Panic("Directory not found")
	}

	files, err := ioutil.ReadDir(articleDirectory)

}
