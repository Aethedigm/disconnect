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
	TeamOwned() Team
	RadioCollider() physics.CircleCollider
}

type Capturer interface {
	TeamOwned() Team
	Collider() physics.CircleCollider
}
