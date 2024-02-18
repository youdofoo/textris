package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

const nextMinoCount = 3

type Game struct {
	board       *Board
	minoFigures []*MinoFigure
	currentMino *Mino
	score       int64

	nextMinos   []*Mino
	nextMinoIdx int
}

func NewGame(b *Board, minoFigures []*MinoFigure) *Game {
	g := &Game{
		board:       b,
		minoFigures: minoFigures,
	}
	g.nextMinos = make([]*Mino, nextMinoCount)
	for i := range g.nextMinos {
		g.nextMinos[i] = g.randomMino()
	}
	return g
}

func (g *Game) spawnMino() {
	g.currentMino = g.nextMinos[g.nextMinoIdx]
	g.nextMinos[g.nextMinoIdx] = g.randomMino()
	g.nextMinoIdx = (g.nextMinoIdx + 1) % len(g.nextMinos)
	g.drawNextMinos()
}

func (g *Game) randomMino() *Mino {
	return &Mino{
		figure: g.minoFigures[rand.Intn(len(g.minoFigures))],
		rot:    0,
		x:      (g.board.Width() - 1) / 2,
		y:      0,
	}
}

func (g *Game) handleKey(k Key) {
	switch k {
	case KeyUp:
		dy := 0
		for g.currentMino.CanMove(g.board, 0, dy+1) {
			dy++
		}
		g.currentMino.Move(0, dy)
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

func (g *Game) gameOver() {
	for i := 0; i < g.board.h; i++ {
		for j := 0; j < g.board.w; j++ {
			if g.board.At(j, i) != 0 {
				g.board.Set(j, i, int(Black))
			}
		}
	}
	g.board.Draw(nil)

	fmt.Printf("\033[1;1H\033[%dB\033[%dCGAME OVER!\033[1;1H\033[%dB", g.board.y+g.board.h/2, g.board.x+(g.board.w*2)/2-5, g.board.y+g.board.h)
}

func (g *Game) Run() error {
	clearDisplay()
	g.drawScore()

	keyInput := make(chan Key)
	go func() {
		err := CaptureKeyInput(keyInput)
		if err != nil {
			log.Fatal(err)
		}
	}()
	drawTicker := time.NewTicker(1000 / 60 * time.Millisecond)
	fallTicker := time.NewTicker(500 * time.Millisecond)

	g.spawnMino()
loop:
	for {
		select {
		case key := <-keyInput:
			if key == KeyEsc {
				g.gameOver()
				break loop
			}
			g.handleKey(key)
		case <-fallTicker.C:
			if g.currentMino.CanMove(g.board, 0, 1) {
				g.currentMino.Move(0, 1)
			} else {
				g.board.Fix(g.currentMino)
				eraced := g.board.EraseLines(g.currentMino)
				g.score += score(eraced)
				g.drawScore()
				g.spawnMino()
				if g.board.HasCollision(g.currentMino) {
					g.board.Fix(g.currentMino)
					g.gameOver()
					break loop
				}
			}
		case <-drawTicker.C:
			g.board.Draw(g.currentMino)
		}
	}
	return nil
}

func score(eraced int) int64 {
	switch eraced {
	case 1:
		return 100
	case 2:
		return 300
	case 3:
		return 500
	case 4:
		return 800
	default:
		return 0
	}
}

const (
	scoreOffsetX     = 6
	scoreY           = 5
	nextMinosOffsetX = 7
	nextMinosY       = 7
)

func (g *Game) drawScore() {
	fmt.Printf("\033[%d;%dHScore: %d", scoreY, boardX+boardWidth*2+scoreOffsetX, g.score)
}

func (g *Game) drawNextMinos() {
	fmt.Printf("\033[%d;%dH====NEXT====", nextMinosY, boardX+boardWidth*2+nextMinosOffsetX)
	for i := 0; i < len(g.nextMinos); i++ {
		idx := (g.nextMinoIdx + i) % len(g.nextMinos)
		g.drawNextMino(g.nextMinos[idx], i)
	}
}

func (g *Game) drawNextMino(mino *Mino, position int) {
	buf := make([]string, max(mino.Size(), 4))
	for i := 0; i < mino.Size(); i++ {
		for j := 0; j < max(mino.Size(), 4); j++ {
			var v int
			if j < mino.Size() {
				v = mino.BlockAt(j, i)
			} else {
				v = 0
			}

			if v == 0 {
				buf[j] = ansiColored("　", Black)
			} else {
				buf[j] = ansiColored("　", mino.Color())
			}
		}
		fmt.Printf("\033[%d;%dH%s", i+nextMinosY+2+position*4, boardX+boardWidth*2+nextMinosOffsetX, strings.Join(buf, ""))
	}
}
