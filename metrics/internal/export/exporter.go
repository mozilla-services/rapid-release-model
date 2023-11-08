package export

import (
	"io"
	"os"
)

type Exporter interface {
	Export(v interface{}) error
}

type WriterExporter struct {
	W       io.Writer
	Encoder Encoder
}

func (s *WriterExporter) Export(v interface{}) error {
	return s.Encoder.Encode(s.W, v)
}

type FileExporter struct {
	Filename string
	Encoder  Encoder
}

func (f *FileExporter) Export(v interface{}) error {
	file, err := os.Create(f.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return f.Encoder.Encode(file, v)
}
