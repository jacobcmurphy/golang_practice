package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
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

	fmt.Println("URL: ", p.URL.String())
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	resp, err := client.Get(p.URL.String())

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
		return n.Type == html.ElementNode && n.Data == "a"
	}

	parseFunc := func(n *html.Node) {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}
			linkToParent := a.Val == p.URL.String()
			contactLink, _ := regexp.MatchString("(mailto:|tel:)", a.Val)
			if linkToParent || contactLink {
				continue
			}

			url, _ := url.Parse(a.Val)
			url = p.URL.ResolveReference(url)
			p.ChildPages = append(p.ChildPages, &Page{URL: url, WordCounts: make(map[string]int)})
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
		isTextNode := n.Type == html.TextNode && strings.TrimSpace(n.Data) != ""
		ancestorsAreValid := true

		for ancestor := n.Parent; ancestor != nil; ancestor = ancestor.Parent {
			if ancestor.Type == html.ElementNode && (ancestor.Data == "script" || ancestor.Data == "style") {
				ancestorsAreValid = false
				break
			}
		}
		return isTextNode && ancestorsAreValid
	}

	parseFunc := func(n *html.Node) {
		words := strings.Fields(n.Data)
		for _, word := range words {
			p.WordCounts[word]++
		}
		// fmt.Println(n.Data)
	}
	return &Parser{
		ShouldParse: condition,
		Parse:       parseFunc,
	}
}
