package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Page - Represents a webpage
type Page struct {
	URL        *url.URL
	Text       string
	WordCounts map[string]int
	ChildPages []*Page
}

// FetchPage - fetches a page
func (p *Page) FetchPage() error {
	if len(p.URL.String()) < 1 {
		return errors.New("the page must have a URL")
	}

	resp, err := http.Get(p.URL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Could not read the page %s", p.URL.String())
		return err
	}

	p.Text = string(body)
	p.parseText()
	return nil
}

func (p *Page) parseText() {
	n, err := html.Parse(strings.NewReader(p.Text))
	if err != nil {
		fmt.Println("Could not parse the HTML.")
	}
	parsers := []*Parser{childFromLinkTags(p), wordCount(p)}
	parseNode(n, parsers)
}

func childFromLinkTags(p *Page) *Parser {
	condition := func(n *html.Node) bool {
		if !(n.Type == html.ElementNode && n.Data == "a") {
			return false
		}
		return true
	}

	parseFunc := func(n *html.Node) {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}
			url, _ := url.Parse(a.Val)

			// TODO: filter out mailto links
			if url.String() == p.URL.String() {
				continue
			}

			newURL := p.URL.ResolveReference(url)
			if newURL.Host == "" {
				fmt.Println("**************", "Base", p.URL.String(), "Link", url)
			}
			p.ChildPages = append(p.ChildPages, &Page{URL: newURL})
			break
		}
	}

	return &Parser{
		ShouldParse: condition,
		Parse:       parseFunc,
	}
}

func wordCount(p *Page) *Parser {
	condition := func(n *html.Node) bool {
		return false
	}

	parseFunc := func(n *html.Node) {
		// TODO
	}
	return &Parser{
		ShouldParse: condition,
		Parse:       parseFunc,
	}
}
