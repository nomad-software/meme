package image

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/output"
)

const (
	UPLOAD_URL = "https://api.imgur.com/3/upload"
)

// Upload the image.
func Upload(opt cli.Options, img image.Image) string {
	stream := encode(img)
	base64 := base64.StdEncoding.EncodeToString(stream)

	return upload(opt, base64)
}

// Encode the image to a byte stream.
func encode(img image.Image) []byte {
	var buffer bytes.Buffer

	writer := bufio.NewWriter(&buffer)
	png.Encode(writer, img)
	writer.Flush()

	return buffer.Bytes()
}

// Perform the request to the storage provider.
func upload(opt cli.Options, base64 string) string {
	req, err := http.NewRequest("POST", UPLOAD_URL, strings.NewReader(base64))
	output.OnError(err, "Could not create upload request")
	req.Header.Set("Authorization", "Client-ID "+opt.ClientId)

	resp, err := http.DefaultClient.Do(req)
	output.OnError(err, "Could not upload image")
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		output.OnError(err, "Could not read response body")

		var imgur ImgurResponse
		err = json.Unmarshal(body, &imgur)
		output.OnError(err, "Could not decode json response")

		return imgur.Data.Link

	} else {
		output.Error("Could not upload image")
	}

	panic("Never reached")
}

type ImgurResponse struct {
	Status  int       `json:"status"`
	Success bool      `json:"success"`
	Data    ImgurData `json:"data"`
}

type ImgurData struct {
	Link string `json:"link"`
}
