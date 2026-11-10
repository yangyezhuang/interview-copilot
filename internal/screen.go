package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/kbinani/screenshot"
)

const maxLongEdge = 1568

func resizeImage(src image.Image, limit int) image.Image {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w <= limit && h <= limit {
		return src
	}

	var newW, newH int
	if w >= h {
		newW = limit
		newH = int(float64(h) * float64(limit) / float64(w))
	} else {
		newH = limit
		newW = int(float64(w) * float64(limit) / float64(h))
	}
	if newW <= 0 {
		newW = 1
	}
	if newH <= 0 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	for y := 0; y < newH; y++ {
		srcY := bounds.Min.Y + (y * h / newH)
		for x := 0; x < newW; x++ {
			srcX := bounds.Min.X + (x * w / newW)
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

func CaptureScreen() (string, error) {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		return "", fmt.Errorf("no active displays found")
	}

	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", fmt.Errorf("failed to capture screen: %w", err)
	}

	scaled := resizeImage(img, maxLongEdge)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, scaled, &jpeg.Options{Quality: 75}); err != nil {
		return "", fmt.Errorf("failed to encode jpeg: %w", err)
	}

	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/jpeg;base64," + b64, nil
}
