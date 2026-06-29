package objects

import (
	"main/camera"
	"main/physics"
	"main/utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Mecha struct {
	Position utils.Vector2

	Controller Controller
	Team       Team

	LowerPart MechaLowerPart
	UpperPart MechaUpperPart

	Health float64

	wc WorldContext

	targetSprite *ebiten.Image
}

func (m *Mecha) Update() (res UpdateResult) {
	mc := MechaContext{
		Position:      m.Position,
		UpperRot:      m.UpperPart.Rotation,
		UpperRotSpeed: m.UpperPart.RotationSpeed,
		LowerRot:      m.LowerPart.Rotation,
		LowerRotSpeed: m.LowerPart.RotationSpeed,
		GunRange:      m.UpperPart.Guns[0].Weapon.Range,
		Health:        m.Health,
	}
	inp := m.Controller.Update(mc, m.wc)

	if inp.Move.Length() > 1 {
		inp.Move = inp.Move.Normalized()
	}

	// Rotate upper part if HasAim
	if inp.HasAim {
		m.UpperPart.Rotation += utils.RotateTowards(inp.AimTargetAngle, m.UpperPart.Rotation, m.UpperPart.RotationSpeed)
	}

	tankMoveDir := utils.Vector2{X: math.Cos(m.LowerPart.Rotation), Y: math.Sin(m.LowerPart.Rotation)}
	tankMoveDir.MulScalar(m.LowerPart.DriveSpeed * inp.Move.Y)
	m.Position.Add(tankMoveDir)

	if inp.Fire {
		// Intent to fire, check conditions if able
		for i := range m.UpperPart.Guns {
			res.Spawn = append(
				res.Spawn,
				m.UpperPart.Guns[i].Fire(m.Position, m.UpperPart.Rotation, m.Team)...,
			)
		}
	}

	// Rotate Lower Part
	m.LowerPart.Rotation += m.LowerPart.RotationSpeed * inp.Move.X

	// Tick time on weapons
	for i := range m.UpperPart.Guns {
		m.UpperPart.Guns[i].Update()
	}

	// Destroyed
	res.Destroy = m.Health < 1

	return
}

func drawPart(screen, sprite *ebiten.Image, pos utils.Vector2, rot float64) {
	op := &ebiten.DrawImageOptions{}

	bounds := sprite.Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(rot)

	cam := camera.GetCamera()
	op.GeoM.Translate(pos.X-cam.Position.X, pos.Y-cam.Position.Y)

	screen.DrawImage(sprite, op)
}

func (m *Mecha) Draw(screen *ebiten.Image) {
	// Draw order matters, Lower then Upper
	drawPart(screen, m.LowerPart.Sprite, m.Position, m.LowerPart.Rotation)
	drawPart(screen, m.UpperPart.Sprite, m.Position, m.UpperPart.Rotation)

	if m.Team == TeamEnemy {
		// Draw red square as a "targeter" for player
		tOp := &ebiten.DrawImageOptions{}
		bounds := m.targetSprite.Bounds()
		w := float64(bounds.Dx())
		h := float64(bounds.Dy())
		tOp.GeoM.Translate(-w/2, -h/2)
		cam := camera.GetCamera()
		tOp.GeoM.Translate(m.Position.X-cam.Position.X, m.Position.Y-cam.Position.Y)
		screen.DrawImage(m.targetSprite, tOp)
	}
}

func (m *Mecha) Collider() physics.CircleCollider {
	bounds := m.LowerPart.Sprite.Bounds()
	return physics.CircleCollider{
		Center: m.Position,
		Radius: max(float64(bounds.Dx())/2, float64(bounds.Dy())/2),
	}
}

func (m *Mecha) Move(delta utils.Vector2) {
	m.Position.Add(delta)
}

func (m *Mecha) TeamOwned() Team {
	return m.Team
}

func (m *Mecha) ApplyDamage(amount float64) {
	m.Health -= amount
}

func (m *Mecha) SetWorldContext(wc WorldContext) {
	m.wc = wc

	m.wc.SelfPosition = m.Position
	m.wc.SelfTeam = m.Team
	m.wc.SelfLowerRot = m.LowerPart.Rotation
	m.wc.SelfUpperRot = m.UpperPart.Rotation
}

func (m *Mecha) HasRadio() bool {
	return m.UpperPart.RadioRange > 0
}

func (m *Mecha) RadioCollider() physics.CircleCollider {
	return physics.CircleCollider{
		Center: m.Position,
		Radius: m.UpperPart.RadioRange,
	}
}
