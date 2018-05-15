package source

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
)

// Source is an origin of Items.
type Source interface {
	Next() (Item, error)
}

// Item represents a stored component.
type Item interface {
	Name() string
	Content() (io.ReadCloser, error)
}

// PathFilter is used to exclude/include files in a FileSource.
type PathFilter func(string) bool

// FilterPrefix filters files by the given suffix.
func FilterSuffix(filter string) PathFilter {
	return func(s string) bool { return strings.HasSuffix(s, filter) }
}

// FilterPrefix filters files by the given prefix.
func FilterPrefix(filter string) PathFilter {
	return func(s string) bool { return strings.HasPrefix(s, filter) }
}

// MockSource is a static Source.
type MockSource struct {
	Items []MockItem
	i     int
}

// Next implements the Source interface.
func (m *MockSource) Next() (Item, error) {
	if m.i == len(m.Items) {
		return nil, nil
	}
	m.i++
	return m.Items[m.i-1], nil
}

// MockItem is a static Item.
type MockItem struct {
	ID       string
	Contents string
}

// Name implements the Item interface.
func (m MockItem) Name() string {
	return m.ID
}

// Content implements the Item interface.
func (m MockItem) Content() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBufferString(m.Contents)), nil
}
