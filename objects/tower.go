package objects

import (
	"main/camera"
	"main/data"
	"main/physics"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	towerCaptureMax = 1000
)

type Tower struct {
	Team Team

	Position utils.Vector2
	Sprite   *ebiten.Image

	CaptureProgress float64
	CapturingTeam   Team
}

func NewNeutralTower(position utils.Vector2) *Tower {
	return &Tower{
		Position: position,
		Sprite:   utils.ImageDecode(data.TowerBase),
	}
}

func (t *Tower) Update() (res UpdateResult) {
	if t.CaptureProgress > towerCaptureMax {
		t.CaptureProgress = towerCaptureMax
	}
	return
}

func (t *Tower) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	bounds := t.Sprite.Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	op.GeoM.Translate(-w/2, -h/2)

	cam := camera.GetCamera()
	op.GeoM.Translate(t.Position.X-cam.Position.X, t.Position.Y-cam.Position.Y)

	// Team Coloring
	if t.Team != TeamNone {
		red := utils.FastBoolConvFloat32(t.Team == TeamEnemy)
		blue := utils.FastBoolConvFloat32(t.Team == TeamFriendly)
		op.ColorScale.Scale(red, 0, blue, 1)
	}

	screen.DrawImage(t.Sprite, op)
}

func (t *Tower) Collider() physics.CircleCollider {
	return physics.CircleCollider{
		Center: t.Position,
		Radius: 28,
	}
}

func (t *Tower) RadioCollider() physics.CircleCollider {
	return physics.CircleCollider{
		Center: t.Position,
		Radius: 500,
	}
}

func (t *Tower) CaptureCollider() physics.CircleCollider {
	return physics.CircleCollider{
		Center: t.Position,
		Radius: 60,
	}
}

func (t *Tower) TeamOwned() Team {
	return t.Team
}

func (t *Tower) HasRadio() bool {
	return true
}
