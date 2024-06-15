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
	charFactory, err := graphics.NewCharFactory(assets.Characters, assets.CharactersTSX, scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	char := charFactory.NewChar(1, []int{10, 111, 304}, 80, 120)
	game := &Game{
		MapScene: town1,
		Player:   char,
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		slog.Error("something went wrong", "err", err)
	}
}
