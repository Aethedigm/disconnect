package objects

import (
	"log"
	"main/utils"
)

type Weapon struct{}

type GunMount struct {
	LocalPosition utils.Vector2
	LocalRotation float64
	Weapon        *Weapon
}

func (g *GunMount) Fire(mechaPos utils.Vector2, upperRot float64, team Team) []GameObject {
	offset := g.LocalPosition.Rotated(upperRot)
	offset.Add(mechaPos)

	return g.Weapon.Fire(offset, upperRot+g.LocalRotation, team)
}

func (w *Weapon) Fire(pos utils.Vector2, rot float64, team Team) []GameObject {
	var projectiles []GameObject

	log.Println("Bang!")

	return projectiles
}
