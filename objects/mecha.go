package objects

import (
	"main/data"
	"main/physics"
	"main/utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Mecha struct {
	Position utils.Vector2

	Controller Controller

	LowerPart MechaLowerPart
	UpperPart MechaUpperPart
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
}

// TODO: Default to AI for this constructor
func NewMecha(position utils.Vector2) *Mecha {
	return &Mecha{
		Position:   position,
		Controller: &AIController{},
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

func (m *Mecha) Update() {
	inp := m.Controller.Update(m.Position)

	// Preserve controller angular input
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

	// Rotate Lower Part
	m.LowerPart.Rotation += m.LowerPart.RotationSpeed * inp.Move.X
}

func drawPart(screen, sprite *ebiten.Image, pos utils.Vector2, rot float64) {
	op := &ebiten.DrawImageOptions{}

	bounds := sprite.Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(rot)
	op.GeoM.Translate(pos.X, pos.Y)

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
