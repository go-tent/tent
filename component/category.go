package component

import (
	"fmt"
	"sort"
	"strings"
)

// Category is a branch node in the tree.
type Category struct {
	ID         string            `yaml:"id"`
	Index      float64           `yaml:"index"`
	Meta       map[string]string `yaml:",inline"`
	Sub        []Category        `yaml:"sub"`
	Components []Component       `yaml:"components"`
}

// Order implements the Component interface.
func (c *Category) Order() float64 { return c.Index }

func (c *Category) String() string {
	return fmt.Sprintf("Category:%s Idx:%v Meta:%v", c.ID, c.Index, c.Meta)
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

// Component represents a leaf node.
type Component interface {
	Order() float64
}
