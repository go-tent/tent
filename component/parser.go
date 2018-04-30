// Package component implement Source parsing in a Category/Component tree.
package component

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"gopkg.in/tent.v1/source"
)

const (
	catName = ".category.yaml"
)

// Parser is a Component parser
type Parser interface {
	Match(name string) bool
	Parse(root *Category, item source.Item) error
}

var baseParsers = []Parser{catParser{}, segParser{}, picParser{}}

// Parse trasforms a Source in a Category tree.
func Parse(src source.Source, parsers ...Parser) (*Category, error) {
	var root = Category{ID: "root"}
	for i, err := src.Next(); i != nil; i, err = src.Next() {
		if err != nil {
			return nil, err
		}
		name := i.Name()
		for _, p := range append(baseParsers, parsers...) {
			if !p.Match(name) {
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

// ExtractMeta looks for "---\n" delimiters, returning what's between.
func ExtractMeta(r *bufio.Reader) (io.Reader, error) {
	row, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if !bytes.Equal([]byte("---\n"), row) {
		return nil, errors.New("Invalid header")
	}
	b := bytes.NewBuffer(nil)
	for {
		row, err := r.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		if bytes.Equal([]byte("---\n"), row) {
			break
		}
		b.Write(row)
	}
	return b, nil
}
