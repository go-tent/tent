package core

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Checks represents a list of Checkboxes.
type Checks struct {
	ID    string            `yaml:"-"`
	Index float64           `yaml:"index"`
	Meta  map[string]string `yaml:",inline"`
	List  []Check           `yaml:"list,omitempty"`
}

// GetID implements the Component interface.
func (c *Checks) GetID() string { return c.ID }

// Order implements the Component interface.
func (c *Checks) Order() float64 { return c.Index }

func (c Checks) String() string {
	return fmt.Sprintf("Checks:%v list:%v", c.ID, len(c.List))
}

// Encode returns Item contents.
func (c *Checks) Encode() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(b).Encode(c); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Format implements the Decoder interface.
func (Checks) Format() (string, []string) { return "c_", []string{".yml"} }

// Decode returns a new Checks with Item contents.
func (c *Checks) Decode(id string, r io.Reader) (Component, error) {
	return c.decode(id, r)
}
func (*Checks) decode(id string, r io.Reader) (*Checks, error) {
	c := Checks{ID: id}
	if err := yaml.NewDecoder(r).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Check is a checkbox.
type Check struct {
	Text     string  `yaml:"text"`
	Label    bool    `yaml:"label,omitempty"`
	Children []Check `yaml:"children,omitempty"`
}
