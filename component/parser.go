// Package component implement Source parsing in a Category/Component tree.
package component

import (
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
	for _, p := range parsers {
		name := matchParser(p, file)
		if name == "" {
			continue
		}
		r, err := i.Content()
		if err != nil {
			return err
		}
		defer r.Close()
		cmp, err := p.Parse(name, r)
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
	parent := root.ensure(path.Dir(path.Clean(dir)))
	parent.Sub = append(parent.Sub, cat)
	return nil
}

func matchParser(p Parser, name string) string {
	ext := path.Ext(name)
	prefix, validExts := p.Format()
	if prefix != "" && !strings.HasPrefix(name, prefix) {
		return ""
	}
	for _, e := range validExts {
		if ext != e {
			continue
		}
		name = strings.TrimPrefix(name, prefix)
		if len(validExts) == 1 {
			name = strings.TrimSuffix(name, ext)
		}
		return name
	}
	return ""
}
