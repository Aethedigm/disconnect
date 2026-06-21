package objects

import (
	"main/camera"
	"main/data"
	"main/physics"
	"main/utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Projectile interface {
	GameObject
	HasTeam
	Collider() physics.CircleCollider
	ProjectileDamage() float64
	Destroy()
	IsDestroyed() bool
}

type Bullet struct {
	Sprite         *ebiten.Image
	Damage         float64
	Range          float64
	Speed          float64
	Angle          float64
	Position       utils.Vector2
	Team           Team
	bShouldDestroy bool
}

func NewBullet(damage, rangeVal, speed, angle float64, position utils.Vector2, team Team) *Bullet {
	return &Bullet{
		Sprite:   utils.ImageDecode(data.Bullet),
		Damage:   damage,
		Range:    rangeVal,
		Speed:    speed,
		Angle:    angle,
		Position: position,
		Team:     team,
	}
}

type Rocket struct {
	Sprite       *ebiten.Image
	ImpactDamage float64
	SplashDamage float64
	SplashRadius float64
	Range        float64
	Speed        float64
	Angle        float64
	Position     utils.Vector2
	Team         Team
}

func (b *Bullet) Update() (res UpdateResult) {
	res.Destroy = b.bShouldDestroy

	// Move position by angle
	targetDir := utils.Vector2{
		X: math.Cos(b.Angle),
		Y: math.Sin(b.Angle),
	}

	targetDir.MulScalar(b.Speed)
	b.Range -= b.Speed

	b.Position.Add(targetDir)

	if b.Range < 0 {
		b.bShouldDestroy = true
	}
	return
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	bounds := b.Sprite.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	op.GeoM.Translate(float64(-w/2), float64(-h/2))
	op.GeoM.Rotate(b.Angle)

	cam := camera.GetCamera()
	op.GeoM.Translate(b.Position.X-cam.Position.X, b.Position.Y-cam.Position.Y)

	screen.DrawImage(b.Sprite, op)
}

func (b *Bullet) TeamOwned() Team {
	return b.Team
}

func (b *Bullet) Collider() physics.CircleCollider {
	return physics.CircleCollider{
		Center: b.Position,
		Radius: float64(b.Sprite.Bounds().Dx()) / 2,
	}
}

func (b *Bullet) ProjectileDamage() float64 {
	return b.Damage
}

func (b *Bullet) Destroy() {
	b.bShouldDestroy = true
}

func (b *Bullet) IsDestroyed() bool {
	return b.bShouldDestroy
}
