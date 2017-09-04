package main

import (
	"math"
)

type Vec struct {
	X float64
	Y float64
	Z float64
}

func (v1 *Vec) Add(v2 *Vec) *Vec {
	return &Vec{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

func (v1 *Vec) Div(v2 *Vec) *Vec {
	return &Vec{v1.X / v2.X, v1.Y / v2.Y, v1.Z / v2.Z}
}

func (v1 *Vec) Sub(v2 *Vec) *Vec {
	return &Vec{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

func (v1 *Vec) Scale(s float64) *Vec {
	return &Vec{v1.X * s, v1.Y * s, v1.Z * s}
}

func (v1 *Vec) Lerp(v2 *Vec, s float64) *Vec {
	return v1.Scale(1 - s).Add(v2.Scale(s))
}

func (v1 *Vec) Dot(v2 *Vec) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 *Vec) Len() float64 {
	return math.Sqrt(math.Pow(v1.X, 2) + math.Pow(v1.Y, 2) + math.Pow(v1.Z, 2))
}

func (v1 *Vec) Normalize() *Vec {
	mag := math.Sqrt(v1.Dot(v1))

	return &Vec{
		v1.X / mag,
		v1.Y / mag,
		v1.Z / mag,
	}
}
