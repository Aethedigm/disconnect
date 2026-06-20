package physics

import (
	"main/utils"
)

type CircleCollider struct {
	Center utils.Vector2
	Radius float64
}

func CircleCollidersCollided(c1, c2 CircleCollider) bool {
	dx := c2.Center.X - c1.Center.X
	dy := c2.Center.Y - c1.Center.Y
	r := c1.Radius + c2.Radius

	return dx*dx+dy*dy < r*r
}

func ResolveCircleOverlap(c1, c2 CircleCollider) (pushDir utils.Vector2) {
	delta := c1.Center.Subbed(c2.Center)
	dist := delta.Length()
	minDist := c1.Radius + c2.Radius

	if dist == 0 || dist >= minDist {
		return utils.Vector2Zero()
	}

	overlap := minDist - dist
	pushDir = delta.Normalized()

	pushDir.MulScalar(overlap)
	return
}
