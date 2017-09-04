package main

import (
	"flag"
	"fmt"
	"math"
	"time"
)

const recurDepth = 4
const samples = 3

var black = &Color{0, 0, 0}

type Scene struct {
	Camera  *Vec
	Spheres []*Sphere
	Lights  []*Light

	X1 *Vec
	X2 *Vec
	X3 *Vec
	X4 *Vec

	Ambient *Color
}

func occluded(scene *Scene, sphere *Sphere, shadowRay *Ray) bool {
	for _, s2 := range scene.Spheres {
		if s2 == sphere {
			continue
		}

		t2, ok := s2.Intersects(shadowRay)

		if !ok {
			continue
		}

		if t2 < 1 && t2 > 0 {
			return true
		}
	}

	return false
}

func rayColor(ray *Ray, scene *Scene, depth int) *Color {
	min := struct {
		*Sphere
		t float64
	}{}

	for _, s := range scene.Spheres {
		t, ok := s.Intersects(ray)

		if !ok {
			continue
		}

		if min.Sphere == nil || t < min.t {
			min.Sphere = s
			min.t = t
		}
	}

	if min.Sphere == nil {
		return black
	}

	color := min.Sphere.Material.Ambient.Times(scene.Ambient)

	intersectPoint := ray.Point.Add(ray.Direction.Scale(min.t))
	normalVec := intersectPoint.Sub(min.Sphere.Center).Normalize()

	for _, l := range scene.Lights {
		lightVec := l.Position.Sub(intersectPoint).Normalize()
		normalLightDot := normalVec.Dot(lightVec)

		if normalLightDot < 0 {
			continue
		}

		shadowRay := &Ray{
			Point:     intersectPoint,
			Direction: l.Position.Sub(intersectPoint),
		}

		if occluded(scene, min.Sphere, shadowRay) {
			continue
		}

		color = color.Add(
			min.Sphere.Material.Diffuse.
				Times(l.Diffuse).
				Scale(normalLightDot),
		)

		reflectanceVec := normalVec.
			Scale(2 * normalLightDot).
			Sub(lightVec)

		viewVec := scene.Camera.
			Sub(intersectPoint).
			Normalize()

		dotViewReflect := reflectanceVec.Dot(viewVec)

		specular := l.Specular.
			Times(min.Sphere.Material.Specular).
			Scale(
				math.Pow(dotViewReflect, min.Sphere.Material.Shininess),
			)

		color = color.Add(specular)

	}

	if depth < recurDepth {
		v := ray.Direction.Scale(-1).Normalize()
		reflectanceVec := normalVec.Scale(2 * v.Dot(normalVec)).Sub(v)

		newRay := &Ray{
			Point:     intersectPoint.Add(&Vec{.0001, .0001, .0001}),
			Direction: reflectanceVec,
		}

		color = color.Add(
			rayColor(newRay, scene, depth+1).
				Times(min.Sphere.Material.Reflectivity),
		)
	}

	return color.Clamp()
}

func colorAt(scene *Scene, alpha, beta float64) *Color {
	top := scene.X1.Lerp(scene.X2, alpha)
	bottom := scene.X3.Lerp(scene.X4, alpha)

	pointOnPlane := top.Lerp(bottom, beta)

	ray := &Ray{
		pointOnPlane,
		pointOnPlane.Sub(scene.Camera),
	}

	return rayColor(ray, scene, 0)
}

func renderAA(minAlpha, maxAlpha, minBeta, maxBeta float64, scene *Scene) *Color {
	color := &Color{}

	for sampX := 0; sampX < samples; sampX++ {
		for sampY := 0; sampY < samples; sampY++ {
			xt := float64(sampX) / float64(samples)
			yt := float64(sampY) / float64(samples)

			alpha := ((1 - xt) * minAlpha) + (maxAlpha * xt)
			beta := ((1 - yt) * minBeta) + (maxBeta * yt)

			color = color.Add(colorAt(scene, alpha, beta))
		}
	}

	color = color.Div(math.Pow(samples, 2))

	return color

}

func render(scene *Scene, screen *Screen) {
	for x := 0; x < screen.W; x++ {
		for y := 0; y < screen.H; y++ {
			minAlpha := float64(x) / float64(screen.W)
			minBeta := float64(y) / float64(screen.H)
			maxAlpha := float64(x+1) / float64(screen.W)
			maxBeta := float64(y+1) / float64(screen.H)

			screen.Set(
				x,
				y,
				renderAA(minAlpha, maxAlpha, minBeta, maxBeta, scene),
			)
		}
	}
}

