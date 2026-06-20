package objects

import (
	"main/data"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Cursor struct {
	Sprite *ebiten.Image
}

func NewCursor() *Cursor {
	return &Cursor{
		Sprite: utils.ImageDecode(data.Crosshair),
	}
}

func (c *Cursor) Update() {}

func (c *Cursor) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	bounds := c.Sprite.Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	mousePosX, mousePosY := ebiten.CursorPosition()

	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Translate(float64(mousePosX), float64(mousePosY))

	screen.DrawImage(c.Sprite, op)
}
