package boards

import (
	"encoding/json"
	"image"
	"image/color"

	"github.com/benwiebe/udb-plugin-library/types"
)

// GradientBoard fills the display with a linear gradient between two configurable colours.
// It demonstrates BoardTypeStatic rendering without a datasource — the image is fully
// computed in Init() and returned unchanged on every Render() call.
//
// Optional JSON config fields:
//   - from:      start colour as a hex string, e.g. "#FF0000" (default black)
//   - to:        end colour as a hex string, e.g. "#0000FF" (default white)
//   - direction: "horizontal" (default), "vertical", or "diagonal"
type GradientBoard struct {
	id          string
	cachedImage image.Image
}

func NewGradientBoard(id string) *GradientBoard {
	return &GradientBoard{id: id}
}

func (b *GradientBoard) GetId() string   { return b.id }
func (b *GradientBoard) GetName() string { return "Gradient" }
func (b *GradientBoard) GetSupportedDimensions() []types.BoardDimensions {
	return []types.BoardDimensions{}
}
func (b *GradientBoard) GetType() types.BoardType  { return types.BoardTypeStatic }
func (b *GradientBoard) GetDatasourceType() string { return "" }

func (b *GradientBoard) Init(cfg json.RawMessage, _ types.Datasource[any], dimensions types.BoardDimensions) error {
	from := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	to := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	direction := "horizontal"

	if len(cfg) > 0 {
		var parsed struct {
			From      string `json:"from"`
			To        string `json:"to"`
			Direction string `json:"direction"`
		}
		if err := json.Unmarshal(cfg, &parsed); err != nil {
			return err
		}
		if parsed.From != "" {
			c, err := parseHexColour(parsed.From)
			if err != nil {
				return err
			}
			from = c
		}
		if parsed.To != "" {
			c, err := parseHexColour(parsed.To)
			if err != nil {
				return err
			}
			to = c
		}
		if parsed.Direction != "" {
			direction = parsed.Direction
		}
	}

	b.cachedImage = buildGradientImage(dimensions, from, to, direction)
	return nil
}

func (b *GradientBoard) Render() image.Image {
	return b.cachedImage
}

func buildGradientImage(dims types.BoardDimensions, from, to color.RGBA, direction string) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, dims.Width, dims.Height))

	// Avoid division by zero for 1×1 displays.
	w := max1(dims.Width - 1)
	h := max1(dims.Height - 1)

	for y := 0; y < dims.Height; y++ {
		for x := 0; x < dims.Width; x++ {
			var t float64
			switch direction {
			case "vertical":
				t = float64(y) / float64(h)
			case "diagonal":
				t = (float64(x)/float64(w) + float64(y)/float64(h)) / 2
			default: // "horizontal"
				t = float64(x) / float64(w)
			}
			img.SetRGBA(x, y, lerpRGBA(from, to, t))
		}
	}
	return img
}

func lerpRGBA(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		G: uint8(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		B: uint8(float64(a.B) + (float64(b.B)-float64(a.B))*t),
		A: uint8(float64(a.A) + (float64(b.A)-float64(a.A))*t),
	}
}

func max1(n int) int {
	if n < 1 {
		return 1
	}
	return n
}
