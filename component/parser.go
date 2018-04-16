// Package component implement Source parsing in a Category/Component tree.
package component

import (
	"gopkg.in/tent.v1/source"
)

const (
	catName = ".category.yaml"
)

// Parser is a Component parser
type Parser interface {
	Match(item source.Item) bool
	Parse(root *Category, item source.Item) error
}

var baseParsers = []Parser{categoryParser{}, segmentParser{}}

// Parse trasforms a Source in a Category tree.
func Parse(src source.Source, parsers ...Parser) (*Category, error) {
	var root = Category{ID: "root"}
	for i, err := src.Next(); i != nil; i, err = src.Next() {
		if err != nil {
			return nil, err
		}
		for _, p := range append(baseParsers, parsers...) {
			if !p.Match(i) {
				continue
			}
			if err := p.Parse(&root, i); err != nil {
				return nil, err
			}
			break
		}
	}
	root.sort()
	return &root, nil
}
