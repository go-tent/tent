package item

import (
	"io"
)

// Item represents a stored component.
type Item interface {
	Name() string
	Content() (io.ReadCloser, error)
}
