package objects

import (
	"main/data"
	"main/utils"
)

func NewEnemyMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position: position,
		Controller: &AIController{
			Team: TeamEnemy,
		},
		Team:         TeamEnemy,
		Health:       100,
		LowerPart:    NewMechaBottomOne(),
		UpperPart:    NewMechaTopOne(),
		targetSprite: utils.ImageDecode(data.SelectorOutline),
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
		LowerPart: NewMechaBottomOne(),
		UpperPart: NewMechaTopOne(),
	}
}

func NewPlayerMecha(position utils.Vector2, upper, lower int) *Mecha {

	var up MechaUpperPart
	var low MechaLowerPart

	switch upper {
	case 1:
		up = NewMechaTopTwo()
	case 2:
		up = NewMechaTopThree()
	case 3:
		up = NewMechaTopFour()
	default:
		up = NewMechaTopOne()
	}

	switch lower {
	case 1:
		low = NewMechaBottomTwo()
	default:
		low = NewMechaBottomOne()
	}

	return &Mecha{
		Position:   position,
		Controller: &PlayerController{},
		Team:       TeamFriendly,
		Health:     100,
		UpperPart:  up,
		LowerPart:  low,
	}
}
