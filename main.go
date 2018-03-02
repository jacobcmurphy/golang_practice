package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Query().Get("url")
	wordCount := getCounts(urlStr)

	fmt.Fprintf(w, urlStr)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCounts(urlStr string) map[string]int {
	match, _ := regexp.MatchString("^[http|//]", urlStr)
	if !match {
		urlStr = "https://" + urlStr
	}
	location, err := url.Parse(urlStr)
	if err != nil {
		fmt.Printf("You must enter a URL. You entered %s.", urlStr)
	}

	p := &Page{
		URL:        location,
		WordCounts: make(map[string]int),
	}

	fetchErr := p.FetchPage()
	if fetchErr != nil {
		fmt.Println("We could not fetch your main page.")
		return
	}

	fmt.Println("PAGE COUNT: ", len(p.ChildPages))

	c := make(chan *Page, len(p.ChildPages))
	errChan := make(chan error)
	for i := range p.ChildPages {
		go func(idx int) {
			child := p.ChildPages[idx]
			err := child.FetchPage()
			if err != nil {
				errChan <- err
			}
			c <- child
			fmt.Println(fmt.Sprintf("FINISHED %d", idx))
		}(i)
	}

	for i := 0; i < len(p.ChildPages); i++ {
		select {
		case e := <-errChan:
			fmt.Println(e)
		case <-c:
			// noop
		}
	}
	close(c)
	close(errChan)

	return countWords(append(p.ChildPages, p))
}

func countWords(pages []*Page) map[string]int {
	counts := make(map[string]int)

	for _, page := range pages {
		for word, count := range page.WordCounts {
			counts[word] = counts[word] + count
		}
	}

	return counts
}
