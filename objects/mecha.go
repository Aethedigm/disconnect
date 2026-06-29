package objects

import (
	"image/color"
	"main/camera"
	"main/physics"
	"main/utils"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SmokePuff struct {
	Position utils.Vector2
	Radius   float64
	Life     int
	MaxLife  int
}

type Mecha struct {
	Position utils.Vector2

	Controller Controller
	Team       Team

	LowerPart MechaLowerPart
	UpperPart MechaUpperPart

	Health float64

	wc WorldContext

	targetSprite *ebiten.Image

	smokePuffs []SmokePuff
	smokeTimer int
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

	if m.Health < 40 {
		m.UpdateDamageSmoke()
	}

	return
}

func (m *Mecha) UpdateDamageSmoke() {
	for i := 0; i < len(m.smokePuffs); i++ {
		m.smokePuffs[i].Radius += 0.2
		m.smokePuffs[i].Life--

		if m.smokePuffs[i].Life <= 0 {
			m.smokePuffs = append(m.smokePuffs[:i], m.smokePuffs[i+1:]...)
			i--
		}
	}

	m.smokeTimer--
	if m.smokeTimer > 0 || len(m.smokePuffs) > 5 {
		return
	}

	m.smokeTimer = 10

	m.smokePuffs = append(m.smokePuffs, SmokePuff{
		Position: m.Position,
		Radius:   3 + rand.Float64()*2,
		Life:     60,
		MaxLife:  60,
	})
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

func drawSmoke(screen *ebiten.Image, pos utils.Vector2) {
	cam := camera.GetCamera()

	x := float32(pos.X - cam.Position.X)
	y := float32(pos.Y - cam.Position.Y)

	vector.FillCircle(screen, x, y, 5, color.RGBA{50, 50, 50, 160}, true)
	vector.FillCircle(screen, x+5, y-3, 4, color.RGBA{80, 80, 80, 120}, true)
	vector.FillCircle(screen, x-4, y-5, 3, color.RGBA{30, 30, 30, 140}, true)
}

func (m *Mecha) DrawSmokePuffs(screen *ebiten.Image) {
	cam := camera.GetCamera()
	for _, smoke := range m.smokePuffs {
		alpha := uint8(180 * float64(smoke.Life) / float64(smoke.MaxLife))
		smokeColor := color.RGBA{80, 80, 80, alpha}
		vector.FillCircle(screen, float32(smoke.Position.X-cam.Position.X), float32(smoke.Position.Y-cam.Position.Y), float32(smoke.Radius), smokeColor, true)
	}
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

	if m.Health < 40 {
		// drawSmoke(screen, m.Position)
		m.DrawSmokePuffs(screen)
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
