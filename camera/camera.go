package camera

import (
	"main/utils"
)

var Cam Camera

type Camera struct {
	Position utils.Vector2
	Bounds   utils.Vector2
}

func GetCamera() *Camera {
	return &Cam
}

func SetCameraBounds(bounds utils.Vector2) {
	cam := GetCamera()
	cam.Bounds = bounds
}

func (c *Camera) Move(res utils.Vector2) {
	c.Position = res

	c.Position.Sub(utils.Vector2{
		X: c.Bounds.X / 2,
		Y: c.Bounds.Y / 2,
	})
}
