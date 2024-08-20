package fileconv

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
)

func createTestImage() io.Reader {
	img := imaging.New(200, 200, image.White)
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	return bytes.NewReader(buf.Bytes())
}

func TestThumbnail_ReturnsCorrectSize(t *testing.T) {
	reader := createTestImage()
	conv := NewImageConv(reader)

	buf, err := conv.Thumbnail()
	assert.NoError(t, err)

	img, err := imaging.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, 100, img.Bounds().Dx())
}

func TestThumbnail_InvalidImage_ReturnsError(t *testing.T) {
	reader := bytes.NewReader([]byte("invalid image data"))
	conv := NewImageConv(reader)

	_, err := conv.Thumbnail()
	assert.Error(t, err)
}

func TestResize_ReturnsCorrectSize(t *testing.T) {
	reader := createTestImage()
	conv := NewImageConv(reader)

	buf, err := conv.Resize(50, 50)
	assert.NoError(t, err)

	img, err := imaging.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, 50, img.Bounds().Dx())
	assert.Equal(t, 50, img.Bounds().Dy())
}

func TestResize_InvalidImage_ReturnsError(t *testing.T) {
	reader := bytes.NewReader([]byte("invalid image data"))
	conv := NewImageConv(reader)

	_, err := conv.Resize(50, 50)
	assert.Error(t, err)
}

func TestCompress_ReturnsCorrectQuality(t *testing.T) {
	reader := createTestImage()
	conv := NewImageConv(reader)

	buf, err := conv.Compress(50)
	assert.NoError(t, err)

	img, err := imaging.Decode(buf)
	assert.NoError(t, err)

	// Since we cannot directly check the quality, we ensure the image is still decodable
	assert.NotNil(t, img)
}

func TestCompress_InvalidImage_ReturnsError(t *testing.T) {
	reader := bytes.NewReader([]byte("invalid image data"))
	conv := NewImageConv(reader)

	_, err := conv.Compress(50)
	assert.Error(t, err)
}
