package objects

import (
	"log"
	"main/camera"
	"main/utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Input struct {
	Move           utils.Vector2
	AimTargetAngle float64
	HasAim         bool
	Fire           bool
}

type WorldContext struct {
	SelfPosition utils.Vector2
	SelfLowerRot float64
	SelfUpperRot float64
	SelfTeam     Team

	NearbyMecha  []ObjectInfo
	NearbyTowers []ObjectInfo
}

type ObjectInfo struct {
	Position utils.Vector2
	Team     Team
	Distance float64
}

type Controller interface {
	Update(utils.Vector2, WorldContext) Input
}

type AIController struct{}
type PlayerController struct{}

func (a *AIController) Update(position utils.Vector2, wc WorldContext) (inp Input) {
	// Should fire?
	// If aiming at enemy, and in range

	// Where to Aim?
	// Are there enemies around?

	// Should take tower?

	log.Print(wc)

	return
}

func (p *PlayerController) Update(position utils.Vector2, _ WorldContext) (inp Input) {
	inp.Move = utils.Vector2{}
	inp.Fire = false

	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) {
		inp.Fire = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		inp.Move.X--
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		inp.Move.X++
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		inp.Move.Y++
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		inp.Move.Y--
	}

	// Set turn target angle
	mPosX, mPosY := ebiten.CursorPosition()
	cam := camera.GetCamera()
	mousePos := utils.Vector2{
		X: float64(mPosX) + cam.Position.X,
		Y: float64(mPosY) + cam.Position.Y,
	}
	aim := mousePos.Subbed(position)

	if !aim.Equals(utils.Vector2Zero()) {
		inp.AimTargetAngle = math.Atan2(aim.Y, aim.X)
		inp.HasAim = true
	}

	return
}
