package item

import (
	"bytes"
	"io"
	"io/ioutil"
)

// Item represents a stored component.
type Item interface {
	Name() string
	Content() (io.ReadCloser, error)
}

// Memory is a static Item.
type Memory struct {
	ID       string
	Contents []byte
}

// Name implements the Item interface.
func (m Memory) Name() string { return m.ID }

// Content implements the Item interface.
func (m Memory) Content() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(m.Contents)), nil
}
