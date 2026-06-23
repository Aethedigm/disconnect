package objects

import (
	"log"
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

type Controller interface {
	Update(MechaContext, WorldContext) Input
}

type AIState int

const (
	CaptureTower AIState = iota
	FightEnemy
	FleeEnemies
	Wander
)

type AIController struct {
	State              AIState
	Team               Team
	TargetTravelVector utils.Vector2
	TargetTravelTime   int
	TargetEnemy        *ObjectInfo
}
type PlayerController struct{}

func GetMechaBreakdowns(wc WorldContext, a *AIController) (enemies, friendlies []*ObjectInfo, enemyFound bool) {
	for _, mecha := range wc.NearbyMecha {
		if mecha.Team == a.Team {
			friendlies = append(friendlies, &mecha)
		} else {
			enemies = append(enemies, &mecha)
			enemyFound = true
		}
	}

	return
}

func GetBestTower(wc WorldContext, mc MechaContext, a *AIController) (bestTower *ObjectInfo, towerFound bool) {
	var closestTower *ObjectInfo

	for _, tower := range wc.NearbyTowers {
		if tower.Position.DistanceTo(mc.Position) < 80 && (tower.Team != a.Team || (tower.Team == a.Team && tower.Progress < 750)) {
			// We're here
			// Our Team Owns it and we need more progress
			// Our team doesn't own it
			bestTower = &tower
			towerFound = true

			return
		}
	}

	for _, tower := range wc.NearbyTowers {
		if tower.Team == a.Team {
			if tower.Progress > 750 {
				continue
			} else {
				bestTower = &tower
				towerFound = true
				continue
			}
		} else {

		}

		if bestTower == nil {
			if closestTower == nil {
				closestTower = &tower
			}

			if tower.Distance < closestTower.Distance {
				closestTower = &tower
			}
		}
	}

	if bestTower == nil && closestTower != nil {
		bestTower = closestTower
		towerFound = true
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
func (a *AIController) StateMachine(mc MechaContext, wc WorldContext) (enemies, friendlies []*ObjectInfo, bestTower *ObjectInfo) {
	enemyFound, towerFound := false, false

	enemies, friendlies, enemyFound = GetMechaBreakdowns(wc, a)
	bestTower, towerFound = GetBestTower(wc, mc, a)

	switch {
	// case enemyFound && len(enemies) > len(friendlies)+1:
	case enemyFound && mc.Health < 25 && len(friendlies) == 0:
		a.State = FleeEnemies
	case enemyFound:
		a.State = FightEnemy
	case towerFound:
		a.State = CaptureTower
	default:
		a.State = Wander
	}

	return
}

func AimAtEnemy(pos1, pos2 utils.Vector2, inp *Input) {
	aim := pos1.Subbed(pos2)
	if !aim.Equals(utils.Vector2Zero()) {
		inp.HasAim = true
		inp.AimTargetAngle = math.Atan2(aim.Y, aim.X)
	}
}

func Contains(o []*ObjectInfo, t *ObjectInfo) bool {
	for _, obj := range o {
		if obj == t {
			return true
		}
	}

	return false
}

func (a *AIController) Update(mc MechaContext, wc WorldContext) (inp Input) {
	a.TargetTravelTime--
	enemies, _, bestTower := a.StateMachine(mc, wc)

	if a.State == FleeEnemies {
		log.Println("Want to flee enemies")
		// run generally away?
	}

	if a.State == FightEnemy {
		log.Println("Want to fight enemies")
		if !Contains(enemies, a.TargetEnemy) {
			// Remove target if not in context anymore
			a.TargetEnemy = nil
		}

		if a.TargetEnemy == nil {
			// Select target randomly from context
			a.TargetEnemy = enemies[rand.Intn(len(enemies))]
		}

		AimAtEnemy(a.TargetEnemy.Position, mc.Position, &inp)

		if a.TargetEnemy.Distance < mc.GunRange/3*2 {
			if math.Abs(inp.AimTargetAngle-mc.UpperRot) < 0.2 {
				inp.Fire = true
			}
		}

		// Set Target Vector for Travel
		// Set arbitrary point within range of enemy
		offset := utils.Vector2{
			X: rand.Float64()*500 - 250,
			Y: rand.Float64()*500 - 250,
		}

		if mc.Position.DistanceTo(a.TargetTravelVector) < 35 || a.TargetTravelTime <= 0 {
			a.TargetTravelVector = a.TargetEnemy.Position.Added(offset)
			a.TargetTravelTime = 30
		}
	}

	if a.State == CaptureTower {
		log.Println("Want to capture tower")
		if a.TargetTravelTime <= 0 {
			a.TargetTravelVector = bestTower.Position
			a.TargetTravelTime = 120
		}
	}

	if a.State == Wander {
		log.Println("Want to wander")
	}

	// Move towards target travel vector
	turnDir := utils.RotateTowardsVectorFromVector(a.TargetTravelVector, mc.Position, mc.LowerRot, mc.LowerRotSpeed)
	targetTravelDistance := mc.Position.DistanceTo(a.TargetTravelVector)

	if targetTravelDistance > 10 {
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
