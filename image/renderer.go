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

// A unit of work containing frame information.
type frameInfo struct {
	base   image.Rectangle
	frame  *image.Paletted
	index  int
	top    string
	bottom string
}

// Process each gif frame.
func processFrame(fi frameInfo, output chan frameInfo) {
	// Expand each frame, if needed, so it's the same size as the base.
	// This is to make it easier to draw and position the text.
	img := image.NewPaletted(fi.base, fi.frame.Palette)
	draw.Draw(img, fi.frame.Bounds(), fi.frame, fi.frame.Bounds().Min, draw.Src)

	// resImg, factor := resizeImage(img)
	// resBounds := resizeBounds(frame.Bounds(), factor)
	// src.Config.Width = resImg.Bounds().Dx()
	// src.Config.Height = resImg.Bounds().Dy()

	// Draw on the text.
	ctx := gfx.NewContext(img)
	if fi.top != "" {
		gfx.TopBanner(ctx, fi.top)
	}
	if fi.bottom != "" {
		gfx.BottomBanner(ctx, fi.bottom)
	}

	// Convert the context to a paletted image.
	img = image.NewPaletted(img.Bounds(), img.Palette)
	draw.FloydSteinberg.Draw(img, img.Bounds(), ctx.Image(), image.ZP)

	fi.frame = img.SubImage(fi.frame.Bounds()).(*image.Paletted)
	output <- fi
}

// RenderGif performs the graphical manipulation of the gif.
func RenderGif(opt cli.Options, st stream.Stream) stream.Stream {
	src := st.DecodeGif()
	queue := make(chan frameInfo)

	for x, frame := range src.Image {
		fi := frameInfo{
			base:   src.Image[0].Bounds(),
			frame:  frame,
			index:  x,
			top:    opt.Top,
			bottom: opt.Bottom,
		}
		go processFrame(fi, queue)
	}

	for range src.Image {
		fi := <-queue
		src.Image[fi.index] = fi.frame
	}

	close(queue)
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
