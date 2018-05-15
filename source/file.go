// Package source provides an iterable Source interface of elements with path and content.
package source

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NewFileSource returns a new FileSource, using the given filters.
func NewFileSource(ctx context.Context, root string, filters ...PathFilter) Source {
	src := FileSource{ch: make(chan fileItem)}
	go src.walk(ctx, root, filters)
	return src
}

// FileSource takes a directory in the filesystem as source.
type FileSource struct {
	ch chan fileItem
}

func (f FileSource) walk(ctx context.Context, root string, filters []PathFilter) {
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
			f.ch <- fileItem{root: root, Path: path, err: err}
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

// fileItem represent an Item in the filesystem.
type fileItem struct {
	Path string
	root string
	err  error
}

// Name implements the Item interface.
func (f fileItem) Name() string {
	return strings.Replace(f.Path[len(f.root)+1:], `\`, `/`, -1)
}

// Content implements the Item interface.
func (f fileItem) Content() (io.ReadCloser, error) {
	file, err := os.Open(filepath.Join(f.root, f.Name()))
	if err != nil {
		return nil, err
	}
	return file, err
}
