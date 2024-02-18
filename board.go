package main

import (
	"fmt"
	"strings"
)

type Board struct {
	x, y   int
	h, w   int
	values [][]int
}

func NewBoard(x, y, h, w int) *Board {
	values := make([][]int, h)
	for i := range values {
		values[i] = make([]int, w)
	}

	return &Board{
		x:      x,
		y:      y,
		h:      h,
		w:      w,
		values: values,
	}
}

func (b *Board) At(x, y int) int {
	if x < 0 || x >= b.w || y < 0 || y >= b.h {
		return -1
	}
	return b.values[y][x]
}

func (b *Board) Set(x, y, v int) {
	if x < 0 || x >= b.w || y < 0 || y >= b.h {
		return
	}
	b.values[y][x] = v
}

func (b *Board) Height() int {
	return b.h
}

func (b *Board) Width() int {
	return b.w
}

func (b *Board) HasCollision(m *Mino) bool {
	for i := 0; i < m.Size(); i++ {
		for j := 0; j < m.Size(); j++ {
			if m.BlockAt(j, i) == 1 && b.At(m.x+j, m.y+i) != 0 {
				return true
			}
		}
	}
	return false
}

func (b *Board) Fix(m *Mino) {
	for i := 0; i < m.Size(); i++ {
		for j := 0; j < m.Size(); j++ {
			if m.BlockAt(j, i) == 1 {
				b.Set(m.x+j, m.y+i, int(m.Color()))
			}
		}
	}
}

func (b *Board) EraseLines(landedMino *Mino) {
	isEraced := make(map[int]bool, b.h)
	eracedCnt := 0
	for i := 0; i < landedMino.Size(); i++ {
		y := landedMino.y + i
		filled := true
		for j := 0; j < b.w; j++ {
			filled = filled && b.At(j, y) != 0
		}
		if filled {
			isEraced[y] = true
			eracedCnt++
		}
	}
	if eracedCnt == 0 {
		return
	}

	newValues := make([][]int, b.h)
	for i := range b.values {
		newValues[i] = make([]int, b.w)
	}
	idx := b.h - 1
	for i := b.h - 1; i >= 0; i-- {
		if isEraced[i] {
			continue
		}
		copy(newValues[idx], b.values[i])
		idx--
	}
	b.values = newValues
}

func (b *Board) Draw(m *Mino) {
	ss := make([]string, b.h)
	buf := make([]string, b.w)
	for i := 0; i < b.h; i++ {
		for j := 0; j < b.w; j++ {
			v := b.At(j, i)
			if m != nil && i >= m.y && i < m.y+m.Size() && j >= m.x && j < m.x+m.Size() && m.BlockAt(j-m.x, i-m.y) == 1 {
				v = int(m.Color())
			}
			if v == 0 {
				buf[j] = ansiColored("　", White)
			} else if v == -1 {
				buf[j] = ansiColored("　", Black)
			} else {
				buf[j] = ansiColored("　", Color(v))
			}
		}
		ss[i] = fmt.Sprintf("\033[%dC%s", b.x, strings.Join(buf, ""))
	}

	fmt.Printf("\033[1;1H\033[%dB%s", b.y, strings.Join(ss, "\n")+"\n")
}
