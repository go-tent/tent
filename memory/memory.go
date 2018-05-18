package memory

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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

// Hash returns the SHA1 for the item.
func (d *Destination) Hash(ctx context.Context, i item.Item) (string, error) {
	h := sha1.New()
	if _, err := h.Write(d.items[i.Name()]); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Create adds a new Item.
func (d *Destination) Create(ctx context.Context, i item.Item) error {
	_, ok := d.items[i.Name()]
	if ok {
		return os.ErrExist
	}
	i.Content()
	return nil
}

// Update writes an existing Item.
func (d *Destination) Update(ctx context.Context, i item.Item, hash string) error {
	name := i.Name()
	b, ok := d.items[name]
	if !ok {
		return os.ErrNotExist
	}
	if err := d.ensureHash(ctx, i, hash); err != nil {
		return err
	}
	d.items[name] = b
	return nil
}

// Delete removes an existing Item.
func (d *Destination) Delete(ctx context.Context, i item.Item, hash string) error {
	name := i.Name()
	if _, ok := d.items[name]; !ok {
		return os.ErrNotExist
	}
	if err := d.ensureHash(ctx, i, hash); err != nil {
		return err
	}
	delete(d.items, name)
	return nil
}

func (d *Destination) ensureHash(ctx context.Context, i item.Item, hash string) error {
	h, err := d.Hash(ctx, i)
	if err != nil {
		return err
	}
	if h != hash {
		return fmt.Errorf("conflict (expected %s, got %s)", h, hash)
	}
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
