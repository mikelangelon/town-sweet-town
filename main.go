package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log/slog"
)

func main() {
	game := &Game{}
	ebiten.SetWindowSize(800, 600)
	if err := ebiten.RunGame(game); err != nil {
		slog.Error("something went wrong", "err", err)
	}
}
