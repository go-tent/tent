package source

import (
	"strings"

	"github.com/go-tent/tent/item"
)

// Source is an origin of Items.
type Source interface {
	Next() (item.Item, error)
}

// PathFilter is used to exclude/include files in a FileSource.
type PathFilter func(string) bool

// FilterSuffix filters files by the given suffix.
func FilterSuffix(filter string) PathFilter {
	return func(s string) bool { return strings.HasSuffix(s, filter) }
}

// FilterPrefix filters files by the given prefix.
func FilterPrefix(filter string) PathFilter {
	return func(s string) bool { return strings.HasPrefix(s, filter) }
}

// Memory is a static Source.
type Memory struct {
	Items []item.Memory
	i     int
}

// Next implements the Source interface.
func (m *Memory) Next() (item.Item, error) {
	if m.i == len(m.Items) {
		return nil, nil
	}
	m.i++
	return m.Items[m.i-1], nil
}
