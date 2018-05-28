package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

// Picture represents an image.
type Picture struct {
	ID   string
	Data []byte
}

// GetID implements the Component interface.
func (p *Picture) GetID() string { return p.ID }

// Encode returns Item contents.
func (p *Picture) Encode() ([]byte, error) {
	return p.Data, nil
}

// Order returns math.MaxFloat64, Pictures are shown last.
func (*Picture) Order() float64 { return math.MaxFloat64 }

func (p Picture) String() string {
	return fmt.Sprintf("Picture:%s Size:%v", p.ID, len(p.Data))
}

// Match implements the Decoder interface.
func (*Picture) Format() (string, []string) {
	return "", []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"}
}

// Decode returns a new Picture with Item contents.
func (p *Picture) Decode(id string, r io.Reader) (Component, error) {
	return p.decode(id, r)
}

func (*Picture) decode(id string, r io.Reader) (*Picture, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Picture{ID: id, Data: data}, nil
}
