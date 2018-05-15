package source

import (
	"testing"
)

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
