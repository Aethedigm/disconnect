package data

import (
	_ "embed"
)

var MenuOptions = []string{"Play", "Options", "Quit"}

//go:embed crosshair044.png
var Crosshair []byte

//go:embed TankTop.png
var TankTopOne []byte

//go:embed TankBottom.png
var TankBottomOne []byte

//go:embed TowerBase.png
var TowerBase []byte

//go:embed MechaTop.png
var MechaTop []byte

//go:embed MechaBottomLegs.png
var MechaBottomLegs []byte

//go:embed Bullet.png
var Bullet []byte

//go:embed MechaTop2.png
var MechaTopTwo []byte

//go:embed MechaTopCommander.png
var MechaTopCommander []byte

//go:embed Rocket.png
var Rocket []byte

//go:embed SelectorArrow.png
var SelectorArrow []byte

//go:embed SelectorRed.png
var SelectorOutline []byte

//go:embed shoot.wav
var SFX_Shoot []byte

//go:embed explosion.wav
var SFX_Explode []byte
