package image

import (
	"image"
	"image/draw"
	"math"

	"github.com/nfnt/resize"
	"github.com/nomad-software/meme/cli"
	gfx "github.com/nomad-software/meme/image/draw"
	"github.com/nomad-software/meme/image/stream"
)

const (
	maxImageSize = 600 // px
)

// Resize the passed image if it's too big.
func resizeImage(img image.Image) (image.Image, float64) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	var factor float64 = 1.0

	if width > height && width > maxImageSize {
		img = resize.Resize(maxImageSize, 0, img, resize.NearestNeighbor)
		factor = float64(maxImageSize) / float64(width)

	} else if height > width && height > maxImageSize {
		img = resize.Resize(0, maxImageSize, img, resize.NearestNeighbor)
		factor = float64(maxImageSize) / float64(width)
	}

	return img, factor
}

// RenderImage performs the graphical manipulation of the image.
func RenderImage(opt cli.Options, st stream.Stream) stream.Stream {
	img, _ := resizeImage(st.DecodeImage())

	// Draw on the text.
	ctx := gfx.NewContext(img)
	if opt.Top != "" {
		gfx.TopBanner(ctx, opt.Top)
	}
	if opt.Bottom != "" {
		gfx.BottomBanner(ctx, opt.Bottom)
	}

	return stream.EncodeImage(ctx.Image())
}

// RenderGif performs the graphical manipulation of the gif.
func RenderGif(opt cli.Options, st stream.Stream) stream.Stream {
	src := st.DecodeGif()
	base := src.Image[0]

	for x, frame := range src.Image {
		// Expand each frame, if needed, so it's the same size as the base.
		// This is to make it easier to draw and position the text.
		img := image.NewPaletted(base.Bounds(), frame.Palette)
		draw.Draw(img, frame.Bounds(), frame, frame.Bounds().Min, draw.Src)

		// resImg, factor := resizeImage(img)
		// resBounds := resizeBounds(frame.Bounds(), factor)
		// src.Config.Width = resImg.Bounds().Dx()
		// src.Config.Height = resImg.Bounds().Dy()

		// Draw on the text.
		ctx := gfx.NewContext(img)
		if opt.Top != "" {
			gfx.TopBanner(ctx, opt.Top)
		}
		if opt.Bottom != "" {
			gfx.BottomBanner(ctx, opt.Bottom)
		}

		// Convert the context to a paletted image.
		img = image.NewPaletted(img.Bounds(), frame.Palette)
		draw.FloydSteinberg.Draw(img, img.Bounds(), ctx.Image(), image.ZP)

		// Replace the frame.
		src.Image[x] = img.SubImage(frame.Bounds()).(*image.Paletted)
	}

	return stream.EncodeGif(src)
}

// Recalculate bound sizes from a passed factor.
func resizeBounds(src image.Rectangle, factor float64) image.Rectangle {
	x0 := int(math.Ceil(float64(src.Min.X) * factor))
	y0 := int(math.Ceil(float64(src.Min.Y) * factor))
	x1 := int(math.Floor(float64(src.Max.X) * factor))
	y1 := int(math.Floor(float64(src.Max.Y) * factor))

	return image.Rect(x0, y0, x1, y1)
}
