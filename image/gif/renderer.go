package gif

import (
	"image"
	"image/draw"

	"github.com/nomad-software/meme/cli"
	gfx "github.com/nomad-software/meme/image/draw"
	"github.com/nomad-software/meme/image/stream"
)

// Render the meme using the base image.
func Render(opt cli.Options, st stream.Stream) stream.Stream {

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
