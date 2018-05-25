package component

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Form represent a Web form
type Form struct {
	ID      string            `yaml:"-"`
	Index   float64           `yaml:"index"`
	Meta    map[string]string `yaml:",inline"`
	Screens []FormScreen      `yaml:"screens"`
}

// Order implements the Component interface.
func (f *Form) Order() float64 {
	return f.Index
}

func (f Form) String() string {
	return fmt.Sprintf("Form:%v screens:%v", f.ID, len(f.Screens))
}

// Encode returns Item contents.
func (f *Form) Encode() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(b).Encode(f); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// FormScreen is a Form Screen
type FormScreen struct {
	Name  string     `yaml:"name"`
	Items []FormItem `yaml:"items"`
}

// FormItem is form input
type FormItem struct {
	Name     string                 `yaml:"name"`
	Type     string                 `yaml:"type"`
	Required bool                   `yaml:"required,omitempty"`
	Meta     map[string]interface{} `yaml:",inline"`
}

// formDecoder is the Decoder for Form.
type formDecoder struct{}

// Format implements the Decoder interface.
func (formDecoder) Format() (string, []string) { return "f_", []string{".yml"} }

// Decode populates the Segment with Item contents.
func (s formDecoder) Decode(id string, r io.Reader) (Component, error) {
	return s.decode(id, r)
}
func (formDecoder) decode(id string, r io.Reader) (*Form, error) {
	c := Form{ID: id}
	if err := yaml.NewDecoder(r).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
