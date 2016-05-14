package main

import (
  "fmt"
  "log"
  "regexp"
  "strings"

  "github.com/PuerkitoBio/goquery"
)

// Channels
// =========
var getPagesChan = make(chan *Subcategory)
var getContentChan = make(chan Subcategory)

// Structs
// ========
type Page struct {
  Name string
  Url string
  Content string
}

func (page *Page) GetContent() {
  defer func() {
       // recover from panic if one occured. Set err to nil otherwise.
       err := recover()

       if err != nil {
         fmt.Println(err)
         return
       }
   }()

  // Go to URL
  doc, err := goquery.NewDocument(page.Url)
  if err != nil {
    log.Fatal(err)
  }

  //fmt.Println("fetched")

  doc.Find(".mw-body").Each(func(i int, s *goquery.Selection) {
    page.Content = s.Text()
  })
}

type Subcategory struct {
  Name string
  Url string
  Pages []Page
}

func (subcategory *Subcategory) GetPages() {
  fmt.Println(subcategory.Name + ": Started")

  // Go to URL
  doc, err := goquery.NewDocument(subcategory.Url)
  if err != nil {
    log.Fatal(err)
  }

  doc.Find(".mw-category-group+ .mw-category-group a").Each(func(i int, s *goquery.Selection) {
    page := Page{}
    page.Name = s.Text()
    page.Url, _ = s.Attr("href")
    page.Url = "https://en.wikipedia.org" + page.Url
    subcategory.Pages = append(subcategory.Pages, page)
  })

  getPagesChan <- subcategory
}

type Category struct {
  Name string
  Subcategories []Subcategory
}

// Functions
// ==========
func GetCategories() []Category {
  doc, err := goquery.NewDocument("https://en.wikipedia.org/wiki/Portal:Contents/Categories")
  if err != nil {
    log.Fatal(err)
  }

  r, err := regexp.Compile(`(?i)([^(]+)`)
  if err != nil {
    log.Fatal(err)
  }

  categories := make([]Category,0)
  doc.Find("big").Each(func(i int, s *goquery.Selection) {
    if i == 0 || i == 1 {
      return
    }

    // Get category name
    name := s.Text()
    name = r.FindAllStringSubmatch(name, -1)[0][0]
    name = strings.Trim(name, " ")
    category := Category{}
    category.Name = name
    categories = append(categories, category)
  })

  return categories
}

func GetSubcategories() []Subcategory {
  doc, err := goquery.NewDocument("https://en.wikipedia.org/wiki/Portal:Contents/Categories")
  if err != nil {
    log.Fatal(err)
  }

  subcategories := make([]Subcategory, 0)
  // Get Subcategories
  doc.Find("#mw-content-text div div .hlist li+ li a").Each(func(i int, s *goquery.Selection) {
    name := s.Text()
    url, _ := s.Attr("href")
    subcategory := Subcategory{}
    subcategory.Name = name
    subcategory.Url  = "https://en.wikipedia.org" + url
    subcategories = append(subcategories, subcategory)
  })

  return subcategories
}

func main() {
  // Get subcategories
  subcategories := GetSubcategories()

  // Fetch pages for each subcategory
  for i := 0; i < len(subcategories); i++ {
    go subcategories[i].GetPages()
  }

  // Wait until all pages are fetched
  subcategories_fetched := 0
  for {
    subcategory := <-getPagesChan
    fmt.Println(subcategory.Name + ": Finished with", len(subcategory.Pages), "pages")
    subcategories_fetched++
    if subcategories_fetched >= len(subcategories) {
      break
    }
  }


  queue := make([]func(), 0)

  // Loop through subcategories and fetch content for each page
  for i := 0; i < len(subcategories); i++ {
    sub := subcategories[i]

    if len(sub.Pages) < 1 {
      continue
    }

    queue = append(queue, func() {
      fmt.Println(sub.Name + ": Started Fetching Pages")

      // Loop through pages
      for j := 0; j < len(sub.Pages); j++ {
        //fmt.Println(j, len(sub.Pages))
        sub.Pages[j].GetContent()
      }
      getContentChan <- sub
    })
  }


  queue_len := len(queue)
  // Start first n threads
  threads := 10
  for t := 0; t < threads; t++ {
    f := queue[0]
    queue = queue[1:]
    go f()
  }

  // Wait until all pages are fetched
  subcategories_fetched = 0

  for {

    subcategory := <-getContentChan
    fmt.Println("============================== " + subcategory.Name + ": Finished Fetching Pages =============================")

    f := queue[0]
    if len(queue) > 1 {
      queue = queue[1:]
    }
    go f()

    fmt.Println(subcategories_fetched, len(subcategories))
    subcategories_fetched++
    if subcategories_fetched >= queue_len {
      break
    }
  }
}
