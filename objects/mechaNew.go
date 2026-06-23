package objects

import (
	"main/utils"
)

func NewEnemyMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position: position,
		Controller: &AIController{
			Team: TeamEnemy,
		},
		Team:      TeamEnemy,
		Health:    100,
		LowerPart: NewMechaBottomTwo(),
		UpperPart: NewMechaTopThree(),
	}
}

func NewFriendlyMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position: position,
		Controller: &AIController{
			Team: TeamFriendly,
		},
		Team:      TeamFriendly,
		Health:    100,
		LowerPart: NewMechaBottomTwo(),
		UpperPart: NewMechaTopOne(),
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
