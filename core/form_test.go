package core

import (
	"bytes"
	"reflect"
	"testing"
)

func TestForm(t *testing.T) {
	f1 := &Form{
		ID:    "a",
		Index: 10,
		Meta:  map[string]string{"title": "form"},
		Screens: []FormScreen{
			{
				Meta: Map{"title": "1"},
				Items: []FormItem{
					{Name: "field1", Type: "text", Required: true, Meta: map[string]interface{}{"label": "Field 1"}},
					{Name: "field2", Type: "select", Meta: map[string]interface{}{"label": "Field 2", "option": []interface{}{"option1", "option2"}}},
				},
			},
			{
				Meta: Map{"title": "2"},
				Items: []FormItem{
					{Name: "field3", Type: "checkbox"},
					{Name: "field4", Type: "date"},
				},
			},
		},
	}
	b, err := f1.Encode()
	if err != nil {
		t.Fatal(err)
	}

	exp := `index: 10
screens:
- items:
  - name: field1
    type: text
    required: true
    label: Field 1
  - name: field2
    type: select
    label: Field 2
    option:
    - option1
    - option2
  title: "1"
- items:
  - name: field3
    type: checkbox
  - name: field4
    type: date
  title: "2"
title: form
`
	if !bytes.Equal(b, []byte(exp)) {
		t.Fatalf("Expected %q, got %q", exp, string(b))
	}
	f2, err := (*Form).decode(nil, f1.ID, bytes.NewBufferString(exp))
	if err != nil {
		t.Fatal(err)
	}
	if f2.ID != f1.ID {
		t.Fatalf("Expected %q checks, got %q", f1.ID, f2.ID)
	}
	if f2.Index != f1.Index {
		t.Fatalf("Expected %v index, got %v", f1.Index, f2.Index)
	}
	if len(f2.Meta) != 1 || f1.Meta["title"] != f2.Meta["title"] {
		t.Fatalf("Expected %v meta, got %v", f1.Meta, f2.Meta)
	}
	if l1, l2 := len(f1.Screens), len(f2.Screens); l1 != l2 {
		t.Fatalf("Expected %v screens, got %v", l1, l2)
	}
	for i := range f1.Screens {
		s1, s2 := f1.Screens[i], f2.Screens[i]
		if s1.Meta["title"] != s2.Meta["title"] {
			t.Fatalf("Expected %v screen Title, got %v", s1.Meta["title"], s2.Meta["title"])
		}
		if l1, l2 := len(s1.Items), len(s2.Items); l1 != l2 {
			t.Fatalf("Expected %v items, got %v", l1, l2)
		}
		for i := range s1.Items {
			if !reflect.DeepEqual(s1.Items[i], s2.Items[i]) {
				t.Fatalf("Expected %v item, got %v", s1.Items[i], s2.Items[i])
			}
		}
	}
}
