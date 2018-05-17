package destination

import (
	"io"
	"os"
	"path"

	"gopkg.in/tent.v1/hash"
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

// Create adds a new Item to the Destination.
func (f *File) Create(i item.Item) error {
	if f.exist(i) {
		return os.ErrExist
	}
	return f.write(i, os.O_CREATE)
}

// Update overwrites an Item in the Destination.
func (f *File) Update(i item.Item, hash []byte) error {
	if !f.exist(i) {
		return os.ErrNotExist
	}
	if err := f.ensureHash(i, hash); err != nil {
		return err
	}
	return f.write(i, os.O_TRUNC)
}

// Delete removes an Item from the Destination.
func (f *File) Delete(i item.Item, hash []byte) error {
	if !f.exist(i) {
		return os.ErrNotExist
	}
	if err := f.ensureHash(i, hash); err != nil {
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

func (f *File) ensureHash(i item.Item, h []byte) error {
	r, err := os.Open(f.path(i))
	if err != nil {
		return err
	}
	defer r.Close()

	if err := hash.Verify(r, h); err != nil {
		return err
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
