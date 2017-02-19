package stream

import (
	"bytes"
	"image"
	"image/gif"
	"image/png"
	"io"
	"io/ioutil"

	"github.com/nomad-software/meme/output"
)

// Stream contains information about a loaded image.
type Stream struct {
	io.Reader
	bytes []byte
	index int
	typ   string
}

// Bytes returns the stream's bytes.
func (st *Stream) Bytes() []byte {
	return st.bytes
}

// Read implements the io.Reader interface for the Stream.
func (st *Stream) Read(b []byte) (n int, err error) {
	if st.index >= len(st.bytes) {
		return 0, io.EOF
	}
	n = copy(b, st.bytes[st.index:])
	st.index += n
	return
}

// IsGif returns true if the loaded image is a gif.
func (st *Stream) IsGif() bool {
	return st.typ == "gif"
}

// IsJpg returns true if the loaded image is a Jpeg.
func (st *Stream) IsJpg() bool {
	return st.typ == "jpeg"
}

// IsPng returns true if the loaded image is a Png.
func (st *Stream) IsPng() bool {
	return st.typ == "png"
}

// FileExt returns the file extension of the image.
func (st *Stream) FileExt() string {
	if st.IsGif() {
		return "gif"
	}
	if st.IsJpg() {
		return "jpg"
	}
	if st.IsPng() {
		return "png"
	}
	panic("File extension not recognised")
}

// NewStream creates a new stream.
func NewStream(stream io.Reader) Stream {
	a, err := ioutil.ReadAll(stream)
	output.OnError(err, "Could not read image bytes")

	b := make([]byte, len(a))
	copy(b, a)

	_, typ, err := image.DecodeConfig(bytes.NewReader(a))
	output.OnError(err, "Could not decode image config")

	return Stream{
		bytes: b,
		typ:   typ,
	}
}

// EncodeImage encodes an image into a stream.
func EncodeImage(img image.Image) Stream {
	var buffer bytes.Buffer
	png.Encode(&buffer, img)
	return NewStream(&buffer)
}

// DecodeImage decodes the byte stream and returns an image.
func (st *Stream) DecodeImage() image.Image {
	img, _, err := image.Decode(st)
	output.OnError(err, "Could not decode image")
	return img
}

// EncodeGif encodes a gif into a stream.
func EncodeGif(img *gif.GIF) Stream {
	var buffer bytes.Buffer
	gif.EncodeAll(&buffer, img)
	return NewStream(&buffer)
}

// DecodeGif decodes the byte stream and returns a gif.
func (st *Stream) DecodeGif() *gif.GIF {
	if !st.IsGif() {
		output.Error("Can't decode stream to gif")
	}
	gif, err := gif.DecodeAll(st)
	output.OnError(err, "Could not decode gif")
	return gif
}
