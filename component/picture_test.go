package component

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestPicture(t *testing.T) {
	p1 := &Picture{ID: "a.jpg", Data: []byte("picbytes")}
	r, err := p1.Encode()
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	exp := "picbytes"
	if !bytes.Equal(b, []byte(exp)) {
		t.Fatalf("Expected %q, got %q", exp, string(b))
	}
	p2, err := picDecoder{}.decode(p1.ID, bytes.NewBufferString(exp))
	if err != nil {
		t.Fatal(err)
	}
	if p2.ID != p1.ID {
		t.Fatalf("Expected %q segments, got %q", p1.ID, p2.ID)
	}
	if !bytes.Equal(p2.Data, p1.Data) {
		t.Fatalf("Expected %v data, got %v", string(p1.Data), string(p2.Data))
	}
}
