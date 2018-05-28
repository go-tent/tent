package core

import (
	"io"
	"log"
	"testing"

	"gopkg.in/tent.v1/item"
	"gopkg.in/tent.v1/source"
	"gopkg.in/yaml.v2"
)

func TestDecodeNestedCategories(t *testing.T) {
	items := []item.Memory{
		{ID: "a/.category.yml", Contents: []byte("index: 2\nm: x")},
		{ID: "a/b/.category.yml", Contents: []byte("index: 7\nm: y")},
		{ID: "a/b/d/.category.yml", Contents: []byte("index: 20\nm: w")},
		{ID: "a/b/c/.category.yml", Contents: []byte("index: 12\nm: z")},
	}
	r, err := NewRoot()
	if err != nil {
		log.Fatalf("%s", err)
	}
	if err := r.Decode(&source.Memory{Items: items}); err != nil {
		t.Fatal(err)
	}
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

func TestDecodeComponent(t *testing.T) {
	items := []item.Memory{
		{ID: "m_hello.mock", Contents: []byte(`index: 10`)},
	}
	r, err := NewRoot(mockCmp{})
	if err != nil {
		log.Fatalf("%s", err)
	}
	if err := r.Decode(&source.Memory{Items: items}); err != nil {
		t.Fatal(err)
	}

	if l := len(r.Components); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	m, ok := r.Components[0].(mockCmp)
	if !ok {
		t.Fatalf("Expected component to be a %T, got %T", new(mockCmp), r.Components[0])
	}
	if id := m.ID; id != "hello" {
		t.Fatalf("Expected %q segments, got %q", "a", id)
	}
	if i := m.Index; i != 10 {
		t.Fatalf("Expected %v index, got %v", 10, i)
	}
}

func TestDecodeNestedComponent(t *testing.T) {
	items := []item.Memory{
		{ID: "cat/m_hello.mock", Contents: []byte(`index: 10`)},
	}
	r, err := NewRoot(mockCmp{})
	if err != nil {
		log.Fatalf("%s", err)
	}
	if err := r.Decode(&source.Memory{Items: items}); err != nil {
		t.Fatal(err)
	}

	if l := len(r.Sub); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	if l := len(r.Sub[0].Components); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	m, ok := r.Sub[0].Components[0].(mockCmp)
	if !ok {
		t.Fatalf("Expected component to be a %T, got %T", mockCmp{}, r.Sub[0].Components[0])
	}
	if id := m.ID; id != "hello" {
		t.Fatalf("Expected %q segments, got %q", "a", id)
	}
	if i := m.Index; i != 10 {
		t.Fatalf("Expected %v index, got %v", 10, i)
	}
}

type mockCmp struct {
	ID    string
	Index float64
}

func (m mockCmp) GetID() string  { return m.ID }
func (m mockCmp) Order() float64 { return m.Index }

func (mockCmp) Encode() ([]byte, error) { return nil, nil }

func (mockCmp) Format() (string, []string) { return "m_", []string{".mock"} }

func (mockCmp) Decode(id string, r io.Reader) (Component, error) {
	m := mockCmp{ID: id}
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
