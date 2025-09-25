// Package main demonstrates Slate.js operations with JSON Patch.
package main

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/pkg/slate"
)

func main() {
	// Text node
	textNode := slate.NewBasicTextNode("Bold text", slate.WithBold(true))
	fmt.Printf("Text: %s\n", textNode.GetText())

	// Element node
	elementNode := slate.NewBasicElementNode([]slate.TextNode{*textNode}, slate.BasicProps{})
	fmt.Printf("Children: %d\n", len(elementNode.GetChildren()))

	// Combined properties
	combinedProps := slate.CombineProps(
		slate.WithBold(true),
		slate.WithItalic(true),
	)
	styledNode := slate.NewBasicTextNode("Styled", combinedProps)
	fmt.Printf("Styled: %+v\n", styledNode.GetProperties())
}
