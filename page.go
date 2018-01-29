package main

import (
	"net/url"
)

// Page - Represents a webpage
type Page struct {
	URL        *url.URL
	Text       string
	WordCounts map[string]int
	ChildPages []Page
}

func (p *Page) getChildren() {

}

func (p *Page) getWordCounts() {

}
