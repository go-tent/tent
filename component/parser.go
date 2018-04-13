package component

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"gopkg.in/tent.v1/source"
	yaml "gopkg.in/yaml.v2"
)

const (
	contentPrefix = "contents_"
	catName       = ".category.yaml"
)

// Parse trasforms a Source in a Category tree.
func Parse(src source.Source) (*Category, error) {
	var root = Category{ID: "root"}
	for i, err := src.Next(); i != nil; i, err = src.Next() {
		if err != nil {
			return nil, err
		}
		path := strings.Split(i.Name(), "/")
		if len(path) < 2 || !strings.HasPrefix(path[0], contentPrefix) {
			return nil, nil
		}
		loc := strings.TrimPrefix(path[0], contentPrefix)
		var cat = root.ensure(loc)
		for _, id := range path[1 : len(path)-1] {
			cat = cat.ensure(id)
		}
		name := path[len(path)-1]
		switch {
		case cat.Match(name):
			if err := cat.Parse(i); err != nil {
				return nil, err
			}
		case (&Segment{}).Match(name):
			var s Segment
			if err := s.Parse(i); err != nil {
				return nil, err
			}
			cat.Segments = append(cat.Segments, s)
		}
	}
	root.sort()
	return &root, nil
}

// parseMeta extracts the meta and the index from a Reader.
// If sep is true it looks for "---\n" delimiters.
func parseMeta(r *bufio.Reader, m *map[string]string, idx *float64, sep bool) error {
	if sep {
		row, err := r.ReadBytes('\n')
		if err != nil {
			return err
		}
		if !bytes.Equal([]byte("---\n"), row) {
			return errInvalidHeader
		}
	}
	var header io.Reader = r
	if sep {
		b := bytes.NewBuffer(nil)
		for {
			row, err := r.ReadBytes('\n')
			if err != nil {
				return err
			}
			if bytes.Equal([]byte("---\n"), row) {
				break
			}
			b.Write(row)
		}
		header = b
	}
	if err := yaml.NewDecoder(header).Decode(m); err != nil {
		return err
	}
	index, err := strconv.ParseFloat((*m)["index"], 64)
	if err != nil {
		return fmt.Errorf("Invalid Index: %s", (*m)["index"])
	}
	delete((*m), "index")
	*idx = index
	return nil
}
