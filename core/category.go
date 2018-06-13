package core

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Category is a branch node in the tree.
type Category struct {
	ID         string            `yaml:"-"`
	Index      float64           `yaml:"index,omitempty"`
	Meta       map[string]string `yaml:",inline"`
	Sub        []Category        `yaml:"-"`
	Components []Component       `yaml:"-"`
}

// GetID implements the Component interface.
func (c *Category) GetID() string { return c.ID }

// Format implements the Component interface.
func (c *Category) Format() (string, []string) { return "", nil }

// Order implements the Component interface.
func (c *Category) Order() float64 { return c.Index }

func (c *Category) String() string {
	return fmt.Sprintf("Category:%s Idx:%v Meta:%v", c.ID, c.Index, c.Meta)
}

// Encode returns Item contents.
func (c *Category) Encode() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(b).Encode(c); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (c *Category) sort() {
	sort.Slice(c.Sub, func(i, j int) bool {
		return c.Sub[i].Index < c.Sub[j].Index
	})
	sort.Slice(c.Components, func(i, j int) bool {
		return c.Components[i].Order() < c.Components[j].Order()
	})
	for i := range c.Sub {
		c.Sub[i].sort()
	}
}

// ensure follows the path to a leaf node, creating all needed ones.
func (c *Category) ensure(path string) *Category {
	if path == "" || path == "." {
		return c
	}
item:
	for _, id := range strings.FieldsFunc(path, func(r rune) bool { return r == '/' }) {
		for i := range c.Sub {
			if c.Sub[i].ID == id {
				c = &c.Sub[i]
				continue item
			}
		}
		c.Sub = append(c.Sub, Category{ID: id})
		c = &c.Sub[len(c.Sub)-1]
	}
	return c
}

// Decode returns a new Category with Item contents.
func (c *Category) Decode(id string, r io.Reader) (Component, error) {
	return c.decode(id, r)
}

func (*Category) decode(id string, r io.Reader) (*Category, error) {
	c := Category{ID: id}
	if err := yaml.NewDecoder(r).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
