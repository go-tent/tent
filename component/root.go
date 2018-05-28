// Package component implement Source parsing in a Category/Component tree.
package component

import (
	"fmt"
	"io"
	"path"
	"strings"

	"gopkg.in/tent.v1/item"
	"gopkg.in/tent.v1/source"
)

// Component represents a leaf node.
type Component interface {
	// Returns the ID for the current component
	GetID() string
	// Format returns filename prefix and allowed extesions
	Format() (prefix string, ext []string)
	// Decode creates and returns a new Component
	Decode(id string, r io.Reader) (Component, error)
	// Order is used for sorting Componenets
	Order() float64
	// Encode returns Component's contents for Item
	Encode() ([]byte, error)
}

func NewItem(prefix []string, cmp Component) (item.Item, error) {
	dir := path.Join(prefix...)
	if cat, ok := cmp.(*Category); ok {
		b, err := cat.Encode()
		if err != nil {
			return nil, err
		}
		return item.Memory{ID: path.Join(dir, ".category.yml"), Contents: b}, nil
	}
	name := cmp.GetID()
	if pre, exts := cmp.Format(); len(exts) == 1 {
		name = pre + name + exts[0]
	}
	b, err := cmp.Encode()
	if err != nil {
		return nil, err
	}
	return item.Memory{ID: path.Join(dir, name), Contents: b}, nil
}

// NewRoot returns a new Root.
func NewRoot(components ...Component) (*Root, error) {
	if err := detectCollisions(components); err != nil {
		return nil, err
	}
	return &Root{Category: new(Category), decoders: components}, nil
}

// Root is a container for a component Tree.
type Root struct {
	*Category
	decoders []Component
}

// IsValid verifies the existence of an Item Component.
func (r *Root) IsValid(i item.Item) error {
	_, file := path.Split(i.Name())
	if file == ".category.yml" {
		_, err := r.decodeCategory(i)
		return err
	}
	cmp, err := r.decodeComponent(i)
	if err != nil {
		return err
	}
	if cmp == nil {
		return fmt.Errorf("No parser for %s", path.Base(i.Name()))
	}
	return nil
}

// Decode trasforms a Source in a Category tree.
func (r *Root) Decode(src source.Source) error {
	root := Category{ID: "root"}
	for i, err := src.Next(); i != nil; i, err = src.Next() {
		if err != nil {
			return err
		}
		name := i.Name()
		dir, file := path.Split(name)
		if file == ".category.yml" {
			cat, err := r.decodeCategory(i)
			if err != nil {
				return err
			}
			parent := root.ensure(path.Dir(path.Clean(dir)))
			parent.Sub = append(parent.Sub, *cat)
			continue
		}
		cmp, err := r.decodeComponent(i)
		if err != nil {
			return err
		}
		parent := root.ensure(dir)
		parent.Components = append(parent.Components, cmp)
	}
	root.sort()
	r.Category = &root
	return nil
}

func (r *Root) decodeCategory(i item.Item) (*Category, error) {
	dir, _ := path.Split(i.Name())
	contents, err := i.Content()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", dir, err)
	}
	defer contents.Close()

	cat, err := (*Category).decode(nil, path.Base(dir), contents)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", dir, err)
	}
	return cat, nil
}

func (r *Root) decodeComponent(i item.Item) (Component, error) {
	_, file := path.Split(i.Name())
	for _, p := range r.decoders {
		name := r.matchDecoder(p, file)
		if name == "" {
			continue
		}
		r, err := i.Content()
		if err != nil {
			return nil, err
		}
		defer r.Close()
		cmp, err := p.Decode(name, r)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", file, err)
		}
		return cmp, nil
	}
	return nil, nil
}

func (r *Root) matchDecoder(p Component, name string) string {
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
