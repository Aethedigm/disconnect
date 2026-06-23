package objects

import (
	"main/camera"
	"main/utils"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Input struct {
	Move           utils.Vector2
	AimTargetAngle float64
	HasAim         bool
	Fire           bool
}

type MechaContext struct {
	Position      utils.Vector2
	UpperRot      float64
	UpperRotSpeed float64
	LowerRot      float64
	LowerRotSpeed float64
	GunRange      float64
	Health        float64
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
	Progress float64
}

type ClosestEnemy struct {
	Found bool
	Enemy ObjectInfo
}

type Controller interface {
	Update(MechaContext, WorldContext) Input
}

type AIState int

const (
	CaptureTower AIState = iota
	FightEnemy
	Wander
)

type AIController struct {
	State              AIState
	Team               Team
	TargetTravelVector utils.Vector2
}
type PlayerController struct{}

func GetClosestEnemy(wc WorldContext, a *AIController) (enemy ClosestEnemy) {
	for _, mecha := range wc.NearbyMecha {
		if mecha.Team != a.Team {
			if !enemy.Found {
				enemy.Enemy = mecha
				enemy.Found = true
			} else if mecha.Distance < enemy.Enemy.Distance {
				enemy.Enemy = mecha
			}
		}
	}

	return
}

func GetClosestTower(wc WorldContext, a *AIController) (closestTower *ObjectInfo) {
	// Should take tower?
	// NOTE: There will never be 0 towers. This game is all about Tower capture
	for i := range wc.NearbyTowers {
		if wc.NearbyTowers[i].Team == a.Team && wc.NearbyTowers[i].Progress > 500 {
			continue
		}

		if closestTower == nil || wc.NearbyTowers[i].Distance < closestTower.Distance {
			closestTower = &wc.NearbyTowers[i]
		}
	}

	return
}

// TODO: Replace AIController.Update()
func (a *AIController) StateMachine(mc MechaContext, wc WorldContext) (inp Input) {
	enemy := GetClosestEnemy(wc, a)
	_ = GetClosestTower(wc, a)

	switch {
	case enemy.Found:
		a.State = FightEnemy
	default:
		a.State = Wander
	}

	return
}

func (a *AIController) Update(mc MechaContext, wc WorldContext) (inp Input) {
	enemy := GetClosestEnemy(wc, a)
	closestTower := GetClosestTower(wc, a)

	if enemy.Found {
		aim := enemy.Enemy.Position.Subbed(mc.Position)
		if !aim.Equals(utils.Vector2Zero()) {
			inp.HasAim = true
			inp.AimTargetAngle = math.Atan2(aim.Y, aim.X)
		}

		// Should fire?
		// If aiming approx at enemy, and in range
		if enemy.Enemy.Distance < mc.GunRange/3*2 {
			// Enemy in range, engage
			a.State = FightEnemy
			// We will want to circle the enemy at ~ 2/3rds our gun range typically
			if math.Abs(inp.AimTargetAngle-mc.UpperRot) < 0.2 {
				inp.Fire = true
			}
		}

		if enemy.Enemy.Distance < mc.GunRange && mc.Health < 40 {
			// Disengage
			pos := enemy.Enemy.Position.Subbed(mc.Position)
			turnDir := math.Atan2(pos.Y, pos.X)

			if turnDir < -0.02 {
				inp.Move.X--
			} else if turnDir > 0.02 {
				inp.Move.X++
			}

			inp.Move.Y++
		}

	} else {
		a.State = CaptureTower
	}

	if closestTower != nil {
		if closestTower.Progress > 500 && closestTower.Team == a.Team {
			a.State = Wander
		}
	} else {
		a.State = Wander
	}

	if a.State == CaptureTower {
		turnDir := utils.RotateTowardsVectorFromVector(closestTower.Position, mc.Position, mc.LowerRot, mc.LowerRotSpeed)

		if closestTower.Distance > 75 {
			if turnDir < -0.02 {
				inp.Move.X--
			} else if turnDir > 0.02 {
				inp.Move.X++
			}

			// Drive forward??
			if math.Abs(turnDir) < .1 {
				inp.Move.Y++
			}
		}
	}

	if a.State == Wander {
		inp.Move.Y++
		RandomMovement(&inp)
	}

	return
}

func RandomMovement(inp *Input) {
	val := rand.Float64()
	if val < 0.4 {
		inp.Move.X--
	} else if val > 0.6 {
		inp.Move.X++
	}
}

func (p *PlayerController) Update(mc MechaContext, _ WorldContext) (inp Input) {
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
	aim := mousePos.Subbed(mc.Position)

	if !aim.Equals(utils.Vector2Zero()) {
		inp.AimTargetAngle = math.Atan2(aim.Y, aim.X)
		inp.HasAim = true
	}

	return
}
