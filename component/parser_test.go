package component

import (
	"io"
	"strings"
	"testing"

	"gopkg.in/tent.v1/source"
)

type MockParser struct {
	prefix string
	exts   []string
}

func (m MockParser) Format() (string, []string)                   { return m.prefix, m.exts }
func (MockParser) Parse(_ string, _ io.Reader) (Component, error) { return nil, nil }

func TestParseAnalyze(t *testing.T) {
	testCases := map[bool][][]Parser{
		true: {
			{
				MockParser{"d_", []string{".a"}},
				MockParser{"m_", []string{".a"}},
				MockParser{"s_", []string{".a"}},
				MockParser{"d_", []string{".b"}},
				MockParser{"m_", []string{".b"}},
				MockParser{"s_", []string{".b"}},
			}, {
				MockParser{"", []string{".a", ".b"}},
				MockParser{"", []string{".c", ".d"}},
			},
		},
		false: {
			{
				MockParser{"m_", []string{}},
			}, {
				MockParser{"", []string{".a"}},
			}, {
				MockParser{"m_", []string{".a", ".b"}},
			}, {
				MockParser{"m_", []string{".a"}},
				MockParser{"m_", []string{".a"}},
			}, {
				MockParser{"", []string{".a", ".b"}},
				MockParser{"", []string{".b", ".c"}},
			},
			{
				MockParser{"d_", []string{".a"}},
				MockParser{"m_", []string{".a"}},
				MockParser{"s_", []string{".a"}},
				MockParser{"", []string{".a", ".b"}},
				MockParser{"", []string{".c", ".d"}},
			},
			{
				MockParser{"", []string{".a", ".b"}},
				MockParser{"", []string{".c", ".d"}},
				MockParser{"d_", []string{".a"}},
				MockParser{"m_", []string{".a"}},
				MockParser{"s_", []string{".a"}},
			},
		}}

	for success, testList := range testCases {
		for i, tc := range testList {
			if err := detectCollisions(tc); (err == nil) != success {
				t.Fatalf("[Test %v] Expected %v, got %v", i, success, err)
			}
		}
	}
}

func TestParseNestedCategories(t *testing.T) {
	src := source.MockSource{
		Items: []source.MockItem{
			{ID: "a/.category.yml", Contents: "index: 2\nm: x"},
			{ID: "a/b/.category.yml", Contents: "index: 7\nm: y"},
			{ID: "a/b/d/.category.yml", Contents: "index: 20\nm: w"},
			{ID: "a/b/c/.category.yml", Contents: "index: 12\nm: z"},
		},
	}
	r, err := Parse(&src)
	if err != nil {
		t.Fatal(err)
	}
	defer printCategory(t, 0, r)

	if l := len(r.Sub); l != 1 {
		t.Fatalf("Expected %d category, got %d", 1, l)
	}
	if id := r.Sub[0].ID; id != "a" {
		t.Fatalf("Expected %q category, got %q", "a", id)
	}
	if l := len(r.Sub[0].Sub); l != 1 {
		t.Fatalf("Expected %d category, got %d", 1, l)
	}
	if id := r.Sub[0].Sub[0].ID; id != "b" {
		t.Fatalf("Expected %q category, got %q", "b", id)
	}
	if l := len(r.Sub[0].Sub[0].Sub); l != 2 {
		t.Fatalf("Expected %d category, got %d", 2, l)
	}
	if id := r.Sub[0].Sub[0].Sub[0].ID; id != "c" {
		t.Fatalf("Expected %q category, got %q", "c", id)
	}
	if id := r.Sub[0].Sub[0].Sub[1].ID; id != "d" {
		t.Fatalf("Expected %q category, got %q", "d", id)
	}
}

func TestParseSegment(t *testing.T) {
	src := source.MockSource{
		Items: []source.MockItem{
			{ID: "s_a.md", Contents: `---
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
	defer printCategory(t, 0, r)

	if l := len(r.Components); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	s, ok := r.Components[0].(*Segment)
	if !ok {
		t.Fatalf("Expected component to be a %T, got %T", new(Segment), r.Components[0])
	}
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
}

func TestParseNestedSegment(t *testing.T) {
	src := source.MockSource{
		Items: []source.MockItem{
			{ID: "cat/s_a.md", Contents: `---
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
	defer printCategory(t, 0, r)

	if l := len(r.Sub); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	if l := len(r.Sub[0].Components); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	s, ok := r.Sub[0].Components[0].(*Segment)
	if !ok {
		t.Fatalf("Expected component to be a %T, got %T", new(Segment), r.Components[0])
	}
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
}

func TestParsePicture(t *testing.T) {
	src := source.MockSource{
		Items: []source.MockItem{
			{ID: "a.jpg", Contents: `somebytes`},
		},
	}
	r, err := Parse(&src)
	if err != nil {
		t.Fatal(err)
	}
	defer printCategory(t, 0, r)

	if l := len(r.Components); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	p, ok := r.Components[0].(*Picture)
	if !ok {
		t.Fatalf("Expected component to be a %T, got %T", new(Picture), r.Components[0])
	}
	if id := p.ID; id != "a.jpg" {
		t.Fatalf("Expected %q segments, got %q", "a.jpg", id)
	}
	if e := "somebytes"; string(p.Data) != e {
		t.Fatalf("Expected %v data, got %v", e, string(p.Data))
	}
}

func printCategory(t *testing.T, deep int, c *Category) {
	t.Logf("%s> %+v", strings.Repeat("  ", deep), c)
	for _, cat := range c.Sub {
		printCategory(t, deep+1, &cat)
	}
	for _, s := range c.Components {
		t.Logf("%s- %+v", strings.Repeat("  ", deep+1), s)
	}
}
