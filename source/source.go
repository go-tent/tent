package source

import (
	"context"
	"io"
	"os"
	"path/filepath"
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

// NewFileSource returns a new FileSource, using the given filters.
func NewFileSource(ctx context.Context, root string, filters ...PathFilter) Source {
	src := FileSource{ch: make(chan FileItem)}
	go src.walk(ctx, root, filters)
	return src
}

// FileSource takes a directory in the filesystem as source.
type FileSource struct {
	ch chan FileItem
}

func (f FileSource) walk(ctx context.Context, root string, filters []PathFilter) {
	done := ctx.Done()
	filepath.Walk(root, func(path string, _ os.FileInfo, err error) error {
		select {
		case <-done:
			return nil
		default:
			for _, f := range filters {
				if !f(path) {
					return nil
				}
			}
			f.ch <- FileItem{root: root, Path: path, err: err}
			return nil
		}
	})
	close(f.ch)
}

// Next implements the Source interface.
func (f FileSource) Next() (Item, error) {
	item, ok := <-f.ch
	if !ok {
		return nil, nil
	}
	if item.err != nil {
		return nil, item.err
	}
	return item, nil
}

// FileItem represent an Item in the filesystem.
type FileItem struct {
	Path string
	root string
	err  error
}

// Name implements the Item interface.
func (f FileItem) Name() string {
	return strings.Replace(f.Path[len(f.root)+1:], `\`, `/`, -1)
}

// Content implements the Item interface.
func (f FileItem) Content() (io.ReadCloser, error) {
	file, err := os.Open(filepath.Join(f.root, f.Name()))
	if err != nil {
		return nil, err
	}
	return file, err
}

// PathFilter is used to exclude/include files in a FileSource.
type PathFilter func(string) bool

// ExtFilter filters files by the given extension.
func ExtFilter(ext string) PathFilter {
	return func(s string) bool { return strings.HasSuffix(s, ext) }
}
