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
	baseTest(t, NewFileSource(context.Background(), wd, ExtFilter(".go")), 4)
}

func TestFileSourceContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	baseTest(t, NewFileSource(ctx, wd, ExtFilter(".go")), 0)
}

func baseTest(t *testing.T, src Source, expected int) {
	var count int
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
