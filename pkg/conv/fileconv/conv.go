package fileconv

import "io"

type IImage interface {
	Thumbnail() (io.ReadWriter, error)
	Resize(width, height int) (io.ReadWriter, error)
	Compress(quality int) (io.ReadWriter, error)
}

type IDocument interface {
	Thumbnail() (io.ReadWriter, error)
	ConvertToPdf() (io.ReadWriter, error)
}

type IVideo interface {
	Thumbnail() (io.ReadWriter, error)
	Compress(quality int) (io.ReadWriter, error)
}

type IAudio interface {
	Compress(quality int) (io.ReadWriter, error)
}
