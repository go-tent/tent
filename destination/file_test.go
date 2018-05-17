package destination

import (
	"os"
	"path"
	"testing"

	"gopkg.in/tent.v1/item"
	"gopkg.in/tent.v1/memory"
)

func TestFile(t *testing.T) {
	var (
		root = "."
		dest = NewFile(root)
		i    = memory.Item{ID: "id", Contents: "contents"}
		hash = []byte{
			0xD1, 0xB2, 0xA5, 0x9F, 0xBE, 0xA7, 0xE2, 0x00,
			0x77, 0xAF, 0x9F, 0x91, 0xB2, 0x7E, 0x95, 0xE8,
			0x65, 0x06, 0x1B, 0x27, 0x0B, 0xE0, 0x3F, 0xF5,
			0x39, 0xAB, 0x3B, 0x73, 0x58, 0x78, 0x82, 0xE8,
		}
	)
	os.Remove(path.Join(root, i.ID))

	var (
		create = func(i item.Item, _ []byte) error { return dest.Create(i) }
		update = dest.Update
		delete = dest.Delete
	)

	var operations = []struct {
		f   func(item.Item, []byte) error
		err error
	}{
		{delete, os.ErrNotExist},
		{update, os.ErrNotExist},
		{create, nil},
		{create, os.ErrExist},
		{update, nil},
		{delete, nil},
	}

	for n, op := range operations {
		if err := op.f(i, hash); err != op.err {
			t.Fatalf("%d) Expected %q, got %q", n, op.err, err)
		}
	}
}
