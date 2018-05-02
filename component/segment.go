package component

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"gopkg.in/tent.v1/source"
	yaml "gopkg.in/yaml.v2"
)

const extSegment = ".md"

// Segment is an Article.
type Segment struct {
	ID    string            `yaml:"id"`
	Index float64           `yaml:"index"`
	Meta  map[string]string `yaml:",inline"`
	Body  []byte
}

// Order implements the Component interface.
func (s *Segment) Order() float64 { return s.Index }

func (s Segment) String() string {
	return fmt.Sprintf("Segment:%s Idx:%v Meta:%v", s.ID, s.Index, s.Meta)
}

// segParser is the Parser for Segment.
type segParser struct{}

// Match tells if it's a Segment from the name.
func (segParser) Match(name string) bool {
	return strings.ToLower(path.Ext(name)) == extSegment
}

// Parse populates the Segment with Item contents.
func (segParser) Parse(c *Category, item source.Item) error {
	contents, err := item.Content()
	if err != nil {
		return err
	}
	defer contents.Close()

	b := bufio.NewReader(contents)
	header, err := ExtractMeta(b)
	if err != nil {
		return err
	}
	s := Segment{ID: strings.TrimSuffix(item.Name(), extSegment)}
	if err := yaml.NewDecoder(header).Decode(&s); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(b)
	if err != nil {
		return err
	}
	s.Body = body

	c.Components = append(c.Components, &s)
	return nil
}
