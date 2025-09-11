package tests

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// GenerateTestImage создает простое тестовое изображение
func GenerateTestImage() ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Заполняем фон
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			r := uint8(x * 255 / 100)
			g := uint8(y * 255 / 100)
			b := uint8((x + y) * 255 / 200)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
