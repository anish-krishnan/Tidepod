package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/h2non/filetype"
)

// These constants define the common MIME types of data
// MIME type is a label used to identify a type of data and
// allows tidepod to transcode different types of media
const (
	MimeTypeJpeg   = "image/jpeg"
	MimeTypePng    = "image/png"
	MimeTypeGif    = "image/gif"
	MimeTypeBitmap = "image/bmp"
	MimeTypeTiff   = "image/tiff"
	MimeTypeHEIF   = "image/heif"
)

var mimeTable = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/gif":  "gif",
	"image/bmp":  "bmp",
	"image/tiff": "tiff",
	"image/heif": "heif",
}

func getExtFromMimeType(mimeType string) string {
	ext, ok := mimeTable[mimeType]
	if !ok {
		panic(fmt.Sprintf("can't map mimetype to extension. mimetype: %s", mimeType))
	}
	return ext
}

// MimeType returns the mime type of a file, empty string if unknown.
func getMimeType(filename string) string {
	handle, err := os.Open("photo_storage/saved/" + filename)

	if err != nil {
		panic(err)
	}

	defer handle.Close()

	// Only the first 261 bytes are used to sniff the content type.
	buffer := make([]byte, 261)

	if _, err := handle.Read(buffer); err != nil {
		panic(err)
		return ""
	} else if t, err := filetype.Get(buffer); err == nil && t != filetype.Unknown {
		return t.MIME.Value
	} else if t := filetype.GetType(NormalizedExt(filename)); t != filetype.Unknown {
		return t.MIME.Value
	} else {
		return ""
	}
}

// NormalizedExt returns the file extension without dot and in lowercase.
func NormalizedExt(fileName string) string {
	if dot := strings.LastIndex(fileName, "."); dot != -1 && len(fileName[dot+1:]) >= 1 {
		return strings.ToLower(fileName[dot+1:])
	}

	return ""
}
