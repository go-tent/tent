package component

import (
	"bytes"
	"reflect"
	"testing"
)

func TestChecks(t *testing.T) {
	c1 := &Checks{
		ID:    "a",
		Index: 10,
		Meta:  map[string]string{"title": "sample checklist"},
		List: []Check{
			{
				Text:  "fruits",
				Label: true,
				Children: []Check{
					{Text: "apple"},
					{Text: "pear"},
					{Text: "melon"},
					{
						Text:  "citrus",
						Label: true,
						Children: []Check{
							{Text: "lemon"},
							{Text: "orange"},
						},
					},
				},
			},
		},
	}
	b, err := c1.Encode()
	if err != nil {
		t.Fatal(err)
	}
	exp := `index: 10
list:
- text: fruits
  label: true
  children:
  - text: apple
  - text: pear
  - text: melon
  - text: citrus
    label: true
    children:
    - text: lemon
    - text: orange
title: sample checklist
`
	if !bytes.Equal(b, []byte(exp)) {
		t.Fatalf("Expected %q, got %q", exp, string(b))
	}
	c2, err := checksDecoder{}.decode(c1.ID, bytes.NewBufferString(exp))
	if err != nil {
		t.Fatal(err)
	}
	if c2.ID != c1.ID {
		t.Fatalf("Expected %q checks, got %q", c1.ID, c2.ID)
	}
	if c2.Index != c1.Index {
		t.Fatalf("Expected %v index, got %v", c1.Index, c2.Index)
	}
	if len(c2.Meta) != 1 || c1.Meta["title"] != c2.Meta["title"] {
		t.Fatalf("Expected %v meta, got %v", c1.Meta, c2.Meta)
	}
	if !reflect.DeepEqual(c2.List, c1.List) {
		t.Fatalf("Expected List:\n%v\nGot:\n%v", c1.List, c2.List)
	}
}
