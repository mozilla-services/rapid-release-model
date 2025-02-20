package export

import (
	"io"
	"os"
)

type Exporter interface {
	Export(v interface{}) error
}

type WriterExporter struct {
	w       io.Writer
	encoder Encoder
}

func (s *WriterExporter) Export(v interface{}) error {
	return s.encoder.Encode(s.w, v)
}

func NewWriterExporter(w io.Writer, e Encoder) (*WriterExporter, error) {
	return &WriterExporter{w: w, encoder: e}, nil
}

type FileExporter struct {
	encoder  Encoder
	filename string
}

func (f *FileExporter) Export(v interface{}) error {
	file, err := os.Create(f.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return f.encoder.Encode(file, v)
}

func NewFileExporter(f string, e Encoder) (*FileExporter, error) {
	return &FileExporter{encoder: e, filename: f}, nil
}
