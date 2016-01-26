package article

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Aggregator struct {
	articles   []Article
	TitleToUrl map[string]string
	UrlToTitle map[string]string
}

// This function aggregates all the functions and creates an Aggregator object that can be
// used for templating on the main page.
func Aggregate() Aggregator {
	files, err := ioutil.ReadDir("src/static/")

	if err != nil {
		log.Panic(err)
	}

	for i, v := range files {
		// load the json from the file
		// these things are file values
		if v.IsDir() {
			continue
		}
		var art Article
		fileContent, err := ioutil.ReadFile("src/static/" + v.Name())

		if err != nil {
			log.Panic(err)
		}

		err = json.Unmarshal(fileContent, &art)
		if err != nil {
			log.Panic(err)
		}

	}

}