package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-yaml/yaml"
	"github.com/gosuri/uiprogress"
)

// GetCategoryLinks returns a map of all categories and their links.
func GetCategoryLinks() map[string][]string {
	dat, err := ioutil.ReadFile("category_links.yml")
	if err != nil {
		panic(err)
	}

	var categoryLinks map[string][]string
	err = yaml.Unmarshal(dat, &categoryLinks)
	if err != nil {
		panic(err)
	}

	return categoryLinks
}

// FetchCategoryDocumentLinks returns all links to documents given a sub-category.
func FetchCategoryDocumentLinks(catName string, links []string) []string {
	fmt.Println("\nFetching", len(links), "document links from "+catName+"...")

	var sublinks []string
	for _, link := range links {
		// Fetch webpage
		doc, err := goquery.NewDocument(link)
		if err != nil {
			panic(err)
		}

		// Find links to documents
		doc.Find("#mw-pages .mw-category-group a").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			sublinks = append(sublinks, "https://en.wikipedia.org"+url)
		})
	}

	return sublinks
}

// FetchCategoryDocuments returns all documents given an array of links to documents.
func FetchCategoryDocuments(catName string, links []string) []string {
	//fmt.Println("\nFetching", len(links), "documents from " + cat_name + "...")
	bar := uiprogress.AddBar(len(links))
	bar.AppendElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return catName + " (" + strconv.Itoa(len(links)) + "): " + strconv.Itoa(b.Current())
	})

	var documents []string
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	for i := 0; i < len(links); i++ {
		link := links[i]
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Fetch webpage
			doc, err := goquery.NewDocument(link)
			if err != nil {
				// bar.Incr()
				return
			}

			// Find document content
			content := ""
			doc.Find("#bodyContent").Each(func(i int, s *goquery.Selection) {
				content += s.Text()
			})

			mutex.Lock()
			documents = append(documents, content)
			mutex.Unlock()
			bar.Incr()
		}()
	}

	wg.Wait()
	//fmt.Println("\nFetched", len(documents), "documents from " + cat_name)
	return documents
}

func main() {
	categoryLinks := GetCategoryLinks()
	categorySublinks := make(map[string][]string)
	categoryDocs := make(map[string][]string)

	for catName, links := range categoryLinks {
		categorySublinks[catName] = FetchCategoryDocumentLinks(catName, links)
	}

	uiprogress.Start()
	for catName, links := range categorySublinks {
		categoryDocs[catName] = FetchCategoryDocuments(catName, links)
	}

	// Save documents to /documents
	for catName, documents := range categoryDocs {
		os.MkdirAll("documents/"+catName, 0755)
		for i, document := range documents {
			fileName := "documents/" + catName + "/" + strconv.Itoa(i)
			ioutil.WriteFile(fileName, []byte(document), 0755)
		}
	}
}
