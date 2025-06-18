package helpers

import (
	"github.com/ohler55/ojg/oj"
)

func Marshal(v any) ([]byte, error) {
	return oj.Marshal(v)
}
func Unmarshal(data []byte, v any) error {
	return oj.Unmarshal(data, v)
}
