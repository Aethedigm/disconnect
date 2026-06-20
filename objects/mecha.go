package objects

import (
	"main/data"
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

type MechaPart interface {
	Draw(screen *ebiten.Image)
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
		Controller: &PlayerController{},
		LowerPart: MechaLowerPart{
			Sprite:        utils.ImageDecode(data.TankBottomOne),
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
		// Apply rotation: Rotation Speed towards Angle
		// Target Bearing : inp.AimTargetAngle
		// Current Bearing: m.UpperPart.Rotation
		diff := inp.AimTargetAngle - m.UpperPart.Rotation
		diff = math.Atan2(math.Sin(diff), math.Cos(diff))

		step := min(math.Abs(diff), m.UpperPart.RotationSpeed)
		m.UpperPart.Rotation += math.Copysign(step, diff)
	}

	// TODO: Replace with tank forward/backward only
	tankMove := utils.Vector2{X: 0, Y: inp.Move.Y}
	m.Position.Add(tankMove)

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
