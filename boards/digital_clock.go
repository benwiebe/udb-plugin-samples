package boards

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"

	"github.com/benwiebe/udb-plugin-library/types"
)

var clockFont *truetype.Font

func init() {
	var err error
	clockFont, err = truetype.Parse(gomono.TTF)
	if err != nil {
		panic("failed to parse embedded clock font: " + err.Error())
	}
}

// DigitalClockBoard displays the current time from a CurrentTimeDatasource.
// Optional JSON config fields:
//   - format:     Go time layout string (default "15:04"; use "15:04:05" to show seconds)
//   - colour:     hex colour for the digits, e.g. "#FF0000" (default white)
//   - blinkColon: toggle colon separators off on odd seconds (default true)
type DigitalClockBoard struct {
	id         string
	datasource types.Datasource[any]
	format     string
	colour     color.Color
	blinkColon bool

	// Pre-computed in Init from the display dimensions and config.
	cachedDims types.BoardDimensions
	cachedFace font.Face
	cachedDot  fixed.Point26_6 // pre-computed baseline origin
	cachedSrc  *image.Uniform  // pre-allocated colour source
}

func NewDigitalClockBoard(id string) *DigitalClockBoard {
	return &DigitalClockBoard{id: id}
}

func (b *DigitalClockBoard) GetId() string   { return b.id }
func (b *DigitalClockBoard) GetName() string { return "Digital Clock" }
func (b *DigitalClockBoard) GetSupportedDimensions() []types.BoardDimensions {
	return []types.BoardDimensions{}
}
func (b *DigitalClockBoard) GetType() types.BoardType  { return types.BoardTypeDynamic }
func (b *DigitalClockBoard) GetDatasourceType() string { return "UdbSamplePlugin/CurrentTime" }

func (b *DigitalClockBoard) Init(cfg json.RawMessage, datasource types.Datasource[any], dimensions types.BoardDimensions) error {
	b.datasource = datasource
	b.format = "15:04"
	b.colour = color.White
	b.blinkColon = true

	if len(cfg) > 0 {
		var parsed struct {
			Format     string `json:"format"`
			Colour     string `json:"colour"`
			BlinkColon *bool  `json:"blinkColon"`
		}
		if err := json.Unmarshal(cfg, &parsed); err != nil {
			return err
		}
		if parsed.Format != "" {
			b.format = parsed.Format
		}
		if parsed.Colour != "" {
			c, err := parseHexColour(parsed.Colour)
			if err != nil {
				return err
			}
			b.colour = c
		}
		if parsed.BlinkColon != nil {
			b.blinkColon = *parsed.BlinkColon
		}
	}

	b.buildCache(dimensions)
	return nil
}

func (b *DigitalClockBoard) Render() types.AnimationFrame {
	t := b.datasource.GetData().(time.Time)

	timeStr := t.Format(b.format)
	if b.blinkColon && t.Second()%2 == 1 {
		timeStr = strings.ReplaceAll(timeStr, ":", " ")
	}

	img := image.NewRGBA(image.Rect(0, 0, b.cachedDims.Width, b.cachedDims.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	d := font.Drawer{
		Dst:  img,
		Src:  b.cachedSrc,
		Face: b.cachedFace,
		Dot:  b.cachedDot,
	}
	d.DrawString(timeStr)

	// Use 500ms when blinkColon is on or when the format includes seconds — sub-second
	// precision is needed in both cases. Otherwise 1s is sufficient.
	dur := time.Second
	if b.blinkColon || strings.Contains(b.format, "05") {
		dur = 500 * time.Millisecond
	}
	return types.AnimationFrame{Img: img, Duration: dur}
}

// buildCache computes and stores the font face, baseline origin, and colour source for the
// given dimensions. Called once from Init.
func (b *DigitalClockBoard) buildCache(dimensions types.BoardDimensions) {
	if b.cachedFace != nil {
		b.cachedFace.Close()
	}
	b.cachedDims = dimensions

	marginX := int(float64(dimensions.Width) * 0.10)
	marginY := int(float64(dimensions.Height) * 0.10)
	inner := image.Rect(marginX, marginY, dimensions.Width-marginX, dimensions.Height-marginY)

	// Use the Go time reference moment so the sample string matches the exact layout
	// (same character count and structure as any real formatted time).
	ref := time.Date(0, 1, 1, 15, 4, 5, 0, time.UTC).Format(b.format)
	fontSize := findClockFontSize(ref, inner.Dx(), inner.Dy())

	b.cachedFace = truetype.NewFace(clockFont, &truetype.Options{Size: fontSize, DPI: 72})

	d := font.Drawer{Face: b.cachedFace}
	textWidth := d.MeasureString(ref).Ceil()
	m := b.cachedFace.Metrics()

	x := inner.Min.X + (inner.Dx()-textWidth)/2
	y := inner.Min.Y + (inner.Dy()-(m.Ascent+m.Descent).Ceil())/2 + m.Ascent.Ceil()
	b.cachedDot = fixed.P(x, y)
	b.cachedSrc = image.NewUniform(b.colour)
}

// findClockFontSize binary-searches for the largest font size where text fits within maxWidth×maxHeight.
func findClockFontSize(text string, maxWidth, maxHeight int) float64 {
	lo, hi := 1.0, float64(maxHeight)*2
	for hi-lo > 0.5 {
		mid := (lo + hi) / 2
		face := truetype.NewFace(clockFont, &truetype.Options{Size: mid, DPI: 72})
		d := font.Drawer{Face: face}
		w := d.MeasureString(text).Ceil()
		m := face.Metrics()
		h := (m.Ascent + m.Descent).Ceil()
		face.Close()
		if w <= maxWidth && h <= maxHeight {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo
}
