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

// GetID implements the Component interface.
func (f *Form) GetID() string { return f.ID }

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

// Format implements the Decoder interface.
func (*Form) Format() (string, []string) { return "f_", []string{".yml"} }

// Decode returns a new Form with Item contents.
func (f *Form) Decode(id string, r io.Reader) (Component, error) {
	return f.decode(id, r)
}

func (*Form) decode(id string, r io.Reader) (*Form, error) {
	c := Form{ID: id}
	if err := yaml.NewDecoder(r).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
