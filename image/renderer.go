package image

import (
	"image"
	"image/draw"
	"image/gif"
	"math"

	"github.com/nfnt/resize"
	"github.com/nomad-software/meme/cli"
	gfx "github.com/nomad-software/meme/image/draw"
	"github.com/nomad-software/meme/image/stream"
)

const (
	maxImageSize = 600 // px
)

// reduceImage will resize an image if any of its dimensions are above the passed max
// size.
func reduceImage(img image.Image, maxSize uint) image.Image {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	if w > h && w > int(maxSize) {
		img = resize.Resize(maxSize, 0, img, resize.NearestNeighbor)
	} else if h > w && h > int(maxSize) {
		img = resize.Resize(0, maxSize, img, resize.NearestNeighbor)
	}

	return img
}

// RenderImage performs the graphical manipulation of the image.
func RenderImage(opt cli.Options, st stream.Stream) stream.Stream {
	img := reduceImage(st.DecodeImage(), maxImageSize)

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

// A unit of work containing frame information for drawing text.
type drawInfo struct {
	bounds image.Rectangle
	frame  *image.Paletted
	index  int
	top    string
	bottom string
}

// RenderGif performs the graphical manipulation of the gif.
func RenderGif(opt cli.Options, st stream.Stream) stream.Stream {
	src := st.DecodeGif()
	src = reduceGif(opt, src, maxImageSize)
	queue := make(chan drawInfo)

	for x, frame := range src.Image {
		fi := drawInfo{
			bounds: src.Image[0].Bounds(),
			frame:  frame,
			index:  x,
			top:    opt.Top,
			bottom: opt.Bottom,
		}
		go processFrameDraw(fi, queue)
	}

	for range src.Image {
		fi := <-queue
		src.Image[fi.index] = fi.frame
	}

	close(queue)
	return stream.EncodeGif(src)
}

// Process drawing on each gif frame.
func processFrameDraw(fi drawInfo, output chan drawInfo) {
	// Expand each frame, if needed, so it's the same size as the base.
	// This is to make it easier to draw and position the text.
	img := image.NewPaletted(fi.bounds, fi.frame.Palette)
	draw.Draw(img, fi.frame.Bounds(), fi.frame, fi.frame.Bounds().Min, draw.Src)

	// Draw on the text.
	ctx := gfx.NewContext(img)
	if fi.top != "" {
		gfx.TopBanner(ctx, fi.top)
	}
	if fi.bottom != "" {
		gfx.BottomBanner(ctx, fi.bottom)
	}

	// Convert the graphic context to a paletted image.
	img = image.NewPaletted(img.Bounds(), img.Palette)
	draw.FloydSteinberg.Draw(img, img.Bounds(), ctx.Image(), image.ZP)

	fi.frame = img.SubImage(fi.frame.Bounds()).(*image.Paletted)
	output <- fi
}

// A unit of work containing frame information for reducing a gif frame.
type reduceInfo struct {
	config *image.Config
	base   *image.RGBA
	frame  *image.Paletted
	fctr   float64
	index  int
}

// reduceGif will resize a gif if any of its dimensions are above the passed max
// size.
func reduceGif(opt cli.Options, src *gif.GIF, maxSize int) *gif.GIF {
	first := src.Image[0]
	fctr := calcFactor(first, maxSize)
	queue := make(chan reduceInfo)
	base := image.NewRGBA(first.Bounds())

	for x, frame := range src.Image {
		rs := reduceInfo{
			config: &src.Config,
			base:   base,
			frame:  frame,
			fctr:   fctr,
			index:  x,
		}
		if opt.MaxAnim {
			src.Image[x] = processFrameResizeMax(rs)
		} else {
			go processFrameResize(rs, queue)
		}
	}

	if !opt.MaxAnim {
		for range src.Image {
			rs := <-queue
			src.Image[rs.index] = rs.frame
		}
	}

	close(queue)
	return src
}

// Process resizing each gif frame at max quality.
func processFrameResizeMax(rs reduceInfo) *image.Paletted {
	if rs.fctr == 1.0 {
		return rs.frame // No reduction needed.
	}

	if rs.index == 0 {
		resBounds := calcBounds(rs.frame.Bounds(), rs.fctr)
		rs.config.Width = resBounds.Dx()
		rs.config.Height = resBounds.Dy()
	}

	if rs.base.Bounds().Dx() == 0 && rs.frame.Bounds().Dy() == 0 {
		return rs.frame // Empty frame, don't change or can cause corruption.
	}

	// Draw over the base.
	draw.Draw(rs.base, rs.frame.Bounds(), rs.frame, rs.frame.Bounds().Min, draw.Over)

	// Resize the base to the required size.
	w := uint(rs.config.Width)
	h := uint(rs.config.Height)
	res := resize.Resize(w, h, rs.base, resize.NearestNeighbor)

	// Create a new frame.
	img := image.NewPaletted(res.Bounds(), rs.frame.Palette)
	draw.Draw(img, res.Bounds(), res, image.ZP, draw.Src)

	return img
}

// Process resizing each gif frame.
func processFrameResize(rs reduceInfo, output chan reduceInfo) {
	if rs.fctr == 1.0 {
		output <- rs // No reduction needed.
		return
	}

	resBounds := calcBounds(rs.frame.Bounds(), rs.fctr)

	if rs.index == 0 {
		rs.config.Width = resBounds.Dx()
		rs.config.Height = resBounds.Dy()
	}

	if resBounds.Dx() == 0 && resBounds.Dy() == 0 {
		output <- rs // Empty frame, don't change or can cause corruption.
		return
	}

	w := uint(resBounds.Bounds().Dx())
	h := uint(resBounds.Bounds().Dy())

	// Resize the frame.
	res := resize.Resize(w, h, rs.frame, resize.NearestNeighbor)
	img := image.NewPaletted(resBounds, rs.frame.Palette)
	draw.Draw(img, img.Bounds(), res, image.ZP, draw.Src)

	rs.frame = img
	output <- rs
}

// Calculate the reduction factor from a desired maximum size.
func calcFactor(img image.Image, maxSize int) float64 {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	fctr := 1.0

	if w > h && w > maxSize {
		fctr = float64(maxSize) / float64(w)

	} else if h > w && h > maxSize {
		fctr = float64(maxSize) / float64(w)
	}

	return fctr
}

// Recalculate bounds sizes using a passed factor.
func calcBounds(src image.Rectangle, fctr float64) image.Rectangle {
	x0 := int(math.Floor(float64(src.Min.X) * fctr))
	y0 := int(math.Floor(float64(src.Min.Y) * fctr))
	x1 := int(math.Floor(float64(src.Max.X) * fctr))
	y1 := int(math.Floor(float64(src.Max.Y) * fctr))

	return image.Rect(x0, y0, x1, y1)
}
