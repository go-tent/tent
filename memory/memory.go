package memory

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/tent.v1/hash"
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

// Destination is a static Destination.
type Destination struct {
	items map[string][]byte
}

// Create adds a new Item.
func (d *Destination) Create(item item.Item) error {
	_, ok := d.items[item.Name()]
	if ok {
		return os.ErrExist
	}
	item.Content()
	return nil
}

// Update writes an existing Item.
func (d *Destination) Update(item item.Item, h []byte) error {
	name := item.Name()
	b, ok := d.items[name]
	if !ok {
		return os.ErrNotExist
	}
	if err := hash.Verify(bytes.NewReader(b), h); err != nil {
		return err
	}
	d.items[name] = b
	return nil
}

// Delete removes an existing Item.
func (d *Destination) Delete(item item.Item, h []byte) error {
	name := item.Name()
	b, ok := d.items[name]
	if !ok {
		return os.ErrNotExist
	}
	if err := hash.Verify(bytes.NewReader(b), h); err != nil {
		return err
	}
	delete(d.items, name)
	return nil
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
