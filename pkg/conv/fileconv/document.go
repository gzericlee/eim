package fileconv

import (
	"io"
)

type DocumentConv struct {
	Reader io.Reader
}

func NewDocumentConv(reader io.Reader) *DocumentConv {
	return &DocumentConv{Reader: reader}
}

func (its *DocumentConv) Thumbnail() (io.ReadWriter, error) {
	return nil, nil
}

func (its *DocumentConv) ConvertToPdf() (io.ReadWriter, error) {
	return nil, nil
}
