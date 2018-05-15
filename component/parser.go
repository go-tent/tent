// Package component implement Source parsing in a Category/Component tree.
package component

import (
	"fmt"
	"io"
	"path"
	"strings"

	"gopkg.in/tent.v1/source"
	"gopkg.in/yaml.v2"
)

// Decoder is a Component decoder.
type Decoder interface {
	// Format returns filename prefix and allowed extesions
	Format() (prefix string, ext []string)
	// Decode creates and returns the Component
	Decode(id string, r io.Reader) (Component, error)
}

// Decode trasforms a Source in a Category tree.
func Decode(src source.Source, extra ...Decoder) (*Category, error) {
	root := Category{ID: "root"}
	decoders := append([]Decoder{segDecoder{}, picDecoder{}}, extra...)
	if err := detectCollisions(decoders); err != nil {
		return nil, err
	}
	for i, err := src.Next(); i != nil; i, err = src.Next() {
		if err != nil {
			return nil, err
		}
		name := i.Name()
		_, file := path.Split(name)
		if file == ".category.yml" {
			if err := parseCategory(&root, i); err != nil {
				return nil, err
			}
			continue
		}
		if err := parseComponent(&root, i, decoders); err != nil {
			return nil, err
		}
	}
	root.sort()
	return &root, nil
}

func parseComponent(root *Category, i source.Item, decoders []Decoder) error {
	dir, file := path.Split(i.Name())
	for _, p := range decoders {
		name := matchDecoder(p, file)
		if name == "" {
			continue
		}
		r, err := i.Content()
		if err != nil {
			return err
		}
		defer r.Close()
		cmp, err := p.Decode(name, r)
		if err != nil {
			return fmt.Errorf("%s: %s", file, err)
		}
		cat := root.ensure(dir)
		cat.Components = append(cat.Components, cmp)
		break
	}
	return nil
}

func parseCategory(root *Category, i source.Item) error {
	dir, _ := path.Split(i.Name())
	contents, err := i.Content()
	if err != nil {
		return fmt.Errorf("%s: %s", dir, err)
	}
	defer contents.Close()
	cat := Category{ID: path.Base(dir)}
	if err := yaml.NewDecoder(contents).Decode(&cat); err != nil {
		return fmt.Errorf("%s: %s", dir, err)
	}
	parent := root.ensure(path.Dir(path.Clean(dir)))
	parent.Sub = append(parent.Sub, cat)
	return nil
}

func matchDecoder(p Decoder, name string) string {
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
