package renderer

import (
	"image"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/font"
	"github.com/nomad-software/meme/output"
)

const (
	FONT_BORDER_RADIUS = 3.0
	FONT_LEADING       = 1.4
	FONT_SIZE_MAX      = 80.0
	IMAGE_MARGIN       = 25.0
)

// Render the meme using the base image.
func Render(options cli.Options, base image.Image) image.Image {
	ctx := gg.NewContextForImage(base)

	if options.Top != "" {
		drawTopBanner(ctx, options.Top)
	}

	if options.Bottom != "" {
		drawBottomBanner(ctx, options.Bottom)
	}

	return ctx.Image()
}

// Draw the top text onto the meme.
func drawTopBanner(ctx *gg.Context, text string) {
	x := float64(ctx.Width()) / 2
	y := IMAGE_MARGIN
	drawText(ctx, text, x, y, 0.5, 0.0)
}

// Draw the bottom text onto the meme.
func drawBottomBanner(ctx *gg.Context, text string) {
	x := float64(ctx.Width()) / 2
	y := float64(ctx.Height()) - IMAGE_MARGIN
	drawText(ctx, text, x, y, 0.5, 1.0)
}

// Draw text onto the meme.
func drawText(ctx *gg.Context, text string, x float64, y float64, ax float64, ay float64) {
	text = strings.ToUpper(text)
	width := float64(ctx.Width()) - (IMAGE_MARGIN * 2)
	height := float64(ctx.Height()) / 3.5
	calculateFontSize(ctx, text, width, height)

	ctx.SetHexColor("#000")
	for angle := 0.0; angle < (2 * math.Pi); angle += 0.35 {
		bx := x + (math.Sin(angle) * FONT_BORDER_RADIUS)
		by := y + (math.Cos(angle) * FONT_BORDER_RADIUS)
		ctx.DrawStringWrapped(text, bx, by, ax, ay, width, FONT_LEADING, gg.AlignCenter)
	}

	ctx.SetHexColor("#FFF")
	ctx.DrawStringWrapped(text, x, y, ax, ay, width, FONT_LEADING, gg.AlignCenter)
}

// Dynamically calculate the correct size needed for text.
func calculateFontSize(ctx *gg.Context, text string, width float64, height float64) {
	for size := FONT_SIZE_MAX; size > 20; size -= 1 {
		var rWidth, rHeight float64
		var lWidth, lHeight float64

		err := ctx.LoadFontFace(font.Path, size)
		output.OnError(err, "Can not load font")
		lines := ctx.WordWrap(text, width)

		for _, line := range lines {
			lWidth, lHeight = ctx.MeasureString(line)

			if lWidth > rWidth {
				rWidth = lWidth
			}
		}

		rHeight = (lHeight * FONT_LEADING) * float64(len(lines))

		if rWidth <= width && rHeight <= height {
			break
		}
	}
}
