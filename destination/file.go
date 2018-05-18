package destination

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"

	"gopkg.in/tent.v1/item"
)

// NewFile returns a new File destination.
func NewFile(root string) *File {
	return &File{root: root}
}

// File takes a directory in the filesystem as destination.
type File struct {
	root string
}

// Hash returns the SHA1 for the item.
func (f *File) Hash(ctx context.Context, i item.Item) (string, error) {
	r, err := os.Open(f.path(i))
	if err != nil {
		return "", err
	}
	defer r.Close()
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Create adds a new Item to the Destination.
func (f *File) Create(ctx context.Context, i item.Item) error {
	if f.exist(i) {
		return os.ErrExist
	}
	return f.write(i, os.O_CREATE)
}

// Update overwrites an Item in the Destination.
func (f *File) Update(ctx context.Context, i item.Item, hash string) error {
	if !f.exist(i) {
		return os.ErrNotExist
	}
	if err := f.ensureHash(ctx, i, hash); err != nil {
		return err
	}
	return f.write(i, os.O_TRUNC)
}

// Delete removes an Item from the Destination.
func (f *File) Delete(ctx context.Context, i item.Item, hash string) error {
	if !f.exist(i) {
		return os.ErrNotExist
	}
	if err := f.ensureHash(ctx, i, hash); err != nil {
		return err
	}
	return os.Remove(f.path(i))
}

func (f *File) path(i item.Item) string {
	return path.Join(f.root, i.Name())
}

func (f *File) exist(i item.Item) bool {
	_, err := os.Stat(f.path(i))
	return err == nil
}

func (f *File) ensureHash(ctx context.Context, i item.Item, hash string) error {
	h, err := f.Hash(ctx, i)
	if err != nil {
		return err
	}
	if h != hash {
		return fmt.Errorf("conflict (expected %s, got %s)", h, hash)
	}
	return nil
}

func (f *File) write(i item.Item, mode int) error {
	r, err := i.Content()
	if err != nil {
		return err
	}
	defer r.Close()

	p := f.path(i)
	if err := os.MkdirAll(path.Dir(p), 0755); err != nil {
		return err
	}
	w, err := os.OpenFile(p, mode|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}
