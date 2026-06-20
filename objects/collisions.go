package objects

import (
	"main/physics"
	"main/utils"
)

type Collisions interface {
	Collider() physics.CircleCollider
}

type DynamicCollisions interface {
	Collisions
	Move(delta utils.Vector2)
}
