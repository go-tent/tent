package component

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

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

// Match implements the Parser interface.
func (segParser) Format() (string, []string) { return "s_", []string{".md"} }

// Parse populates the Segment with Item contents.
func (segParser) Parse(id string, r io.Reader) (Component, error) {
	b := bufio.NewReader(r)
	header, err := ExtractMeta(b)
	if err != nil {
		return nil, err
	}
	s := Segment{ID: id}
	if err := yaml.NewDecoder(header).Decode(&s); err != nil {
		return nil, err
	}
	s.Body, err = ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
