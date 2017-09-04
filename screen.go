package main

import (
	"image"
	"image/png"
	"os"
)

type Screen struct {
	cells [][]*Color
	W     int
	H     int
}

func NewScreen(w, h int) *Screen {
	s := &Screen{
		cells: make([][]*Color, h),
		W:     w,
		H:     h,
	}

	for y, _ := range s.cells {
		s.cells[y] = make([]*Color, w)
	}

	return s
}

func (s *Screen) Push() {
	final := ""

	for _, row := range s.cells {
		for _, pix := range row {
			final += pix.ToChar()
		}

		final += "\n"
	}

	os.Stdout.Write([]byte("\033[?25l" + "\033[0;0H" + final))
}

func (s *Screen) PushFile(to string) {
	rendered := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: s.W, Y: s.H},
	})

	for y, row := range s.cells {
		for x, pix := range row {
			rendered.Set(x, y, pix.ToRGBA())
		}
	}

	file, err := os.Create(to)
	if err != nil {
		panic(err)
	}

	err = png.Encode(file, rendered)
	if err != nil {
		panic(err)
	}

	file.Close()
}

func (s *Screen) Set(x, y int, c *Color) {
	s.cells[y][x] = c
}
