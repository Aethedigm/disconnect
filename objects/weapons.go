package objects

import (
	"main/utils"
)

type Weapon struct {
	FireRate     float64
	FireCooldown float64
}

type GunMount struct {
	LocalPosition utils.Vector2
	LocalRotation float64
	Weapon        *Weapon
	AmmoCount     int
}

func (g *GunMount) Update() {
	if g.Weapon.FireCooldown > 0 {
		g.Weapon.FireCooldown--
	}
}

func (g *GunMount) Fire(mechaPos utils.Vector2, upperRot float64, team Team) []GameObject {
	offset := g.LocalPosition.Rotated(upperRot)
	offset.Add(mechaPos)

	return g.Weapon.Fire(offset, upperRot+g.LocalRotation, team)
}

func (w *Weapon) Fire(pos utils.Vector2, rot float64, team Team) []GameObject {
	var projectiles []GameObject

	if w.FireCooldown <= 0 {
		w.FireCooldown = w.FireRate
		projectiles = append(projectiles, NewBullet(
			1, 500, 10, rot, pos, team,
		))
	}

	return projectiles
}
