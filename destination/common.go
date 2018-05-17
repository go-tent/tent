package destination

import "gopkg.in/tent.v1/item"

// Destination stores Items.
type Destination interface {
	Create(item item.Item) error
	Update(item item.Item, hash []byte) error
	Delete(item item.Item, hash []byte) error
}
