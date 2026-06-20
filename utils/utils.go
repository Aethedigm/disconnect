package utils

import (
	"bytes"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	epsilon float64 = 0.001
)

func ApproxEquals(val1, val2 float64) bool {
	if math.Abs(val1-val2) < epsilon {
		return true
	}

	return false
}

func ImageDecode(spriteData []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(spriteData))
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}
