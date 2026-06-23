package objects

import (
	"main/utils"
)

type Weapon struct {
	FireRate     float64
	FireCooldown float64
	Damage       float64
	Range        float64
	Speed        float64
}

func NewAutoGun() *Weapon {
	return &Weapon{
		FireRate: 10,
		Damage:   1,
		Range:    500,
		Speed:    10,
	}
}

func NewRocketLauncher() *Weapon {
	return &Weapon{
		FireRate: 1000,
		Damage:   30,
		Range:    950,
		Speed:    50,
	}
}

func NewSniper() *Weapon {
	return &Weapon{
		FireRate: 75,
		Damage:   18,
		Range:    5000,
		Speed:    50,
	}
}

func (w *Weapon) Fire(pos utils.Vector2, rot float64, team Team) []GameObject {
	var projectiles []GameObject

	if w.FireCooldown <= 0 {
		w.FireCooldown = w.FireRate
		projectiles = append(projectiles, NewBullet(
			w.Damage, w.Range, w.Speed, rot, pos, team,
		))
	}

	return projectiles
}
