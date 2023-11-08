package export

import (
	"encoding/json"
	"fmt"
	"io"
)

// Interface for CSV, JSON and other encoders
type Encoder interface {
	Encode(w io.Writer, v interface{}) error
}

type JSONEcoder struct{}

func (j *JSONEcoder) Encode(w io.Writer, v interface{}) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	return e.Encode(v)
}

func NewJSONEncoder() (*JSONEcoder, error) {
	return &JSONEcoder{}, nil
}

type PlainEncoder struct{}

func (p *PlainEncoder) Encode(w io.Writer, v interface{}) error {
	_, err := fmt.Fprint(w, v)
	return err
}

func NewPlainEncoder() (*PlainEncoder, error) {
	return &PlainEncoder{}, nil
}
