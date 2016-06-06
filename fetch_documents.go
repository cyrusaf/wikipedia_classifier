package main

import (
  "fmt"
  "io/ioutil"
  "sync"

  "github.com/go-yaml/yaml"
  "github.com/PuerkitoBio/goquery"
)

func GetCategoryLinks() map[string][]string {
  dat, err := ioutil.ReadFile("category_links.yml")
  if err != nil {
    panic(err)
  }

  var category_links map[string][]string
  err = yaml.Unmarshal(dat, &category_links)
  if err != nil {
    panic(err)
  }

  return category_links
}

func FetchCategoryDocumentLinks(cat_name string, links []string) []string{
  fmt.Println("\nFetching", len(links), "document links from " + cat_name + "...")

  sublinks := make([]string, 0)
  for _, link := range links {

    // Fetch webpage
    doc, err := goquery.NewDocument(link)
    if err != nil {
      panic(err)
    }

    // Find links to documents
    doc.Find("#mw-pages .mw-category-group a").Each(func(i int, s *goquery.Selection) {
      url, _ := s.Attr("href")
      sublinks = append(sublinks, "https://en.wikipedia.org" + url)
    })
  }

  return sublinks
}

func FetchCategoryDocuments(cat_name string, links[]string) []string {
  fmt.Println("\nFetching", len(links), "documents from " + cat_name + "...")

  documents := make([]string, 0)
  var wg sync.WaitGroup

  for _, link := range links {
    wg.Add(1)
    go func() {
      defer wg.Done()

      // Fetch webpage
      doc, err := goquery.NewDocument(link)
      if err != nil {
        panic(err)
      }

      // Find document content
      content := ""
      doc.Find("#bodyContent").Each(func(i int, s *goquery.Selection) {
        content += s.Text()
      })

      documents = append(documents, content)
      fmt.Printf(".")
    }()
  }

  wg.Wait()
  fmt.Println("\nFetched", len(documents), "documents from " + cat_name)
  return documents
}

func main() {
  category_links    := GetCategoryLinks()
  category_sublinks := make(map[string][]string)
  category_docs     := make(map[string][]string)

  for cat_name, links := range category_links {
    category_sublinks[cat_name] = FetchCategoryDocumentLinks(cat_name, links)
  }

  for cat_name, links := range category_sublinks {
    category_docs[cat_name] = FetchCategoryDocuments(cat_name, links)
  }

}
