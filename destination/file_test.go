package destination

import (
	"context"
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
		ctx  = context.Background()
	)
	os.Remove(path.Join(root, i.ID))

	var (
		create = func(cxt context.Context, i item.Item, _ string) error { return dest.Create(ctx, i) }
		update = dest.Update
		delete = dest.Delete
	)

	var operations = []struct {
		f   func(context.Context, item.Item, string) error
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
		hash, _ := dest.Hash(ctx, i)
		if err := op.f(ctx, i, hash); err != op.err {
			t.Fatalf("%d) Expected %q, got %q", n, op.err, err)
		}
	}
}
