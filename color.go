package main

import (
	"image/color"
	"math"
	"strconv"
)

type Color struct {
	R float64
	G float64
	B float64
}

func (c1 *Color) Times(c2 *Color) *Color {
	return &Color{
		c1.R * c2.R,
		c1.G * c2.G,
		c1.B * c2.B,
	}
}

func (c *Color) Scale(s float64) *Color {
	return &Color{
		c.R * s,
		c.G * s,
		c.B * s,
	}
}

func (c1 *Color) Add(c2 *Color) *Color {
	return &Color{
		c1.R + c2.R,
		c1.G + c2.G,
		c1.B + c2.B,
	}
}

func (c1 *Color) Div(s float64) *Color {
	return &Color{
		c1.R / s,
		c1.G / s,
		c1.B / s,
	}
}

func (c *Color) Clamp() *Color {
	r := math.Min(255, math.Max(0, c.R))
	g := math.Min(255, math.Max(0, c.G))
	b := math.Min(255, math.Max(0, c.B))

	return &Color{
		r,
		g,
		b,
	}
}

func (c *Color) ToRGBA() color.Color {
	r := uint8(math.Min(c.R*255, 255))
	g := uint8(math.Min(c.G*255, 255))
	b := uint8(math.Min(c.B*255, 255))

	return color.RGBA{
		r,
		g,
		b,
		255,
	}
}

const shadeChars = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
const maxBright = 3 * 255

func (c *Color) ToChar() string {
	r := int(math.Min(c.R*255, 255))
	g := int(math.Min(c.G*255, 255))
	b := int(math.Min(c.B*255, 255))

	shade := (1 - float64(r+g+b)/float64(maxBright)) * float64((len(shadeChars) - 1))

	char := shadeChars[int(shade)]

	return "\x1b[38;2;" +
		strconv.Itoa(r) + ";" +
		strconv.Itoa(g) + ";" +
		strconv.Itoa(b) + "m" +
		string(char) + "\x1b[0m"
}
