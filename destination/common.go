package destination

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/go-tent/tent/item"
)

// Destination stores Items.
type Destination interface {
	Hash(ctx context.Context, i item.Item) (string, error)
	Create(ctx context.Context, i item.Item) error
	Update(ctx context.Context, i item.Item, hash string) error
	Delete(ctx context.Context, i item.Item, hash string) error
}

// Memory is a static Destination.
type Memory struct {
	items map[string][]byte
}

// Hash returns the SHA1 for the item.
func (m *Memory) Hash(ctx context.Context, i item.Item) (string, error) {
	h := sha1.New()
	if _, err := h.Write(m.items[i.Name()]); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Create adds a new Item.
func (m *Memory) Create(ctx context.Context, i item.Item) error {
	_, ok := m.items[i.Name()]
	if ok {
		return os.ErrExist
	}
	i.Content()
	return nil
}

// Update writes an existing Item.
func (m *Memory) Update(ctx context.Context, i item.Item, hash string) error {
	name := i.Name()
	b, ok := m.items[name]
	if !ok {
		return os.ErrNotExist
	}
	if err := m.ensureHash(ctx, i, hash); err != nil {
		return err
	}
	m.items[name] = b
	return nil
}

// Delete removes an existing Item.
func (m *Memory) Delete(ctx context.Context, i item.Item, hash string) error {
	name := i.Name()
	if _, ok := m.items[name]; !ok {
		return os.ErrNotExist
	}
	if err := m.ensureHash(ctx, i, hash); err != nil {
		return err
	}
	delete(m.items, name)
	return nil
}

func (m *Memory) ensureHash(ctx context.Context, i item.Item, hash string) error {
	h, err := m.Hash(ctx, i)
	if err != nil {
		return err
	}
	if h != hash {
		return fmt.Errorf("conflict (expected %s, got %s)", h, hash)
	}
	return nil
}
