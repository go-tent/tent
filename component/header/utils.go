// Package header contains accessories functions for Component header and metada parsing.
package header

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

// ParseMeta decodes metadata from a Reader, if sep it looks for "---\n" delimiters.
func ParseMeta(r *bufio.Reader, m *map[string]string, sep bool) error {
	if sep {
		row, err := r.ReadBytes('\n')
		if err != nil {
			return err
		}
		if !bytes.Equal([]byte("---\n"), row) {
			return errors.New("Invalid header")
		}
	}
	var header io.Reader = r
	if sep {
		b := bytes.NewBuffer(nil)
		for {
			row, err := r.ReadBytes('\n')
			if err != nil {
				return err
			}
			if bytes.Equal([]byte("---\n"), row) {
				break
			}
			b.Write(row)
		}
		header = b
	}
	return yaml.NewDecoder(header).Decode(m)
}

// ParseIndex extracts the Index from metadata
func ParseIndex(m map[string]string, idx *float64) error {
	index, err := strconv.ParseFloat(m["index"], 64)
	if err != nil {
		return fmt.Errorf("Invalid Index: %s", m["index"])
	}
	delete(m, "index")
	*idx = index
	return nil
}
