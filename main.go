package main

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please enter a URL.")
		return
	}

	urlStr := os.Args[1]
	match, _ := regexp.MatchString("^[http|//]", urlStr)
	if !match {
		urlStr = "https://" + urlStr
	}
	location, err := url.Parse(urlStr)
	if err != nil {
		fmt.Printf("You must enter a URL. You entered %s.", urlStr)
	}

	p := &Page{
		URL: location,
	}

	fetchErr := p.FetchPage()
	if fetchErr != nil {
		fmt.Println("We could not fetch your main page.")
		return
	}

	c := make(chan *Page, len(p.ChildPages))
	errChan := make(chan error)
	for i := range p.ChildPages {
		go func(child *Page) {
			err := child.FetchPage()
			if err != nil {
				errChan <- err
			}
			c <- child
		}(p.ChildPages[i])
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

	countWords(append(p.ChildPages, p))
}

func countWords(pages []*Page) map[string]int {
	counts := make(map[string]int)

	return counts
}
