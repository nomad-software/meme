package image

import (
	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/image/draw"
	"github.com/nomad-software/meme/image/stream"
)

// Render the meme using the base image.
func Render(opt cli.Options, st stream.Stream) stream.Stream {

	base := Resize(st.DecodeImage())
	ctx := draw.NewContext(base)

	if opt.Top != "" {
		draw.TopBanner(ctx, opt.Top)
	}

	if opt.Bottom != "" {
		draw.BottomBanner(ctx, opt.Bottom)
	}

	return stream.EncodeImage(ctx.Image())
}
