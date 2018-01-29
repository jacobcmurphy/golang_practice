package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please enter a URL.")
		return
	}

	urlStr := os.Args[1]
	location, err := url.Parse(urlStr)
	if err != nil {
		fmt.Printf("You must enter a URL. You entered %s.", urlStr)
	}

	if location.Scheme == "" {
		location.Scheme = "https"
	}

	p, err := fetchPage(location)
	if err != nil {
		return
	}

	// fmt.Println(p.Text)

	for _, childPage := range p.ChildPages {
		go func(u *url.URL) {
			fmt.Println(u)
		}(childPage.URL)
	}
}

func fetchPage(url *url.URL) (*Page, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		fmt.Printf("Could not fetch the page %s", url)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Could not read the page %s", url)
		return nil, err
	}

	return &Page{
		URL:  url,
		Text: string(bodyBytes),
	}, nil
}

func extractHTML(string) string {
	return ""
}

func countWords(text string) map[string]int {
	var counts map[string]int

	return counts
}
