package image

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"math"
	"math/rand"
	"time"

	"github.com/nfnt/resize"
	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/data"
	gfx "github.com/nomad-software/meme/image/draw"
	"github.com/nomad-software/meme/image/stream"
)

const (
	maxImageSize   = 650 // px
	shakeDelay     = 2   // 100ths of a second
	shakeFrames    = 10  // Number of frames to create when shaking static images
	shakeIntensity = 8   // px
)

// RenderImage performs the graphical manipulation of the image.
func RenderImage(opt cli.Options, st stream.Stream) stream.Stream {
	if opt.Trigger {
		st = shake(st)
		st = trigger(st)
	} else if opt.Shake {
		st = shake(st)
	}

	if opt.Trigger || opt.Shake || (opt.Gif && st.IsGif()) {
		st = renderGif(opt, st)
	} else {
		st = renderImage(opt, st)
	}

	return st
}

// Trigger adds the triggered banner.
func trigger(st stream.Stream) stream.Stream {
	src := st.DecodeGif()
	width := uint(src.Config.Width + (shakeIntensity * 2))
	ds := Decal(data.TriggeredDecal)
	decal := resize.Resize(width, 0, ds.DecodeImage(), resize.NearestNeighbor)
	queue := make(chan triggerInfo)

	for x, frame := range src.Image {
		ti := triggerInfo{
			frame: frame,
			decal: decal,
			index: x,
		}
		go processTrigger(ti, queue)
	}

	for range src.Image {
		ti := <-queue
		src.Image[ti.index] = ti.frame
	}

	close(queue)
	return stream.EncodeGif(src)
}

// A unit of work containing frame information for triggering gifs.
type triggerInfo struct {
	frame *image.Paletted
	decal image.Image
	index int
}

// Draw the triggered banner on the frame.
func processTrigger(ti triggerInfo, output chan triggerInfo) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rx := r.Intn(shakeIntensity*2) - shakeIntensity
	ry := r.Intn(shakeIntensity*2) - shakeIntensity

	fHeight := ti.frame.Bounds().Dy()
	dHeight := ti.decal.Bounds().Dy()

	x := shakeIntensity + rx
	y := (fHeight - (dHeight - shakeIntensity)) + ry

	point := image.Point{x, -y}

	draw.FloydSteinberg.Draw(ti.frame, ti.frame.Bounds(), ti.decal, point)
	output <- ti
}

// Shake randomly shakes an image.
func shake(st stream.Stream) stream.Stream {
	if st.IsGif() {
		return shakeGif(st)
	}
	return shakeImage(st)
}

// Create a random point for shaking the image.
func pointShaker() func() image.Point {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func() image.Point {
		x := r.Intn(shakeIntensity + 1)
		y := r.Intn(shakeIntensity + 1)
		return image.Point{x, y}
	}
}

// Create a rectangle to crop the shaking image to not expose the borders.
func shakeBounds(b image.Rectangle) image.Rectangle {
	x := b.Bounds().Max.X - shakeIntensity*2
	y := b.Bounds().Max.Y - shakeIntensity*2
	return image.Rect(0, 0, x, y)
}

// shakeGif randomly shakes a gif.
// This function can't use concurrency because normalising the gif needs a
// shared base image for all frames.
func shakeGif(st stream.Stream) stream.Stream {
	src := st.DecodeGif()
	base := src.Image[0]
	images := make([]*image.Paletted, len(src.Image))
	delays := make([]int, len(src.Image))
	crop := shakeBounds(base.Bounds())
	shakePoint := pointShaker()

	for x, frame := range src.Image {
		// Draw on the base.
		draw.Draw(base, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)

		// Create a new frame.
		img := image.NewPaletted(base.Bounds(), frame.Palette)
		draw.Draw(img, img.Bounds(), base, shakePoint(), draw.Src)

		images[x] = img.SubImage(crop).(*image.Paletted)
		delays[x] = shakeDelay
	}

	dst := &gif.GIF{
		Image: images,
		Delay: delays,
	}

	return stream.EncodeGif(dst)
}

// A unit of work containing frame information for shaking an image.
type shakeImgInfo struct {
	img   image.Image
	frame *image.Paletted
	point image.Point
	crop  image.Rectangle
	index int
}

// shakeImage randomly shakes an image creating a gif animation.
// This function can use concurrency because nothing is shared between frames.
func shakeImage(st stream.Stream) stream.Stream {
	src := st.DecodeImage()
	images := make([]*image.Paletted, shakeFrames)
	delays := make([]int, shakeFrames)
	crop := shakeBounds(src.Bounds())
	queue := make(chan shakeImgInfo)
	shakePoint := pointShaker()

	for x := range images {
		si := shakeImgInfo{
			img:   src,
			point: shakePoint(),
			crop:  crop,
			index: x,
		}
		go processImageShake(si, queue)
	}

	for range images {
		si := <-queue
		images[si.index] = si.frame
		delays[si.index] = shakeDelay
	}

	close(queue)

	dst := &gif.GIF{
		Image: images,
		Delay: delays,
	}

	return stream.EncodeGif(dst)
}

// Process shaking an image.
func processImageShake(si shakeImgInfo, output chan shakeImgInfo) {
	img := image.NewPaletted(si.img.Bounds(), palette.Plan9)
	draw.FloydSteinberg.Draw(img, img.Bounds(), si.img, si.point)
	si.frame = img.SubImage(si.crop).(*image.Paletted)

	output <- si
}

// RenderImage performs the graphical manipulation of the image.
func renderImage(opt cli.Options, st stream.Stream) stream.Stream {
	img := st.DecodeImage()
	img = reduceImage(img, maxImageSize)

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

// A unit of work containing frame information for drawing text.
type drawInfo struct {
	bounds image.Rectangle
	frame  *image.Paletted
	index  int
	top    string
	bottom string
}

// RenderGif performs the graphical manipulation of the gif.
func renderGif(opt cli.Options, st stream.Stream) stream.Stream {
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
	draw.Draw(img, img.Bounds(), ctx.Image(), image.ZP, draw.Src)

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

	// No reduction necessary.
	if fctr == 1.0 {
		return src
	}

	base := image.NewRGBA(first.Bounds())

	for x, frame := range src.Image {
		rs := reduceInfo{
			config: &src.Config,
			base:   base,
			frame:  frame,
			fctr:   fctr,
			index:  x,
		}
		src.Image[x] = processFrameReduction(rs)
	}

	return src
}

// Process resizing each gif frame at max quality.
func processFrameReduction(rs reduceInfo) *image.Paletted {
	if rs.index == 0 {
		resBounds := reduceBounds(rs.frame.Bounds(), rs.fctr)
		rs.config.Width = resBounds.Dx()
		rs.config.Height = resBounds.Dy()
	}

	if rs.frame.Bounds().Dx() == 0 || rs.frame.Bounds().Dy() == 0 {
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
func reduceBounds(src image.Rectangle, factor float64) image.Rectangle {
	x0 := int(math.Floor(float64(src.Min.X) * factor))
	y0 := int(math.Floor(float64(src.Min.Y) * factor))
	x1 := int(math.Floor(float64(src.Max.X) * factor))
	y1 := int(math.Floor(float64(src.Max.Y) * factor))

	return image.Rect(x0, y0, x1, y1)
}
