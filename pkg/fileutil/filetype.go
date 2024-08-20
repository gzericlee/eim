package fileutil

import "github.com/h2non/filetype"

func IsImage(buf []byte) bool {
	return filetype.IsImage(buf)
}

func IsDocument(buf []byte) bool {
	return filetype.IsDocument(buf)
}

func IsVideo(buf []byte) bool {
	return filetype.IsVideo(buf)
}

func IsAudio(buf []byte) bool {
	return filetype.IsAudio(buf)
}
