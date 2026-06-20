package objects

import "main/physics"

type Collisions interface {
	Collider() physics.CircleCollider
}
