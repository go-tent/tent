// Package component implement Source parsing in a Category/Component tree.
package component

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"path"
	"strings"

	"gopkg.in/tent.v1/source"
	yaml "gopkg.in/yaml.v2"
)

// Parser is a Component parser.
type Parser interface {
	// Format returns filename prefix and allowed extesions
	Format() (prefix string, ext []string)
	// Parse creates and returns the Component
	Parse(id string, r io.Reader) (Component, error)
}

// Parse trasforms a Source in a Category tree.
func Parse(src source.Source, extra ...Parser) (*Category, error) {
	root := Category{ID: "root"}
	parsers := append([]Parser{segParser{}, picParser{}}, extra...)
	if err := detectCollisions(parsers); err != nil {
		return nil, err
	}
	for i, err := src.Next(); i != nil; i, err = src.Next() {
		if err != nil {
			return nil, err
		}
		name := i.Name()
		_, file := path.Split(name)
		if file == ".category.yml" {
			parseCategory(&root, i)
			continue
		}
		if err := parseComponent(&root, i, parsers); err != nil {
			return nil, err
		}
	}
	root.sort()
	return &root, nil
}

func parseComponent(root *Category, i source.Item, parsers []Parser) error {
	dir, file := path.Split(i.Name())
	ext := path.Ext(file)
parsers:
	for _, p := range parsers {
		prefix, validExts := p.Format()
		if prefix != "" && !strings.HasPrefix(file, prefix) {
			continue
		}
		for i, e := range validExts {
			if ext == e {
				break
			}
			if i == len(validExts)-1 {
				continue parsers
			}
		}
		file = strings.TrimPrefix(file, prefix)
		if len(validExts) == 1 {
			file = strings.TrimSuffix(file, ext)
		}
		r, err := i.Content()
		if err != nil {
			return err
		}
		defer r.Close()
		cmp, err := p.Parse(file, r)
		if err != nil {
			return err
		}
		cat := root.ensure(dir)
		cat.Components = append(cat.Components, cmp)
		break
	}
	return nil
}

func parseCategory(root *Category, i source.Item) error {
	contents, err := i.Content()
	if err != nil {
		return err
	}
	defer contents.Close()
	dir, _ := path.Split(i.Name())
	cat := Category{ID: path.Base(dir)}
	if err := yaml.NewDecoder(contents).Decode(&cat); err != nil {
		return err
	}
	parent := root
	if dir := path.Dir(path.Clean(dir)); dir != "." {
		parent = parent.ensure(dir)
	}
	parent.Sub = append(parent.Sub, cat)
	return nil
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
