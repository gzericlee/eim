package fileconv

import (
	"bytes"
	"fmt"
	"io"

	"github.com/disintegration/imaging"
	_ "github.com/disintegration/imaging"
)

type ImageConv struct {
	Reader io.Reader
}

func NewImageConv(reader io.Reader) *ImageConv {
	return &ImageConv{Reader: reader}
}

func (its *ImageConv) Thumbnail() (io.ReadWriter, error) {
	img, err := imaging.Decode(its.Reader, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("imaging decode -> %w", err)
	}

	const thumbnailWidth = 100
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()
	aspectRatio := float64(originalHeight) / float64(originalWidth)
	newHeight := int(aspectRatio * float64(thumbnailWidth))
	thumbnail := imaging.Thumbnail(img, thumbnailWidth, newHeight, imaging.Lanczos)

	buf := bytes.NewBuffer(nil)
	err = imaging.Encode(buf, thumbnail, imaging.JPEG)
	if err != nil {
		return nil, fmt.Errorf("imaging encode -> %w", err)
	}

	return buf, err
}

func (its *ImageConv) Resize(width, height int) (io.ReadWriter, error) {
	img, err := imaging.Decode(its.Reader, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("imaging decode -> %w", err)
	}

	thumbnail := imaging.Thumbnail(img, width, height, imaging.Lanczos)

	buf := bytes.NewBuffer(nil)
	err = imaging.Encode(buf, thumbnail, imaging.JPEG)
	if err != nil {
		return nil, fmt.Errorf("imaging encode -> %w", err)
	}

	return buf, err
}

func (its *ImageConv) Compress(quality int) (io.ReadWriter, error) {
	img, err := imaging.Decode(its.Reader, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("imaging decode -> %w", err)
	}

	buf := bytes.NewBuffer(nil)
	err = imaging.Encode(buf, img, imaging.JPEG, imaging.JPEGQuality(quality))
	if err != nil {
		return nil, fmt.Errorf("imaging encode -> %w", err)
	}
	return buf, nil
}
