// Package component implement Source parsing in a Category/Component tree.
package component

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"path"

	"gopkg.in/tent.v1/source"
)

// Parser is a Component parser.
type Parser interface {
	Match(name string) bool
	Parse(parent *Category, item source.Item) error
}

// Parse trasforms a Source in a Category tree.
func Parse(src source.Source, extra ...Parser) (*Category, error) {
	var (
		root    = Category{ID: "root"}
		parsers = make([]Parser, 0, 3+len(extra))
	)
	parsers = append(parsers, catParser{}, leaf{segParser{}}, leaf{picParser{}})
	for _, p := range extra {
		parsers = append(parsers, leaf{p})
	}
	for i, err := src.Next(); i != nil; i, err = src.Next() {
		if err != nil {
			return nil, err
		}
		name := i.Name()
		for _, p := range parsers {
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

// leaf is a wrapper for leaf node Parsers.
type leaf struct {
	Parser
}

// Parse calls the underlying Parser using the proper branch node and leafNode.
func (l leaf) Parse(root *Category, item source.Item) error {
	dir, _ := path.Split(item.Name())
	return l.Parser.Parse(root.ensure(dir), leafItem{item})
}

type leafItem struct {
	source.Item
}

func (l leafItem) Name() string { return path.Base(l.Item.Name()) }
