package component

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
func (picDecoder) Decode(id string, r io.Reader) (Component, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Picture{ID: id, Data: data}, nil
}
