package boards

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	xdraw "golang.org/x/image/draw"

	"github.com/benwiebe/udb-plugin-library/types"
)

//go:embed assets/sample_image.png
var defaultSpriteBytes []byte

// ImageBoard renders a static image scaled to the display dimensions.
// It demonstrates BoardTypeStatic with an external (or embedded) asset: the source
// image is loaded once in Init() and scaled to fit, then returned unchanged on every
// Render() call.
//
// Optional JSON config fields:
//   - path: filesystem path to a PNG or JPEG file; uses the embedded sample if omitted
//   - fit:  "fit" (default) — scale to fit inside the display, centred with black bars
//     "fill" — scale to fill the display, cropping edges to preserve aspect ratio
//     "stretch" — ignore aspect ratio and stretch to the exact display size
type ImageBoard struct {
	cachedImage image.Image
}

func NewSpriteBoard() *ImageBoard {
	return &ImageBoard{}
}

func (b *ImageBoard) GetId() string   { return "sprite" }
func (b *ImageBoard) GetName() string { return "Image" }
func (b *ImageBoard) GetSupportedDimensions() []types.BoardDimensions {
	return []types.BoardDimensions{}
}
func (b *ImageBoard) GetType() types.BoardType  { return types.BoardTypeStatic }
func (b *ImageBoard) GetDatasourceType() string { return "" }

func (b *ImageBoard) Init(cfg json.RawMessage, _ types.Datasource, dimensions types.BoardDimensions) error {
	path := ""
	fit := "fit"

	if len(cfg) > 0 {
		var parsed struct {
			Path string `json:"path"`
			Fit  string `json:"fit"`
		}
		if err := json.Unmarshal(cfg, &parsed); err != nil {
			return err
		}
		path = parsed.Path
		if parsed.Fit != "" {
			fit = parsed.Fit
		}
	}

	src, err := loadImage(path)
	if err != nil {
		return err
	}

	b.cachedImage = scaleImage(src, dimensions, fit)
	return nil
}

func (b *ImageBoard) Render() image.Image {
	return b.cachedImage
}

// loadImage decodes the image at path, or falls back to the embedded sample when path is empty.
func loadImage(path string) (image.Image, error) {
	if path == "" {
		img, _, err := image.Decode(bytes.NewReader(defaultSpriteBytes))
		if err != nil {
			return nil, fmt.Errorf("image: decode embedded default: %w", err)
		}
		return img, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("image: open %q: %w", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("image: decode %q: %w", path, err)
	}
	return img, nil
}

// scaleImage scales src to fit within dims according to the chosen mode.
func scaleImage(src image.Image, dims types.BoardDimensions, fit string) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, dims.Width, dims.Height))

	srcW := src.Bounds().Dx()
	srcH := src.Bounds().Dy()

	var dstRect image.Rectangle
	switch fit {
	case "fill":
		// Scale so the shorter axis fills the display; crop the other.
		scaleX := float64(dims.Width) / float64(srcW)
		scaleY := float64(dims.Height) / float64(srcH)
		scale := scaleX
		if scaleY > scaleX {
			scale = scaleY
		}
		scaledW := int(float64(srcW) * scale)
		scaledH := int(float64(srcH) * scale)
		ox := (dims.Width - scaledW) / 2
		oy := (dims.Height - scaledH) / 2
		dstRect = image.Rect(ox, oy, ox+scaledW, oy+scaledH)
	case "stretch":
		dstRect = image.Rect(0, 0, dims.Width, dims.Height)
	default: // "fit"
		// Scale so the longer axis fits inside the display; letterbox/pillarbox the other.
		scaleX := float64(dims.Width) / float64(srcW)
		scaleY := float64(dims.Height) / float64(srcH)
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}
		scaledW := int(float64(srcW) * scale)
		scaledH := int(float64(srcH) * scale)
		ox := (dims.Width - scaledW) / 2
		oy := (dims.Height - scaledH) / 2
		dstRect = image.Rect(ox, oy, ox+scaledW, oy+scaledH)
	}

	xdraw.CatmullRom.Scale(dst, dstRect, src, src.Bounds(), xdraw.Over, nil)
	return dst
}
