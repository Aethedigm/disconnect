package objects

import (
	"main/camera"
	"main/data"
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
}

type MechaLowerPart struct {
	Sprite        *ebiten.Image
	Rotation      float64
	RotationSpeed float64
	DriveSpeed    float64 // Forward/Backward speed
}

type MechaUpperPart struct {
	Sprite        *ebiten.Image
	Rotation      float64
	RotationSpeed float64
	Guns          []GunMount
}

func NewEnemyMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position:   position,
		Controller: &AIEnemyController{},
		Team:       TeamEnemy,
		Health:     100,
		LowerPart: MechaLowerPart{
			Sprite:        utils.ImageDecode(data.TankBottomOne),
			DriveSpeed:    1,
			RotationSpeed: 2 * math.Pi / 180,
		},
		UpperPart: MechaUpperPart{
			Sprite:        utils.ImageDecode(data.TankTopOne),
			RotationSpeed: 2 * math.Pi / 180,
		},
	}
}

func NewFriendlyMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position:   position,
		Controller: &AIFriendlyController{},
		Team:       TeamFriendly,
		Health:     100,
		LowerPart: MechaLowerPart{
			Sprite:        utils.ImageDecode(data.TankBottomOne),
			DriveSpeed:    1,
			RotationSpeed: 2 * math.Pi / 180,
		},
		UpperPart: MechaUpperPart{
			Sprite:        utils.ImageDecode(data.TankTopOne),
			RotationSpeed: 2 * math.Pi / 180,
		},
	}
}

func NewPlayerMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position:   position,
		Controller: &PlayerController{},
		Team:       TeamFriendly,
		Health:     100,
		LowerPart: MechaLowerPart{
			Sprite:        utils.ImageDecode(data.MechaBottomLegs),
			DriveSpeed:    1,
			RotationSpeed: 2 * math.Pi / 180,
		},
		UpperPart: MechaUpperPart{
			Sprite:        utils.ImageDecode(data.MechaTop),
			RotationSpeed: 2 * math.Pi / 180,
			Guns: []GunMount{
				{
					LocalPosition: utils.Vector2{X: 32, Y: 7},
					LocalRotation: 0,
					Weapon: &Weapon{
						FireRate: 10,
					},
				},
				{
					LocalPosition: utils.Vector2{X: 32, Y: -7},
					LocalRotation: 0,
					Weapon: &Weapon{
						FireRate: 10,
					},
				},
			},
		},
	}
}

func (m *Mecha) Update() (res UpdateResult) {
	inp := m.Controller.Update(m.Position)

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
