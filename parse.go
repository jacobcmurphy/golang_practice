package main

import (
	"golang.org/x/net/html"
)

// ParseCondition - determines if we want to parse a Node
type ParseCondition func(*html.Node) bool

// ParseFunction - logic for parsing a Node
type ParseFunction func(*html.Node)

// Parser - determines if a Node should be parsed, then parses it
type Parser struct {
	ShouldParse ParseCondition
	Parse       ParseFunction
}

func parseNode(n *html.Node, parsers []*Parser) {
	// fmt.Println("***", n.Data, "****")
	for _, parser := range parsers {
		if parser.ShouldParse(n) {
			parser.Parse(n)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNode(c, parsers)
	}
}
