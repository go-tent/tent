package core

import (
	"io"
	"testing"
)

type MockDecoder struct {
	prefix string
	exts   []string
}

func (MockDecoder) GetID() string                                   { return "" }
func (m MockDecoder) Format() (string, []string)                    { return m.prefix, m.exts }
func (MockDecoder) Decode(_ string, _ io.Reader) (Component, error) { return nil, nil }
func (MockDecoder) Order() float64                                  { return 0 }
func (MockDecoder) Encode() ([]byte, error)                         { return nil, nil }

func TestDecodeAnalyze(t *testing.T) {
	testCases := map[bool][][]Component{
		true: {
			{
				MockDecoder{"d_", []string{".a"}},
				MockDecoder{"m_", []string{".a"}},
				MockDecoder{"s_", []string{".a"}},
				MockDecoder{"d_", []string{".b"}},
				MockDecoder{"m_", []string{".b"}},
				MockDecoder{"s_", []string{".b"}},
			}, {
				MockDecoder{"", []string{".a", ".b"}},
				MockDecoder{"", []string{".c", ".d"}},
			},
		},
		false: {
			{
				MockDecoder{"m_", []string{}},
			}, {
				MockDecoder{"", []string{".a"}},
			}, {
				MockDecoder{"m_", []string{".a", ".b"}},
			}, {
				MockDecoder{"m_", []string{".a"}},
				MockDecoder{"m_", []string{".a"}},
			}, {
				MockDecoder{"", []string{".a", ".b"}},
				MockDecoder{"", []string{".b", ".c"}},
			},
			{
				MockDecoder{"d_", []string{".a"}},
				MockDecoder{"m_", []string{".a"}},
				MockDecoder{"s_", []string{".a"}},
				MockDecoder{"", []string{".a", ".b"}},
				MockDecoder{"", []string{".c", ".d"}},
			},
			{
				MockDecoder{"", []string{".a", ".b"}},
				MockDecoder{"", []string{".c", ".d"}},
				MockDecoder{"d_", []string{".a"}},
				MockDecoder{"m_", []string{".a"}},
				MockDecoder{"s_", []string{".a"}},
			},
		}}

	for success, testList := range testCases {
		for i, tc := range testList {
			if err := detectCollisions(tc); (err == nil) != success {
				t.Fatalf("[Test %v] Expected %v, got %v", i, success, err)
			}
		}
	}
}
