package component

import (
	"fmt"
)

func detectCollisions(parsers []Parser) error {
	var taken = make(map[string]struct {
		unique Parser
		prefix map[string]Parser
	})
	for _, p := range parsers {
		pre, exts := p.Format()
		if err := validateFormat(pre, exts); err != nil {
			return fmt.Errorf("%T: %s", p, err)
		}
		if ext := exts[0]; len(exts) == 1 {
			if v, ok := taken[ext]; !ok {
				v.prefix = make(map[string]Parser)
				taken[ext] = v
			}
			if name := taken[ext].unique; name != nil {
				return fmt.Errorf("%T: ext %q used by %T", p, ext, name)
			}
			if name, ok := taken[ext].prefix[pre]; ok {
				return fmt.Errorf("%T: prefix %q, ext %q used by %T", p, pre, ext, name)
			}
			taken[ext].prefix[pre] = p
			continue
		}
		for _, ext := range exts {
			if taken[ext].prefix != nil {
				return fmt.Errorf("%T: ext %q used with prefix", p, ext)
			}
			if name := taken[ext].unique; name != nil {
				return fmt.Errorf("%T: ext %q used by %T", p, ext, name)
			}
			v := taken[ext]
			v.unique = p
			taken[ext] = v
		}
	}
	return nil
}

func validateFormat(pre string, exts []string) error {
	switch size, hasPre := len(exts), pre != ""; {
	case size == 0:
		return fmt.Errorf("no ext specified")
	case size == 1 && !hasPre:
		return fmt.Errorf("no prefix specified for unique ext")
	case size > 1 && hasPre:
		return fmt.Errorf("no prefix allowed for multiple ext")
	default:
		return nil
	}
}
