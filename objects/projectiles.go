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

type Explosion struct {
	Position     utils.Vector2
	SplashDamage float64
	SplashRadius float64
	Team         Team
}

func (e *Explosion) Update() (res UpdateResult) {
	return
}

func (e *Explosion) Draw(_ *ebiten.Image) {}

func (e *Explosion) TeamOwned() Team {
	return e.Team
}

func (e *Explosion) Explosion() physics.CircleCollider {
	return physics.CircleCollider{
		Center: e.Position,
		Radius: e.SplashRadius,
	}
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

type Rocket struct {
	Sprite         *ebiten.Image
	Damage         float64
	SplashDamage   float64
	SplashRadius   float64
	Range          float64
	Speed          float64
	Angle          float64
	Position       utils.Vector2
	Team           Team
	bShouldDestroy bool
}

func NewRocket(damage, rangeVal, speed, angle float64, position utils.Vector2, team Team) *Rocket {
	return &Rocket{
		Sprite:       utils.ImageDecode(data.Rocket),
		Damage:       damage,
		SplashDamage: damage / 2,
		SplashRadius: damage * 2,
		Range:        rangeVal,
		Speed:        speed,
		Angle:        angle,
		Position:     position,
		Team:         team,
	}
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

func (r *Rocket) Explode() (res UpdateResult) {
	res.Spawn = append(res.Spawn, &Explosion{
		Position:     r.Position,
		SplashDamage: r.SplashDamage,
		SplashRadius: r.SplashRadius,
		Team:         r.Team,
	})

	return
}

func (r *Rocket) Update() (res UpdateResult) {
	res.Destroy = r.bShouldDestroy

	targetDir := utils.Vector2{
		X: math.Cos(r.Angle),
		Y: math.Sin(r.Angle),
	}

	targetDir.MulScalar(r.Speed)
	r.Range -= r.Speed

	r.Position.Add(targetDir)

	if r.Range < 0 {
		r.bShouldDestroy = true
		// Spawn explosion
		res.Spawn = append(res.Spawn, &Explosion{
			Position:     r.Position,
			SplashDamage: r.SplashDamage,
			SplashRadius: r.SplashRadius,
			Team:         r.Team,
		})
	}

	return
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

func (r *Rocket) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	bounds := r.Sprite.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	op.GeoM.Translate(float64(-w/2), float64(-h/2))
	op.GeoM.Rotate(r.Angle)

	cam := camera.GetCamera()
	op.GeoM.Translate(r.Position.X-cam.Position.X, r.Position.Y-cam.Position.Y)

	screen.DrawImage(r.Sprite, op)
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

func (r *Rocket) TeamOwned() Team {
	return r.Team
}

func (b *Bullet) TeamOwned() Team {
	return b.Team
}

func (r *Rocket) Collider() physics.CircleCollider {
	return physics.CircleCollider{
		Center: r.Position,
		Radius: float64(r.Sprite.Bounds().Dx()) / 2,
	}
}

func (b *Bullet) Collider() physics.CircleCollider {
	return physics.CircleCollider{
		Center: b.Position,
		Radius: float64(b.Sprite.Bounds().Dx()) / 2,
	}
}

func (r *Rocket) ProjectileDamage() float64 {
	return r.Damage
}

func (b *Bullet) ProjectileDamage() float64 {
	return b.Damage
}

func (r *Rocket) Destroy() {
	r.bShouldDestroy = true
}

func (b *Bullet) Destroy() {
	b.bShouldDestroy = true
}

func (r *Rocket) IsDestroyed() bool {
	return r.bShouldDestroy
}

func (b *Bullet) IsDestroyed() bool {
	return b.bShouldDestroy
}
