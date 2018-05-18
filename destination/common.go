package destination

import (
	"context"

	"gopkg.in/tent.v1/item"
)

// Destination stores Items.
type Destination interface {
	Hash(ctx context.Context, i item.Item) (string, error)
	Create(ctx context.Context, i item.Item) error
	Update(ctx context.Context, i item.Item, hash string) error
	Delete(ctx context.Context, i item.Item, hash string) error
}
