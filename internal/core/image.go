package core

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

const (
	formatPNG  = "png"
	formatJPEG = "jpeg"
)

// processImage decodes a PNG screenshot, optionally resizes it to maxWidth
// and re-encodes it as PNG or JPEG. It returns the encoded bytes and MIME type.
func processImage(src []byte, format string, quality, maxWidth int) ([]byte, string, error) {
	img, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, "", fmt.Errorf("decode screenshot: %w", err)
	}

	if maxWidth > 0 {
		bounds := img.Bounds()
		if bounds.Dx() > maxWidth {
			dstHeight := max(1, int(float64(bounds.Dy())*float64(maxWidth)/float64(bounds.Dx())))

			dst := image.NewRGBA(image.Rect(0, 0, maxWidth, dstHeight))
			draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
			img = dst
		}
	}

	var buf bytes.Buffer
	switch format {
	case formatJPEG:
		quality = max(1, min(100, quality))
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, "", fmt.Errorf("encode jpeg: %w", err)
		}
		return buf.Bytes(), "image/jpeg", nil
	default:
		if err := png.Encode(&buf, img); err != nil {
			return nil, "", fmt.Errorf("encode png: %w", err)
		}
		return buf.Bytes(), "image/png", nil
	}
}
