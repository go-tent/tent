package source

import (
	"strings"

	"gopkg.in/tent.v1/item"
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
