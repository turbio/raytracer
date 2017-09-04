package main

import (
	"math"
)

type Sphere struct {
	Center   *Vec
	Radius   float64
	Material *Material
}

func (s *Sphere) Intersects(r *Ray) (float64, bool) {
	rayToSphere := r.Point.Sub(s.Center)

	a := math.Pow(r.Direction.Len(), 2)
	b := r.Direction.Dot(rayToSphere) * 2
	c := math.Pow(rayToSphere.Len(), 2) - math.Pow(s.Radius, 2)

	disc := math.Pow(b, 2) - (4 * a * c)

	if disc < 0 {
		return 0, false
	}

	t1 := (-b + math.Sqrt(math.Pow(b, 2)-4*a*c)) / (2 * a)
	t2 := (-b - math.Sqrt(math.Pow(b, 2)-4*a*c)) / (2 * a)

	t := math.Min(t1, t2)

	if t < 0 {
		return 0, false
	}

	return t, true
}
