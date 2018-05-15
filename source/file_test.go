package source

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var wd string

func init() {
	v, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wd = v
}

func TestFileSource(t *testing.T) {
	baseTest(t, NewFileSource(context.Background(), wd, FilterSuffix(".go"), func(s string) bool {
		return strings.HasPrefix(filepath.Base(s), "file")
	}), 2)
}

func TestFileSourceContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	baseTest(t, NewFileSource(ctx, wd), 0)
}
