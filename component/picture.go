package component

import (
	"bytes"
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

// Encode returns Item contents.
func (p *Picture) Encode() (io.Reader, error) {
	return bytes.NewBuffer(p.Data), nil
}

// Order returns math.MaxFloat64, Pictures are shown last.
func (*Picture) Order() float64 { return math.MaxFloat64 }

func (p Picture) String() string {
	return fmt.Sprintf("Picture:%s Size:%v", p.ID, len(p.Data))
}

// picDecoder is the Decoder for Picture.
type picDecoder struct{}

// Match implements the Decoder interface.
func (picDecoder) Format() (string, []string) {
	return "", []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"}
}

// Decode populates the Picture with Item contents.
func (p picDecoder) Decode(id string, r io.Reader) (Component, error) {
	return p.decode(id, r)
}

func (picDecoder) decode(id string, r io.Reader) (*Picture, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Picture{ID: id, Data: data}, nil
}
