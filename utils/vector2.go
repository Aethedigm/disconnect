package utils

import "math"

type Vector2 struct {
	X, Y float64
}

func (v *Vector2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector2) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v *Vector2) Normalized() Vector2 {
	leng := v.Length()

	if leng == 0 {
		return *v
	}

	return Vector2{
		X: v.X / leng,
		Y: v.Y / leng,
	}
}

// Mutables //

func (v *Vector2) Add(res Vector2) {
	v.X += res.X
	v.Y += res.Y
}

func (v *Vector2) Sub(res Vector2) {
	v.X -= res.X
	v.Y -= res.Y
}

func (v *Vector2) Mul(res Vector2) {
	v.X *= res.X
	v.Y *= res.Y
}

func (v *Vector2) Normalize() {
	n := v.Normalized()

	v.X = n.X
	v.Y = n.Y
}