func main() {
	width := flag.Int("width", 0, "width render target")
	height := flag.Int("height", 0, "height of render target")
	target := flag.String("target", "-", "height of render target")

	flag.Parse()

	screen := NewScreen(*width, *height)
	ratio := float64(*height) / float64(*width)

	scene := &Scene{
		Spheres: []*Sphere{
			// floor
			&Sphere{
				&Vec{0, -200, -17},
				200,
				&Material{
					Ambient:      &Color{0.3, 0.3, 0.3},
					Diffuse:      &Color{0.50, 0.50, 0.50},
					Specular:     &Color{0.00, 0.00, 0.00},
					Reflectivity: &Color{1.00, 1.00, 1.00},
					Shininess:    0,
				},
			},

			// red
			&Sphere{
				&Vec{-3.00, 1.00, -4.50},
				1,
				&Material{
					Ambient:      &Color{0.60, 0.20, 0.20},
					Diffuse:      &Color{0.60, 0.20, 0.20},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.80, 0.80},
					Shininess:    200,
				},
			},

			// yellow
			&Sphere{
				&Vec{-1.50, 2.75, -4.50},
				1,
				&Material{
					Ambient:      &Color{0.60, 0.60, 0.20},
					Diffuse:      &Color{0.60, 0.60, 0.20},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.80, 0.80},
					Shininess:    200,
				},
			},

			// green
			&Sphere{
				&Vec{0.00, 1.00, -4.50},
				1,
				&Material{
					Ambient:      &Color{0.20, 0.60, 0.20},
					Diffuse:      &Color{0.20, 0.60, 0.20},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.80, 0.80},
					Shininess:    200,
				},
			},

			// cyan
			&Sphere{
				&Vec{1.50, 2.75, -4.50},
				1,
				&Material{
					Ambient:      &Color{0.20, 0.60, 0.60},
					Diffuse:      &Color{0.20, 0.60, 0.60},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.80, 0.80},
					Shininess:    200,
				},
			},

			// blue
			&Sphere{
				&Vec{3.00, 1.00, -4.50},
				1,
				&Material{
					Ambient:      &Color{0.10, 0.20, 0.60},
					Diffuse:      &Color{0.10, 0.20, 0.60},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.90, 0.80},
					Shininess:    200,
				},
			},

			// purple
			&Sphere{
				&Vec{2.00, 0.00, -1.50},
				.2,
				&Material{
					Ambient:      &Color{0.60, 0.20, 0.60},
					Diffuse:      &Color{0.60, 0.20, 0.60},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.90, 0.80},
					Shininess:    200,
				},
			},
			// purple
			&Sphere{
				&Vec{2.25, 0.00, -1.50},
				.2,
				&Material{
					Ambient:      &Color{0.60, 0.20, 0.60},
					Diffuse:      &Color{0.60, 0.20, 0.60},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.90, 0.80},
					Shininess:    200,
				},
			},
			// purple
			&Sphere{
				&Vec{2.50, 0.00, -1.50},
				.2,
				&Material{
					Ambient:      &Color{0.60, 0.20, 0.60},
					Diffuse:      &Color{0.60, 0.20, 0.60},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.90, 0.80},
					Shininess:    200,
				},
			},
			// purple
			&Sphere{
				&Vec{2.75, 0.00, -1.50},
				.2,
				&Material{
					Ambient:      &Color{0.60, 0.20, 0.60},
					Diffuse:      &Color{0.60, 0.20, 0.60},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.90, 0.80},
					Shininess:    200,
				},
			},
			// purple
			&Sphere{
				&Vec{3.00, 0.00, -1.50},
				.2,
				&Material{
					Ambient:      &Color{0.60, 0.20, 0.60},
					Diffuse:      &Color{0.60, 0.20, 0.60},
					Specular:     &Color{0.90, 0.90, 0.90},
					Reflectivity: &Color{0.80, 0.90, 0.80},
					Shininess:    200,
				},
			},
		},
		Lights: []*Light{
			&Light{
				Position: &Vec{0, 0, 2},
				Diffuse:  &Color{0.10, 0.10, 0.10},
				Specular: &Color{0.10, 0.10, 0.10},
			},
			&Light{
				Position: &Vec{0, 2, 1},
				Diffuse:  &Color{0.80, 0.80, 0.80},
				Specular: &Color{0.80, 0.80, 0.80},
			},
		},
		Ambient: &Color{0.3, 0.3, 0.3},
		Camera:  &Vec{0, 0, 2},
		X1:      &Vec{-1, ratio, 1},
		X2:      &Vec{1, ratio, 1},
		X3:      &Vec{-1, -ratio, 1},
		X4:      &Vec{1, -ratio, 1},
	}

	renderStart := time.Now()
	render(scene, screen)
	renderDur := time.Since(renderStart).Nanoseconds()

	pushStart := time.Now()
	if *target == "-" {
		screen.Push()
	} else {
		screen.PushFile(*target)
	}
	pushDur := time.Since(pushStart).Nanoseconds()

	fmt.Println("render:", renderDur/1000000, "push:", pushDur/1000000)
}
