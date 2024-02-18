package main

import (
	"fmt"
	"log"
)

type Color uint8

const (
	Black Color = iota + 1
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func ansiColored(s string, c Color) string {
	return fmt.Sprintf("\033[%dm%s\033[m", 39+c, s)
}

func clearDisplay() {
	fmt.Print("\033[2J")
}

const (
	boardX      = 2
	boardY      = 2
	boardWidth  = 10
	boardHeight = 20
)

func main() {
	board := NewBoard(boardX, boardY, boardWidth, boardHeight)
	minoFigures := makeMinoFigures()
	g := NewGame(board, minoFigures)

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
