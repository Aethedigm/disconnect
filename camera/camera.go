package camera

import (
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

var Cam Camera

type Camera struct {
	Position utils.Vector2
}

func GetCamera() *Camera {
	return &Cam
}

func (c *Camera) Move(res utils.Vector2) {
	c.Position = res

	w, h := ebiten.WindowSize()
	c.Position.Sub(utils.Vector2{
		X: float64(w / 2),
		Y: float64(h / 2),
	})
}
