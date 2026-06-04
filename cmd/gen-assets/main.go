// gen-assets generates the sample image files embedded by the boards package.
// Run via: go run ./cmd/gen-assets/
package main

import (
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math"
	"os"
)

func main() {
	generateSprite("boards/assets/sample_sprite.png")
	generateGIF("boards/assets/sample.gif")
}

// generateSprite writes a 32×32 PNG that looks like a simple pixel-art sprite:
// a coloured diamond on a dark background.
func generateSprite(path string) {
	const size = 32
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	cx, cy := float64(size)/2, float64(size)/2

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Manhattan distance from centre, normalised to [0,1].
			dx := math.Abs(float64(x)+0.5-cx) / cx
			dy := math.Abs(float64(y)+0.5-cy) / cy
			dist := dx + dy

			var c color.RGBA
			switch {
			case dist < 0.35:
				c = color.RGBA{R: 255, G: 220, B: 0, A: 255} // gold centre
			case dist < 0.65:
				c = color.RGBA{R: 255, G: 80, B: 0, A: 255} // orange ring
			case dist < 0.90:
				c = color.RGBA{R: 180, G: 0, B: 200, A: 255} // purple outer
			default:
				c = color.RGBA{R: 10, G: 10, B: 30, A: 255} // near-black corners
			}
			img.SetRGBA(x, y, c)
		}
	}

	f, err := os.Create(path)
	must(err)
	defer f.Close()
	must(png.Encode(f, img))
}

// generateGIF writes an animated GIF: a bright dot orbiting the centre of a dark 32×32 canvas.
func generateGIF(path string) {
	const (
		size   = 32
		frames = 16
		dotR   = 2 // dot radius in pixels
	)

	palette := buildPalette()
	var images []*image.Paletted
	var delays []int

	cx, cy := float64(size)/2-0.5, float64(size)/2-0.5
	orbitR := cx * 0.6

	for i := 0; i < frames; i++ {
		angle := 2 * math.Pi * float64(i) / float64(frames)
		dotX := cx + orbitR*math.Cos(angle)
		dotY := cy + orbitR*math.Sin(angle)

		// Hue rotates with the frame so the dot changes colour.
		hue := float64(i) / float64(frames) * 360
		dotColour := hsvToRGBA(hue)

		frame := image.NewPaletted(image.Rect(0, 0, size, size), palette)

		// Fill background (palette index 0 = dark navy).
		for j := range frame.Pix {
			frame.Pix[j] = 0
		}

		// Draw dot.
		dotIdx := nearestPaletteIndex(palette, dotColour)
		for py := 0; py < size; py++ {
			for px := 0; px < size; px++ {
				dx := float64(px) + 0.5 - dotX
				dy := float64(py) + 0.5 - dotY
				if math.Sqrt(dx*dx+dy*dy) <= float64(dotR) {
					frame.SetColorIndex(px, py, dotIdx)
				}
			}
		}

		images = append(images, frame)
		delays = append(delays, 6) // 6 × 10ms = 60ms per frame ≈ 16 fps
	}

	f, err := os.Create(path)
	must(err)
	defer f.Close()

	must(gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	}))
}

// buildPalette returns a 256-entry palette: index 0 is the background, indices 1–255
// are evenly-spaced hues at full saturation and brightness.
func buildPalette() color.Palette {
	p := make(color.Palette, 256)
	p[0] = color.RGBA{R: 10, G: 10, B: 30, A: 255}
	for i := 1; i < 256; i++ {
		hue := float64(i-1) / 255.0 * 360
		p[i] = hsvToRGBA(hue)
	}
	return p
}

// nearestPaletteIndex finds the palette entry closest (by Euclidean RGB distance) to c.
func nearestPaletteIndex(p color.Palette, c color.RGBA) uint8 {
	best := 0
	bestDist := math.MaxFloat64
	for i, pc := range p {
		pr, pg, pb, _ := pc.RGBA()
		dr := float64(pr>>8) - float64(c.R)
		dg := float64(pg>>8) - float64(c.G)
		db := float64(pb>>8) - float64(c.B)
		d := dr*dr + dg*dg + db*db
		if d < bestDist {
			bestDist = d
			best = i
		}
	}
	return uint8(best)
}

func hsvToRGBA(h float64) color.RGBA {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	c := 1.0
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return color.RGBA{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255), A: 255}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
