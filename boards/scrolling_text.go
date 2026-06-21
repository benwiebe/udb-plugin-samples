package boards

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/benwiebe/udb-plugin-library/types"
)

// ScrollingTextBoard scrolls a line of text horizontally across the display in a
// continuous loop. It demonstrates BoardTypeDynamic with internal animation state:
// the text is pre-rendered to an off-screen buffer in Init(), and each Render() call
// simply advances an x-offset and blits the visible slice onto the output image.
//
// Optional JSON config fields:
//   - text:       the string to scroll (default "Hello, World!")
//   - colour:     text colour as a hex string (default white)
//   - background: background colour as a hex string (default black)
//   - speed:      pixels to advance per frame (default 2)
type ScrollingTextBoard struct {
	speed     int
	dims      types.BoardDimensions
	textBuf   *image.RGBA // full pre-rendered text row, height == dims.Height
	bgUniform *image.Uniform
	offset    int // x position of the text's left edge relative to the display
}

func NewScrollingTextBoard() *ScrollingTextBoard {
	return &ScrollingTextBoard{}
}

func (b *ScrollingTextBoard) GetId() string   { return "scrolling-text" }
func (b *ScrollingTextBoard) GetName() string { return "Scrolling Text" }
func (b *ScrollingTextBoard) GetSupportedDimensions() []types.BoardDimensions {
	return []types.BoardDimensions{}
}
func (b *ScrollingTextBoard) GetType() types.BoardType  { return types.BoardTypeDynamic }
func (b *ScrollingTextBoard) GetDatasourceType() string { return "" }

func (b *ScrollingTextBoard) Init(cfg json.RawMessage, _ types.Datasource, dimensions types.BoardDimensions) error {
	text := "Hello, World!"
	textColour := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	bgColour := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	b.speed = 2
	b.dims = dimensions

	if len(cfg) > 0 {
		var parsed struct {
			Text       string `json:"text"`
			Colour     string `json:"colour"`
			Background string `json:"background"`
			Speed      int    `json:"speed"`
		}
		if err := json.Unmarshal(cfg, &parsed); err != nil {
			return err
		}
		if parsed.Text != "" {
			text = parsed.Text
		}
		if parsed.Colour != "" {
			c, err := parseHexColour(parsed.Colour)
			if err != nil {
				return err
			}
			textColour = c
		}
		if parsed.Background != "" {
			c, err := parseHexColour(parsed.Background)
			if err != nil {
				return err
			}
			bgColour = c
		}
		if parsed.Speed != 0 {
			b.speed = parsed.Speed
		}
	}

	b.bgUniform = &image.Uniform{bgColour}

	// Find the largest font that fits within the display height (with a small margin).
	marginV := int(float64(dimensions.Height) * 0.1)
	innerH := dimensions.Height - 2*marginV
	fontSize := findScrollFontSize(text, innerH)
	face := truetype.NewFace(clockFont, &truetype.Options{Size: fontSize, DPI: 72})
	defer face.Close()

	// Measure the rendered text width so we can allocate exactly the right buffer.
	drawer := &font.Drawer{Face: face}
	textW := drawer.MeasureString(text).Ceil()
	m := face.Metrics()
	ascent := m.Ascent.Ceil()

	// Render text onto a buffer exactly as wide as the text itself.
	b.textBuf = image.NewRGBA(image.Rect(0, 0, textW, dimensions.Height))
	draw.Draw(b.textBuf, b.textBuf.Bounds(), &image.Uniform{bgColour}, image.Point{}, draw.Src)

	drawer = &font.Drawer{
		Dst:  b.textBuf,
		Src:  &image.Uniform{textColour},
		Face: face,
		Dot:  fixed.P(0, marginV+ascent),
	}
	drawer.DrawString(text)

	// Start with the text just off the right edge of the display.
	b.offset = dimensions.Width
	return nil
}

func (b *ScrollingTextBoard) Render() types.AnimationFrame {
	img := image.NewRGBA(image.Rect(0, 0, b.dims.Width, b.dims.Height))
	draw.Draw(img, img.Bounds(), b.bgUniform, image.Point{}, draw.Src)

	textW := b.textBuf.Bounds().Dx()

	// Blit the visible portion of the text buffer onto the output image.
	// offset is the x coordinate of the text's left edge in display space.
	srcX := clamp0(-b.offset, textW-1)
	dstX := clamp0(b.offset, b.dims.Width-1)
	copyW := min2(b.dims.Width-dstX, textW-srcX)
	if copyW > 0 {
		src := b.textBuf.SubImage(image.Rect(srcX, 0, srcX+copyW, b.dims.Height))
		draw.Draw(img, image.Rect(dstX, 0, dstX+copyW, b.dims.Height), src, image.Pt(srcX, 0), draw.Src)
	}

	b.offset -= b.speed
	if b.offset <= -textW {
		b.offset = b.dims.Width
	}

	return types.AnimationFrame{Img: img, Duration: 33 * time.Millisecond}
}

// findScrollFontSize returns the largest font size whose rendered text height fits within maxH.
func findScrollFontSize(text string, maxH int) float64 {
	lo, hi := 1.0, float64(maxH)*2
	for hi-lo > 0.5 {
		mid := (lo + hi) / 2
		face := truetype.NewFace(clockFont, &truetype.Options{Size: mid, DPI: 72})
		m := face.Metrics()
		h := (m.Ascent + m.Descent).Ceil()
		face.Close()
		if h <= maxH {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo
}

func clamp0(v, max int) int {
	if v < 0 {
		return 0
	}
	if v > max {
		return max
	}
	return v
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}
