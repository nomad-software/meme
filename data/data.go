package data

import (
	"embed"
)

const (
	// ImagePath is the path to built-in templates.
	ImagePath = "images"

	// ImageExtension is the file extension of the built-in templates.
	ImageExtension = ".jpg"

	// Font is the location of the built-in font.
	Font = "fonts/impact.ttf"

	// TriggeredDecal is the banner used on triggered memes.
	TriggeredDecal = "decals/triggered.jpg"
)

//go:embed decals/*
//go:embed fonts/*
//go:embed images/*
var Files embed.FS
