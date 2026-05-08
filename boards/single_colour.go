package boards

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"strings"

	"github.com/benwiebe/udb-plugin-library/types"
)

// SingleColourBoard is a trivial example of a UDB board which displays a single colour
// on the entire display.
type SingleColourBoard struct {
	id     string
	colour color.Color
}

// NewSingleColourBoard creates a SingleColourBoard with the given ID.
func NewSingleColourBoard(id string) *SingleColourBoard {
	return &SingleColourBoard{id: id}
}

func (b *SingleColourBoard) GetId() string {
	return b.id
}

func (b *SingleColourBoard) GetName() string {
	return "Single Colour"
}

func (b *SingleColourBoard) GetSupportedDimensions() []types.BoardDimensions {
	return []types.BoardDimensions{} // Deliberately empty: supports all sizes
}

func (b *SingleColourBoard) GetType() types.BoardType {
	return types.BoardTypeStatic
}

func (b *SingleColourBoard) GetDatasourceType() string {
	return "" // no datasource required
}

// Init accepts an optional config JSON with a "colour" field containing a hex colour string
// (e.g. "#FF0000" or "FF0000"). Defaults to white if omitted.
func (b *SingleColourBoard) Init(cfg json.RawMessage, datasource types.Datasource[any]) error {
	if len(cfg) > 0 {
		var parsed struct {
			Colour string `json:"colour"`
		}
		if err := json.Unmarshal(cfg, &parsed); err != nil {
			return err
		}
		if parsed.Colour != "" {
			c, err := parseHexColour(parsed.Colour)
			if err != nil {
				return err
			}
			b.colour = c
		}
	}
	if b.colour == nil {
		b.colour = color.White
	}
	return nil
}

func (b *SingleColourBoard) Render(dimensions types.BoardDimensions) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, dimensions.Width, dimensions.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{b.colour}, image.Point{}, draw.Src)
	return img
}

// parseHexColour parses a "#RRGGBB" or "RRGGBB" hex string into a color.RGBA.
func parseHexColour(s string) (color.RGBA, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex colour %q: expected 6 hex digits", s)
	}
	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid hex colour %q: %w", s, err)
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid hex colour %q: %w", s, err)
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid hex colour %q: %w", s, err)
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}
