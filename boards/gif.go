package boards

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"time"

	"github.com/benwiebe/udb-plugin-library/types"
)

//go:embed assets/sample.gif
var defaultGIFBytes []byte

// GifBoard plays an animated GIF on the display.
// It demonstrates BoardTypeAnimated: all frames are decoded and composited in Init(),
// producing a pre-baked []types.AnimationFrame that udb-core cycles through without
// calling Render() again. Frame durations are taken directly from the GIF delay table.
//
// Optional JSON config fields:
//   - path: filesystem path to a GIF file; uses the embedded sample if omitted
type GifBoard struct {
	id     string
	frames types.Animation
}

func NewGifBoard(id string) *GifBoard {
	return &GifBoard{id: id}
}

func (b *GifBoard) GetId() string                                   { return b.id }
func (b *GifBoard) GetName() string                                 { return "Animated GIF" }
func (b *GifBoard) GetSupportedDimensions() []types.BoardDimensions { return []types.BoardDimensions{} }
func (b *GifBoard) GetType() types.BoardType                        { return types.BoardTypeAnimated }
func (b *GifBoard) GetDatasourceType() string                       { return "" }

func (b *GifBoard) Init(cfg json.RawMessage, _ types.Datasource, dimensions types.BoardDimensions) error {
	path := ""
	if len(cfg) > 0 {
		var parsed struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal(cfg, &parsed); err != nil {
			return err
		}
		path = parsed.Path
	}

	g, err := decodeGIF(path)
	if err != nil {
		return err
	}

	b.frames, err = buildFrames(g, dimensions)
	return err
}

func (b *GifBoard) Render() types.Animation {
	return b.frames
}

// decodeGIF decodes the GIF at path, or the embedded sample when path is empty.
func decodeGIF(path string) (*gif.GIF, error) {
	if path == "" {
		g, err := gif.DecodeAll(bytes.NewReader(defaultGIFBytes))
		if err != nil {
			return nil, fmt.Errorf("gif: decode embedded default: %w", err)
		}
		return g, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("gif: open %q: %w", path, err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, fmt.Errorf("gif: decode %q: %w", path, err)
	}
	return g, nil
}

// buildFrames composites GIF frames onto a canvas respecting the disposal method,
// then scales each composited canvas to the display dimensions.
func buildFrames(g *gif.GIF, dims types.BoardDimensions) (types.Animation, error) {
	if len(g.Image) == 0 {
		return nil, fmt.Errorf("gif: file contains no frames")
	}

	canvasW := g.Config.Width
	canvasH := g.Config.Height
	if canvasW == 0 || canvasH == 0 {
		// Fall back to the bounds of the first frame.
		canvasW = g.Image[0].Bounds().Max.X
		canvasH = g.Image[0].Bounds().Max.Y
	}

	bgColor := color.RGBA{}
	if g.Config.ColorModel != nil {
		if p, ok := g.Config.ColorModel.(color.Palette); ok && len(p) > 0 {
			r, gr, b, a := p[0].RGBA()
			bgColor = color.RGBA{R: uint8(r >> 8), G: uint8(gr >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
		}
	}

	canvas := image.NewRGBA(image.Rect(0, 0, canvasW, canvasH))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// prevCanvas holds a copy of the canvas before the current frame was drawn,
	// needed for disposal method 3 ("restore to previous").
	prevCanvas := image.NewRGBA(canvas.Bounds())

	frames := make(types.Animation, 0, len(g.Image))

	for i, srcFrame := range g.Image {
		disposal := byte(0)
		if i < len(g.Disposal) {
			disposal = g.Disposal[i]
		}

		// Snapshot canvas before drawing, for disposal method 3.
		copy(prevCanvas.Pix, canvas.Pix)

		// Composite this frame onto the canvas.
		draw.Draw(canvas, srcFrame.Bounds(), srcFrame, srcFrame.Bounds().Min, draw.Over)

		// Snapshot the composited result for this frame's output image.
		snap := image.NewRGBA(canvas.Bounds())
		copy(snap.Pix, canvas.Pix)

		delay := 10 * time.Millisecond
		if i < len(g.Delay) && g.Delay[i] > 0 {
			// GIF delay is in centiseconds.
			delay = time.Duration(g.Delay[i]) * 10 * time.Millisecond
		}

		// Scale the snapshot to the display dimensions.
		scaled := scaleImage(snap, dims, "fit")
		frames = append(frames, types.AnimationFrame{Img: scaled, Duration: delay})

		// Apply disposal for the next frame.
		switch disposal {
		case 2: // restore to background
			draw.Draw(canvas, srcFrame.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
		case 3: // restore to previous
			copy(canvas.Pix, prevCanvas.Pix)
			// 0 and 1: leave canvas as-is
		}
	}

	return frames, nil
}
