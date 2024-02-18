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

func main() {
	board := NewBoard(2, 2, 20, 10)
	minoFigures := makeMinoFigures()
	g := NewGame(board, minoFigures)

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
