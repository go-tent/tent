package source

import (
	"context"
	"os"
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
	baseTest(context.Background(), t, 2)
}

func TestFileSourceContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	baseTest(ctx, t, 0)
}

func baseTest(ctx context.Context, t *testing.T, expected int) {
	var count int
	src := NewFileSource(ctx, wd, ExtFilter(".go"))
	for item, err := src.Next(); item != nil; item, err = src.Next() {
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v", item.Name())
		r, err := item.Content()
		if err != nil {
			t.Error("open", item.Name(), err)
		} else {
			r.Close()
		}
		count++
	}
	if count != expected {
		t.Errorf("Expected %d items, got %d", expected, count)
	}
}
