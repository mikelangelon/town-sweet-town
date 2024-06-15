package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mikelangelon/town-sweet-town/assets"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"log/slog"
)

func main() {
	town1, err := graphics.NewMapScene(assets.TileMapPacked, assets.Town1, assets.TileMapPackedTSX, screenWidth, screenHeight, scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	game := &Game{
		MapScene: town1,
	}
	ebiten.SetWindowSize(800, 600)
	if err := ebiten.RunGame(game); err != nil {
		slog.Error("something went wrong", "err", err)
	}
}
