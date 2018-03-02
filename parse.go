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
	nodes := []*html.Node{n}

	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		for i := range parsers {
			if parsers[i].ShouldParse(node) {
				parsers[i].Parse(node)
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			nodes = append(nodes, c)
		}
	}
}
