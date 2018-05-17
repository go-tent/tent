package memory

import (
	"bytes"
	"io"
	"io/ioutil"

	"gopkg.in/tent.v1/item"
)

// Item is a static Item.
type Item struct {
	ID       string
	Contents string
}

// Name implements the Item interface.
func (i Item) Name() string { return i.ID }

// Content implements the Item interface.
func (i Item) Content() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBufferString(i.Contents)), nil
}

// Source is a static Source.
type Source struct {
	Items []Item
	i     int
}

// Next implements the Source interface.
func (m *Source) Next() (item.Item, error) {
	if m.i == len(m.Items) {
		return nil, nil
	}
	m.i++
	return m.Items[m.i-1], nil
}
