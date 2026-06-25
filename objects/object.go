package objects

import (
	"main/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

type UpdateResult struct {
	Spawn   []GameObject
	Destroy bool
}

type GameObject interface {
	Update() UpdateResult
	Draw(*ebiten.Image)
}

type HasTeam interface {
	TeamOwned() Team
}

type RadioNode interface {
	HasTeam
	RadioCollider() physics.CircleCollider
	HasRadio() bool
}

type Capturer interface {
	HasTeam
	Collider() physics.CircleCollider
}

type DamageTarget interface {
	DynamicCollisions
	HasTeam
	ApplyDamage(amount float64)
}
