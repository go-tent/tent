// Package source provides an iterable Source interface of elements with path and content.
package source

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/tent.v1/item"
)

// NewFile returns a new File source, using the given filters.
func NewFile(ctx context.Context, root string, filters ...PathFilter) File {
	src := File{ch: make(chan fileItem)}
	go src.walk(ctx, root, filters)
	return src
}

// File takes a directory in the filesystem as source.
type File struct {
	ch chan fileItem
}

func (f File) walk(ctx context.Context, root string, filters []PathFilter) {
	done := ctx.Done()
	filepath.Walk(root, func(path string, _ os.FileInfo, err error) error {
		select {
		case <-done:
			return nil
		default:
			info, _ := os.Stat(path)
			if info.IsDir() {
				return nil
			}
			for _, f := range filters {
				if !f(path) {
					return nil
				}
			}
			f.ch <- fileItem{Root: root, Path: path, Err: err}
			return nil
		}
	})
	close(f.ch)
}

// Next implements the Source interface.
func (f File) Next() (item.Item, error) {
	item, ok := <-f.ch
	if !ok {
		return nil, nil
	}
	if item.Err != nil {
		return nil, item.Err
	}
	return item, nil
}

// fileItem represent an Item in the filesystem.
type fileItem struct {
	Path string
	Root string
	Err  error
}

// Name implements the Item interface.
func (f fileItem) Name() string {
	return strings.Replace(f.Path[len(f.Root)+1:], `\`, `/`, -1)
}

// Content implements the Item interface.
func (f fileItem) Content() (io.ReadCloser, error) {
	file, err := os.Open(filepath.Join(f.Root, f.Name()))
	if err != nil {
		return nil, err
	}
	return file, nil
}
