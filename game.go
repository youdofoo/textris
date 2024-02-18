package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Game struct {
	board       *Board
	minoFigures []*MinoFigure
	currentMino *Mino
	score       int64
}

func NewGame(b *Board, minoFigures []*MinoFigure) *Game {
	return &Game{
		board:       b,
		minoFigures: minoFigures,
	}
}

func (g *Game) spawnMino() {
	idx := rand.Intn(len(g.minoFigures))
	g.currentMino = &Mino{
		figure: g.minoFigures[idx],
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
	scoreOffsetX = 6
	scoreY       = 5
)

func (g *Game) drawScore() {
	fmt.Printf("\033[%d;%dHScore: %d", scoreY, boardX+boardWidth*2+scoreOffsetX, g.score)
}
