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
