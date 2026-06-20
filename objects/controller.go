package objects

import (
	"log"
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

type Controller interface {
	Update(utils.Vector2) Input
}

type AIEnemyController struct{}
type AIFriendlyController struct{}
type PlayerController struct{}

func (a *AIEnemyController) Update(position utils.Vector2) (inp Input) {
	return
}

func (a *AIFriendlyController) Update(position utils.Vector2) (inp Input) {
	return
}

func (p *PlayerController) Update(position utils.Vector2) (inp Input) {
	inp.Move = utils.Vector2{}
	inp.Fire = false

	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) {
		inp.Fire = true
		log.Println("Should fire...")
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
	mousePos := utils.Vector2{X: float64(mPosX), Y: float64(mPosY)}
	aim := mousePos.Subbed(position)

	if !aim.Equals(utils.Vector2Zero()) {
		inp.AimTargetAngle = math.Atan2(aim.Y, aim.X)
		inp.HasAim = true
	}

	return
}
