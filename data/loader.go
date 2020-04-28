package data

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"strconv"

	"github.com/ahmedaly113/roadbook/data/internal"
	. "github.com/ahmedaly113/roadbook/model"
	"github.com/ahmedaly113/roadbook/ui"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func resizePNG(src []byte, imgSize uint) ([]byte, error) {
	i, _, err := image.Decode(bytes.NewReader(src))
	b := i.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y

	if err != nil {
		return nil, errors.Wrap(err, "problem decoding PNG during resize")
	}

	if imgWidth > imgHeight {
		i = resize.Resize(imgSize, 0, i, resize.Bicubic)
	} else {
		i = resize.Resize(0, imgSize, i, resize.Bicubic)
	}

	var w bytes.Buffer
	err = png.Encode(&w, i)
	if err != nil {
		return nil, errors.Wrap(err, "unable to encode PNG during resize")
	}

	return w.Bytes(), nil
}

func FromGPX(src []byte) (*Model, error) {
	data := &Model{
		Speed:   float32(rand.Intn(25) + 25),
		Heading: float32(rand.Intn(180) + 30),
		Idx:     0,
		IsSpeedZone: false,
		SpeedLimit: 0.0,
	}

	gpx := internal.GPX{}
	err := xml.Unmarshal(src, &gpx)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing GPX")
	}

	for i, w := range gpx.Waypoints {
		distance, err := strconv.ParseFloat(w.Distance, 32)
		if err != nil {
			distance = -1
		}

		s, err := internal.CheatBase64(w.Tulip)
		if err != nil {
			log.Printf("Tulip raw: %s\n", w.Tulip)
			return nil, errors.Wrap(err, fmt.Sprintf("unable to parse tulip on waypoint index %d", i))
		}

		tulip, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, errors.Wrap(err, "invalid base64 in SVG")
		}

		tulipResized, err := resizePNG(tulip, ui.TulipHeight)
		if err != nil {
			return nil, errors.Wrap(err, "error resizing image during GPX loading")
		}

		s, err = internal.CheatBase64(w.Notes)
		if err != nil {
			log.Printf("Tulip raw: %s\n", w.Tulip)
			return nil, errors.Wrap(err, fmt.Sprintf("unable to parse notes on waypoint index %d", i))
		}

		notes, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, errors.Wrap(err, "invalid base64 in SVG")
		}

		notesResized, err := resizePNG(notes, ui.NotesHeight)
		if err != nil {
			return nil, errors.Wrap(err, "error resizing image during GPX loading")
		}

		data.Book = append(data.Book, Waypoint{
			Background: color.White,
			Distance:   float32(distance),
			Tulip:      tulipResized,
			Notes:      notesResized,
			DZ:         w.DZ != nil,
			FZ:         w.FZ != nil,
		})
	}
	return data, nil
}
