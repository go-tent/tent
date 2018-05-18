package component

import (
	"bytes"
	"testing"
)

func TestSegment(t *testing.T) {
	s1 := &Segment{ID: "a", Index: 10, Meta: map[string]string{"title": "segment"}, Body: []byte("# Title\n\ntext")}
	b, err := s1.Encode()
	if err != nil {
		t.Fatal(err)
	}
	exp := "---\nindex: 10\ntitle: segment\n---\n# Title\n\ntext"
	if !bytes.Equal(b, []byte(exp)) {
		t.Fatalf("Expected %q, got %q", exp, string(b))
	}
	s2, err := segDecoder{}.decode(s1.ID, bytes.NewBufferString(exp))
	if err != nil {
		t.Fatal(err)
	}
	if s2.ID != s1.ID {
		t.Fatalf("Expected %q segments, got %q", s1.ID, s2.ID)
	}
	if s2.Index != s1.Index {
		t.Fatalf("Expected %v index, got %v", s1.Index, s2.Index)
	}
	if len(s2.Meta) != 1 || s1.Meta["title"] != s2.Meta["title"] {
		t.Fatalf("Expected %v meta, got %v", s1.Meta, s2.Meta)
	}
	if !bytes.Equal(s2.Body, s1.Body) {
		t.Fatalf("Expected %v body, got %v", string(s1.Body), string(s2.Body))
	}
}
