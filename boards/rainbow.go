package boards

import (
	"encoding/json"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/benwiebe/udb-plugin-library/types"
)

// RainbowBoard renders a full-spectrum rainbow that scrolls across the display each frame.
// It demonstrates BoardTypeDynamic with no datasource: all state lives in the board itself,
// and each Render() call advances an internal hue offset to produce motion.
//
// Optional JSON config fields:
//   - speed:    degrees of hue shift per frame (default 2.0; negative reverses direction)
//   - vertical: if true, the rainbow runs top-to-bottom instead of left-to-right (default false)
type RainbowBoard struct {
	speed    float64
	vertical bool
	dims     types.BoardDimensions
	offset   float64
}

func NewRainbowBoard() *RainbowBoard {
	return &RainbowBoard{}
}

func (b *RainbowBoard) GetId() string   { return "rainbow" }
func (b *RainbowBoard) GetName() string { return "Rainbow" }
func (b *RainbowBoard) GetSupportedDimensions() []types.BoardDimensions {
	return []types.BoardDimensions{}
}
func (b *RainbowBoard) GetType() types.BoardType  { return types.BoardTypeDynamic }
func (b *RainbowBoard) GetDatasourceType() string { return "" }

func (b *RainbowBoard) Init(cfg json.RawMessage, _ types.Datasource, dimensions types.BoardDimensions) error {
	b.dims = dimensions
	b.speed = 2.0
	b.vertical = false

	if len(cfg) > 0 {
		var parsed struct {
			Speed    float64 `json:"speed"`
			Vertical bool    `json:"vertical"`
		}
		if err := json.Unmarshal(cfg, &parsed); err != nil {
			return err
		}
		if parsed.Speed != 0 {
			b.speed = parsed.Speed
		}
		b.vertical = parsed.Vertical
	}
	return nil
}

func (b *RainbowBoard) Render() types.AnimationFrame {
	img := image.NewRGBA(image.Rect(0, 0, b.dims.Width, b.dims.Height))

	span := float64(max1(b.dims.Width - 1))
	if b.vertical {
		span = float64(max1(b.dims.Height - 1))
	}

	for y := 0; y < b.dims.Height; y++ {
		for x := 0; x < b.dims.Width; x++ {
			pos := float64(x)
			if b.vertical {
				pos = float64(y)
			}
			hue := math.Mod(b.offset+(pos/span)*360, 360)
			img.SetRGBA(x, y, hsvToRGBA(hue, 1.0, 1.0))
		}
	}

	b.offset = math.Mod(b.offset+b.speed, 360)
	return types.AnimationFrame{Img: img, Duration: 33 * time.Millisecond}
}

// hsvToRGBA converts HSV (h in [0,360), s and v in [0,1]) to color.RGBA.
func hsvToRGBA(h, s, v float64) color.RGBA {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r, g, bl float64
	switch {
	case h < 60:
		r, g, bl = c, x, 0
	case h < 120:
		r, g, bl = x, c, 0
	case h < 180:
		r, g, bl = 0, c, x
	case h < 240:
		r, g, bl = 0, x, c
	case h < 300:
		r, g, bl = x, 0, c
	default:
		r, g, bl = c, 0, x
	}

	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((bl + m) * 255),
		A: 255,
	}
}
