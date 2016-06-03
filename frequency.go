package main

import (
  "fmt"
  "io/ioutil"
  // "os"
  "strings"
  // "sort"
)

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

type Word struct {
  Freq int
  Word string
}

type Words []Word

func (words Words) Len() int {
    return len(words)
}

func (words Words) Less(i, j int) bool {
  return words[i].Freq < words[j].Freq
}

func (words Words) Swap(i, j int) {
  words[i], words[j] = words[j], words[i]
}

func main() {
  stop_list := []string{"the", "►", "category", "and", "this", "from", "for",
    "wikipedia", "not", "pages", "categories", "with", "can", "was", "were"}


  categories, _ := ioutil.ReadDir("data")

  for _, category := range categories {

    total_words := 0
    articles, _ := ioutil.ReadDir("data/" + category.Name())

    for _, article := range articles {
      b, err := ioutil.ReadFile("data/" + category.Name() + "/" + article.Name())
      if err != nil {
        panic(err)
      }

      content := string(b)

      words := strings.Fields(content)

      for i := 0; i < len(words); i++ {
        words[i] = strings.ToLower(words[i])
        words[i] = strings.Replace(words[i], ".", "", -1)
        words[i] = strings.Replace(words[i], "!", "", -1)
        words[i] = strings.Trim(words[i], " ")
      }

      for _, word := range words {

        if len(word) < 3 || len(word) > 15 {
          continue
        }

        if stringInSlice(word, stop_list) {
          continue
        }

        total_words++
      }
    }

    fmt.Println(category.Name() + ":", total_words)


  }



  return
  /*
  args := os.Args
  if len(args) < 2 {
    panic("Must supply category")
  }

  stop_list := []string{"the", "►", "category", "and", "this", "from", "for",
    "wikipedia", "not", "pages", "categories", "with", "can", "was", "were"}

  files, err := ioutil.ReadDir("data/" + args[1])
  if err != nil {
    panic(err)
  }

  freq := make(map[string]int)

  for _, f := range files {
    b, err := ioutil.ReadFile("data/" + args[1] + "/" + f.Name())
    if err != nil {
      panic(err)
    }

    content := string(b)

    words := strings.Fields(content)

    for i := 0; i < len(words); i++ {
      words[i] = strings.ToLower(words[i])
      words[i] = strings.Replace(words[i], ".", "", -1)
      words[i] = strings.Replace(words[i], "!", "", -1)
      words[i] = strings.Trim(words[i], " ")
    }

    for _, word := range words {

      if len(word) < 3 || len(word) > 15 {
        continue
      }

      if stringInSlice(word, stop_list) {
        continue
      }

      if _, ok := freq[word]; ok {
        freq[word]++
      } else {
        freq[word] = 1
      }
    }
  }

  list := make(Words, 0)
  for k, v := range freq {
    fmt.Println(v, k)
    new_word :=  Word{}
    new_word.Word = k
    new_word.Freq = v
    list = append(list, new_word)
  }

  sort.Sort(list)
  fmt.Println(list)
  */
}
