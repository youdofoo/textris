package main

import (
	"fmt"

	"github.com/mattn/go-tty"
)

type Key int

const (
	KeyUp Key = iota
	KeyDown
	KeyRight
	KeyLeft
	KeyEsc
	KeyA
	KeyD
	KeyUnknown = -1
)

func CaptureKeyInput(input chan<- Key) error {
	tty, err := tty.Open()
	if err != nil {
		return fmt.Errorf("failed to open TTY: %w", err)
	}
	defer tty.Close()

	for {
		buf := make([]rune, 0, 10)
		for {
			r, err := tty.ReadRune()
			if err != nil {
				return fmt.Errorf("failed to read rune from TTY: %w", err)
			}
			buf = append(buf, r)
			if !tty.Buffered() {
				break
			}
		}
		k := toKey(buf)
		if k != KeyUnknown {
			input <- k
		}
		if k == KeyEsc {
			break
		}
	}
	return nil
}

func toKey(rs []rune) Key {
	if len(rs) == 3 && rs[0] == 0x1b && rs[1] == 0x5b {
		// 方向キー
		switch rs[2] {
		case 0x41:
			return KeyUp
		case 0x42:
			return KeyDown
		case 0x43:
			return KeyRight
		case 0x44:
			return KeyLeft
		default:
			return KeyUnknown
		}
	}
	if len(rs) == 1 {
		switch rs[0] {
		case 0x1b:
			return KeyEsc
		case 'a':
			return KeyA
		case 'd':
			return KeyD
		}
	}

	return KeyUnknown
}
