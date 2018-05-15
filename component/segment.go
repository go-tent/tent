package component

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
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

// segDecoder is the Decoder for Segment.
type segDecoder struct{}

// Match implements the Decoder interface.
func (segDecoder) Format() (string, []string) { return "s_", []string{".md"} }

// Decode populates the Segment with Item contents.
func (segDecoder) Decode(id string, r io.Reader) (Component, error) {
	b := bufio.NewReader(r)
	header, err := extractMeta(b)
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

// extractMeta looks for "---\n" delimiters, returning what's between.
func extractMeta(r *bufio.Reader) (io.Reader, error) {
	row, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	if !bytes.Equal([]byte("---"), bytes.TrimSuffix(row, []byte("\r"))) {
		return nil, errors.New("Invalid header")
	}
	b := bytes.NewBuffer(nil)
	for {
		row, _, err := r.ReadLine()
		if err != nil {
			return nil, err
		}
		if bytes.Equal([]byte("---"), bytes.TrimSuffix(row, []byte("\r"))) {
			break
		}
		fmt.Fprintln(b, string(row))
	}
	return b, nil
}
