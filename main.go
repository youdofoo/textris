package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	BoardWidth       = 10
	BoardHeight      = 20
	BoardMargin      = 2
	BoardArrayHeight = BoardHeight + 2*BoardMargin
	BoardArrayWidth  = BoardWidth + 2*BoardMargin
)

type Color uint8

const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func ansiColored(s string, c Color) string {
	return fmt.Sprintf("\033[%dm%s\033[m", 40+c, s)
}

type Board struct {
	x, y   int
	h, w   int
	margin int
	values [][]int
}

func NewBoard(x, y, h, w, margin int) *Board {
	values := make([][]int, h+2*margin)
	for i := range values {
		values[i] = make([]int, w+2*margin)
	}
	for i := 0; i < h+2*margin; i++ {
		for m := 0; m < margin; m++ {
			values[i][m] = -1
			values[i][w+m] = -1
		}
	}
	for j := 0; j < w+2*margin; j++ {
		for m := 0; m < margin; m++ {
			values[m][j] = -1
			values[h+m][j] = -1
		}
	}

	return &Board{
		x:      x,
		y:      y,
		h:      h,
		w:      w,
		margin: margin,
		values: values,
	}
}

func (b *Board) At(x, y int) int {
	if x < 0 || x >= b.w || y < 0 || y >= b.h {
		return -1
	}
	return b.values[y+b.margin][x+b.margin]
}

func (b *Board) Set(x, y, v int) {
	if x < 0 || x >= b.w || y < 0 || y >= b.h {
		return
	}
	b.values[y+b.margin][x+b.margin] = v
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
			if i >= m.y && i < m.y+m.Size() && j >= m.x && j < m.x+m.Size() && m.BlockAt(j-m.x, i-m.y) == 1 {
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

func clearDisplay() {
	fmt.Print("\033[2J")
}

type Game struct {
	board       *Board
	minoFigures []*MinoFigure
	currentMino *Mino
}

func NewGame(b *Board, minoFigures []*MinoFigure) *Game {
	return &Game{
		board:       b,
		minoFigures: minoFigures,
	}
}

func (g *Game) SpawnMino() {
	idx := rand.Intn(len(g.minoFigures))
	g.currentMino = &Mino{
		figure: g.minoFigures[idx],
		rot:    0,
		x:      (g.board.Width() - 1) / 2,
		y:      0,
	}
}

func (g *Game) FallMino() bool {
	if g.currentMino.CanMove(g.board, 0, 1) {
		g.currentMino.Move(0, 1)
		return true
	} else {
		return false
	}
}

func (g *Game) OnMinoLanded() {
	g.board.Fix(g.currentMino)
	g.board.EraseLines(g.currentMino)
}

func (g *Game) FixMino() {
	g.board.Fix(g.currentMino)
}

func (g *Game) Draw() {
	g.board.Draw(g.currentMino)
}

func (g *Game) HandleKey(k Key) {
	switch k {
	case KeyUp:
		// TODO: ハードドロップ
	case KeyDown:
		if g.currentMino.CanMove(g.board, 0, 1) {
			g.currentMino.Move(0, 1)
		}
	case KeyRight:
		if g.currentMino.CanMove(g.board, 1, 0) {
			g.currentMino.Move(1, 0)
		}
	case KeyLeft:
		if g.currentMino.CanMove(g.board, -1, 0) {
			g.currentMino.Move(-1, 0)
		}
	case KeyA:
		if g.currentMino.CanRotateL(g.board) {
			g.currentMino.RotateL()
		}
	case KeyD:
		if g.currentMino.CanRotateR(g.board) {
			g.currentMino.RotateR()
		}
	}
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

func main() {
	clearDisplay()
	keyInput := make(chan Key)
	go func() {
		err := CaptureKeyInput(keyInput)
		if err != nil {
			log.Fatal(err)
		}
	}()
	drawTicker := time.NewTicker(1000 / 60 * time.Millisecond)
	fallTicker := time.NewTicker(500 * time.Millisecond)

	board := NewBoard(2, 2, 20, 10, 0)
	minoFigures := makeMinoFigures()
	g := NewGame(board, minoFigures)
	g.SpawnMino()

loop:
	for {
		select {
		case key := <-keyInput:
			if key == KeyEsc {
				break loop
			}
			g.HandleKey(key)
		case <-fallTicker.C:
			if !g.FallMino() {
				g.OnMinoLanded()
				g.SpawnMino()
			}
		case <-drawTicker.C:
			g.Draw()
		default:
		}
	}
}
