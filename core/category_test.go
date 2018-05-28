package core

import (
	"bytes"
	"testing"
)

func TestCategory(t *testing.T) {
	c1 := Category{ID: "a", Index: 7, Meta: map[string]string{"title": "hello"}}
	b, err := c1.Encode()
	if err != nil {
		t.Fatal(err)
	}
	exp := "index: 7\ntitle: hello\n"
	if !bytes.Equal(b, []byte(exp)) {
		t.Fatalf("Expected %q, got %q", exp, string(b))
	}
	c2, err := (*Category).decode(nil, c1.ID, bytes.NewBufferString(exp))
	if err != nil {
		t.Fatal(err)
	}
	if c2.ID != c1.ID {
		t.Fatalf("Expected %q segments, got %q", c1.ID, c2.ID)
	}
	if c2.Meta["title"] != c1.Meta["title"] {
		t.Fatalf("Expected %v data, got %v", c1.Meta, c2.Meta)
	}
}
