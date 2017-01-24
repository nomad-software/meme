package image

import (
	"image"

	"github.com/nfnt/resize"
)

const (
	maxImageSize = 600 // px
)

// Resize the passed image if it's too big.
func Resize(img image.Image) image.Image {
	if img.Bounds().Dx() > maxImageSize {
		img = resize.Resize(maxImageSize, 0, img, resize.Bilinear)
	}

	if img.Bounds().Dy() > maxImageSize {
		img = resize.Resize(0, maxImageSize, img, resize.Bilinear)
	}

	return img
}
