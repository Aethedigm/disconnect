package objects

import (
	"main/data"
	"main/utils"
	"math"
)

func NewEnemyMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position:   position,
		Controller: &AIController{},
		Team:       TeamEnemy,
		Health:     100,
		LowerPart: MechaLowerPart{
			Sprite:        utils.ImageDecode(data.TankBottomOne),
			DriveSpeed:    1,
			RotationSpeed: 2 * math.Pi / 180,
		},
		UpperPart: MechaUpperPart{
			Sprite:        utils.ImageDecode(data.TankTopOne),
			RotationSpeed: 2 * math.Pi / 180,
		},
	}
}

func NewFriendlyMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position:   position,
		Controller: &AIController{},
		Team:       TeamFriendly,
		Health:     100,
		LowerPart: MechaLowerPart{
			Sprite:        utils.ImageDecode(data.TankBottomOne),
			DriveSpeed:    1,
			RotationSpeed: 2 * math.Pi / 180,
		},
		UpperPart: MechaUpperPart{
			Sprite:        utils.ImageDecode(data.TankTopOne),
			RotationSpeed: 2 * math.Pi / 180,
		},
	}
}

func NewPlayerMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position:   position,
		Controller: &PlayerController{},
		Team:       TeamFriendly,
		Health:     100,
		LowerPart:  NewMechaBottomOne(),
		// LowerPart: NewMechaBottomTwo(),
		UpperPart: NewMechaTopOne(),
		// UpperPart: NewMechaTopThree(),
	}
}
