package component

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"gopkg.in/tent.v1/source"
)

func TestParseCategories(t *testing.T) {
	src := testSource{
		items: []testItem{
			{"contents_en/a/.category.yaml", "index: 2\nm: x"},
			{"contents_en/a/b/.category.yaml", "index: 7\nm: y"},
			{"contents_en/a/b/d/.category.yaml", "index: 20\nm: w"},
			{"contents_en/a/b/c/.category.yaml", "index: 12\nm: z"},
		},
	}
	r, err := Parse(&src)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(r.Sub); l != 1 {
		t.Fatalf("Expected %d locale, got %d", 1, l)
	}
	if r.Sub[0].ID != "en" {
		t.Fatalf("Expected %q locale, got %q", "en", r.Sub[0].ID)
	}
	if l := len(r.Sub[0].Sub); l != 1 {
		t.Fatalf("Expected %d category, got %d", 1, l)
	}
	if id := r.Sub[0].Sub[0].ID; id != "a" {
		t.Fatalf("Expected %q category, got %q", "a", id)
	}
	if l := len(r.Sub[0].Sub[0].Sub); l != 1 {
		t.Fatalf("Expected %d category, got %d", 1, l)
	}
	if id := r.Sub[0].Sub[0].Sub[0].ID; id != "b" {
		t.Fatalf("Expected %q category, got %q", "b", id)
	}
	if l := len(r.Sub[0].Sub[0].Sub[0].Sub); l != 2 {
		t.Fatalf("Expected %d category, got %d", 2, l)
	}
	if id := r.Sub[0].Sub[0].Sub[0].Sub[0].ID; id != "c" {
		t.Fatalf("Expected %q category, got %q", "c", id)
	}
	if id := r.Sub[0].Sub[0].Sub[0].Sub[1].ID; id != "d" {
		t.Fatalf("Expected %q category, got %q", "d", id)
	}
	printCategory(t, 0, r.Sub[0].Sub[0])
}

func TestParseSegments(t *testing.T) {
	src := testSource{
		items: []testItem{
			{"contents_en/a.md", `---
index: 10
title: segment title
---
# Title

text`},
		},
	}
	r, err := Parse(&src)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(r.Sub); l != 1 {
		t.Fatalf("Expected %d locale, got %d", 1, l)
	}
	if r.Sub[0].ID != "en" {
		t.Fatalf("Expected %q locale, got %q", "en", r.Sub[0].ID)
	}
	if l := len(r.Sub[0].Segments); l != 1 {
		t.Fatalf("Expected %d segments, got %d", 1, l)
	}
	s := r.Sub[0].Segments[0]
	if id := s.ID; id != "a" {
		t.Fatalf("Expected %q segments, got %q", "a", id)
	}
	if i := s.Index; i != 10 {
		t.Fatalf("Expected %v index, got %v", 10, i)
	}
	if e := map[string]string{"title": "segment title"}; len(s.Meta) != 1 || e["title"] != s.Meta["title"] {
		t.Fatalf("Expected %v meta, got %v", e, s.Meta)
	}

	if e := "# Title\n\ntext"; string(s.Body) != e {
		t.Fatalf("Expected %v body, got %v", e, string(s.Body))
	}
	printCategory(t, 0, r.Sub[0])
}

func printCategory(t *testing.T, deep int, c Category) {
	t.Logf("%s> %+v", strings.Repeat("  ", deep), c)
	for _, cat := range c.Sub {
		printCategory(t, deep+1, cat)
	}
	for _, s := range c.Segments {
		t.Logf("%s- %+v", strings.Repeat("  ", deep+1), s)
	}
}

type testSource struct {
	items []testItem
	i     int
}

func (t *testSource) Next() (source.Item, error) {
	if t.i == len(t.items) {
		return nil, nil
	}
	t.i++
	return t.items[t.i-1], nil
}

type testItem struct {
	name     string
	contents string
}

func (t testItem) Name() string {
	return t.name
}

func (t testItem) Content() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBufferString(t.contents)), nil
}
