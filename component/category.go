package component

import (
	"bufio"
	"fmt"
	"sort"

	"gopkg.in/tent.v1/source"
)

// Category is a parent node in the tree.
type Category struct {
	ID    string
	Index float64
	Meta  map[string]string
	Sub   []Category
}

func (c Category) String() string {
	return fmt.Sprintf("ID:%s Idx:%v Meta:%v", c.ID, c.Index, c.Meta)
}

func (c *Category) ensure(id string) *Category {
	for i := range c.Sub {
		if c.Sub[i].ID == id {
			return &c.Sub[i]
		}
	}
	c.Sub = append(c.Sub, Category{ID: id})
	return &c.Sub[len(c.Sub)-1]
}

func (c Category) sort() {
	sort.Slice(c.Sub, func(i, j int) bool { return c.Sub[i].Index < c.Sub[j].Index })
	for i := range c.Sub {
		c.Sub[i].sort()
	}
}

// Match tells if it's a Category from the name.
func (c *Category) Match(name string) bool {
	return name == catName
}

// Parse populates the Category with Item contents.
func (c *Category) Parse(item source.Item) error {
	contents, err := item.Content()
	if err != nil {
		return err
	}
	defer contents.Close()
	return parseMeta(bufio.NewReader(contents), &c.Meta, &c.Index, false)
}
