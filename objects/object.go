package objects

import (
	"main/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameObject interface {
	Update()
	Draw(*ebiten.Image)
}

type TeamOwned interface {
	TeamOwnership() Team
}

type RadioNode interface {
	TeamOwned() Team
	RadioCollider() physics.CircleCollider
}

type Capturer interface {
	TeamOwned() Team
	Collider() physics.CircleCollider
}
