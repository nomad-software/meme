package image

import (
	"image"
	"image/draw"

	"github.com/nfnt/resize"
	"github.com/nomad-software/meme/cli"
	gfx "github.com/nomad-software/meme/image/draw"
	"github.com/nomad-software/meme/image/stream"
)

const (
	maxImageSize = 600 // px
)

// Resize the passed image if it's too big.
func resizeImage(img image.Image) image.Image {
	if img.Bounds().Dx() > maxImageSize {
		img = resize.Resize(maxImageSize, 0, img, resize.Bilinear)
	}

	if img.Bounds().Dy() > maxImageSize {
		img = resize.Resize(0, maxImageSize, img, resize.Bilinear)
	}

	return img
}

// RenderImage performs the graphical manipulation of the image.
func RenderImage(opt cli.Options, st stream.Stream) stream.Stream {

	base := resizeImage(st.DecodeImage())
	ctx := gfx.NewContext(base)

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

	gif := st.DecodeGif()
	base := gif.Image[0]

	for x := 0; x < len(gif.Image); x++ {
		frame := gif.Image[x]

		// Resize each frame so it's the same as the base.
		img := image.NewPaletted(base.Bounds(), frame.Palette)
		// draw.Draw(img, img.Bounds(), image.Transparent, image.ZP, draw.Src)
		draw.Draw(img, frame.Bounds(), frame, frame.Bounds().Min, draw.Src)

		// Draw on it.
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
		gif.Image[x] = img.SubImage(frame.Bounds()).(*image.Paletted)
	}

	return stream.EncodeGif(gif)
}
