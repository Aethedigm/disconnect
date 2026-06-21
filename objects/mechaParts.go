package objects

import (
	"main/data"
	"main/utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type MechaUpperPart struct {
	Sprite        *ebiten.Image
	Rotation      float64
	RotationSpeed float64
	Guns          []GunMount
}

// MechaTopOne is a dual Autogunner
func NewMechaTopOne() MechaUpperPart {
	return MechaUpperPart{
		Sprite:        utils.ImageDecode(data.MechaTop),
		RotationSpeed: 2 * math.Pi / 180,
		Guns: []GunMount{
			{
				LocalPosition: utils.Vector2{X: 32, Y: 7},
				Weapon:        NewAutoGun(),
			},
			{
				LocalPosition: utils.Vector2{X: 32, Y: -7},
				Weapon:        NewAutoGun(),
			},
		},
	}
}

// MechaTopTwo is a dual rocketeer
func NewMechaTopTwo() MechaUpperPart {
	return MechaUpperPart{
		Sprite:        utils.ImageDecode(data.MechaTopTwo),
		RotationSpeed: 2 * math.Pi / 180,
		Guns: []GunMount{
			{
				LocalPosition: utils.Vector2{},
				Weapon:        NewRocketLauncher(),
			},
			{
				LocalPosition: utils.Vector2{},
				Weapon:        NewRocketLauncher(),
			},
		},
	}
}

// MechaTopThree is a single Sniper
func NewMechaTopThree() MechaUpperPart {
	return MechaUpperPart{
		Sprite:        utils.ImageDecode(data.TankTopOne),
		RotationSpeed: 2 * math.Pi / 180,
		Guns: []GunMount{
			{
				LocalPosition: utils.Vector2{X: 40, Y: 0},
				Weapon:        NewSniper(),
			},
		},
	}
}

type MechaLowerPart struct {
	Sprite        *ebiten.Image
	Rotation      float64
	RotationSpeed float64
	DriveSpeed    float64 // Forward/Backward speed
}

// Legged Mecha Bottom
func NewMechaBottomOne() MechaLowerPart {
	return MechaLowerPart{
		Sprite:        utils.ImageDecode(data.MechaBottomLegs),
		DriveSpeed:    1,
		RotationSpeed: 2 * math.Pi / 180,
	}
}

// Cat Tread Bottom
func NewMechaBottomTwo() MechaLowerPart {
	return MechaLowerPart{
		Sprite:        utils.ImageDecode(data.TankBottomOne),
		DriveSpeed:    3,
		RotationSpeed: math.Pi / 180,
	}
}
