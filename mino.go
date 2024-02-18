package main

import (
	"errors"
	"log"
)

type Blocks [][]int

func NewBlocks(size int) Blocks {
	blocks := make([][]int, size)
	for i := range blocks {
		blocks[i] = make([]int, size)
	}
	return blocks
}

type MinoFigure struct {
	rotations []Blocks
	size      int
	color     Color
}

func NewMinoFigure(base Blocks, color Color) (*MinoFigure, error) {
	if len(base) == 0 {
		return nil, errors.New("empty blocks")
	}
	for i := range base {
		if len(base[i]) != len(base) {
			return nil, errors.New("invalid shape, blocks must be a square")
		}
	}
	size := len(base)
	rotations := make([]Blocks, 4)
	rotations[0] = NewBlocks(size)
	copy(rotations[0], base)

	prevBlocks := rotations[0]
	for r := 1; r < 4; r++ {
		blocks := NewBlocks(size)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				blocks[i][j] = prevBlocks[size-j-1][i]
			}
		}
		rotations[r] = blocks
		prevBlocks = rotations[r]
	}
	return &MinoFigure{
		rotations: rotations,
		size:      size,
		color:     color,
	}, nil
}

type Mino struct {
	figure *MinoFigure
	rot    int
	x, y   int
}

func (m *Mino) BlockAt(x, y int) int {
	return m.figure.rotations[m.rot][y][x]
}

func (m *Mino) Blocks() Blocks {
	return m.figure.rotations[m.rot]
}

func (m *Mino) Size() int {
	return m.figure.size
}

func (m *Mino) Color() Color {
	return m.figure.color
}

func (m *Mino) Move(dx, dy int) {
	m.x += dx
	m.y += dy
}

func (m *Mino) RotateL() {
	m.rot = (m.rot - 1 + 4) % 4
}

func (m *Mino) RotateR() {
	m.rot = (m.rot + 1) % 4
}

func (m *Mino) CanMove(b *Board, dx, dy int) bool {
	cp := *m
	cp.x += dx
	cp.y += dy
	return !b.HasCollision(&cp)
}

func (m *Mino) CanRotateL(b *Board) bool {
	cp := *m
	cp.rot = (cp.rot - 1 + 4) % 4
	return !b.HasCollision(&cp)
}

func (m *Mino) CanRotateR(b *Board) bool {
	cp := *m
	cp.rot = (cp.rot + 1) % 4
	return !b.HasCollision(&cp)
}

func makeMinoFigures() []*MinoFigure {
	values := []struct {
		base  Blocks
		color Color
	}{
		{
			// Iミノ
			base: Blocks{
				{0, 0, 0, 0},
				{1, 1, 1, 1},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			color: Cyan,
		},
		{
			// Oミノ
			base: Blocks{
				{1, 1},
				{1, 1},
			},
			color: Yellow,
		},
		{
			// Sミノ
			base: Blocks{
				{0, 1, 1},
				{1, 1, 0},
				{0, 0, 0},
			},
			color: Green,
		},
		{
			// Zミノ
			base: Blocks{
				{1, 1, 0},
				{0, 1, 1},
				{0, 0, 0},
			},
			color: Red,
		},
		{
			// Jミノ
			base: Blocks{
				{0, 0, 1},
				{1, 1, 1},
				{0, 0, 0},
			},
			color: Blue,
		},
		{
			// Lミノ
			base: Blocks{
				{1, 0, 0},
				{1, 1, 1},
				{0, 0, 0},
			},
			color: Yellow,
		},
		{
			// Tミノ
			base: Blocks{
				{0, 1, 0},
				{1, 1, 1},
				{0, 0, 0},
			},
			color: Magenta,
		},
	}

	var err error
	figures := make([]*MinoFigure, len(values))
	for i := range values {
		figures[i], err = NewMinoFigure(values[i].base, values[i].color)
		if err != nil {
			log.Fatal(err)
		}
	}
	return figures
}
