package utils

import (
	"bytes"
	"fmt"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

func ConvertToGIF(reader io.Reader, fileExt string) (io.Reader, error) {
	var img image.Image
	var err error
	switch fileExt {
	case ".png":
		img, err = png.Decode(reader)
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(reader)
	case ".webp":
		img, err = webp.Decode(reader)
	case ".gif":
		return reader, nil
	default:
		return nil, fmt.Errorf("non supported file extension: %s", fileExt)
	}

	if err != nil {
		return nil, fmt.Errorf("error while decoding: %v", err)
	}

	var gifBuffer bytes.Buffer
	if err := gif.Encode(&gifBuffer, img, nil); err != nil {
		return nil, fmt.Errorf("error while encoding: %v", err)
	}

	return bytes.NewReader(gifBuffer.Bytes()), nil

}
