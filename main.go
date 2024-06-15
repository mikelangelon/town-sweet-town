package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/assets"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/scenes"
	"log/slog"
)

func main() {

	town1, err := graphics.NewMapScene(assets.TileMapPacked, assets.Town1, assets.TileMapPackedTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	town2, err := graphics.NewMapScene(assets.TileMapPacked, assets.Town2, assets.TileMapPackedTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	charFactory, err := graphics.NewCharFactory(assets.Characters, assets.CharactersTSX, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	char := charFactory.NewChar(1, []int{10, 111, 304}, 16*6, 16*6)

	state := scenes.State{
		Player: char,
	}
	sm := stagehand.NewSceneManager[scenes.State](&scenes.Town1Scene{
		MapScene: town1,
		TransitionPoints: map[common.Position]stagehand.Scene[scenes.State]{
			common.Position{X: 24 * 16, Y: 6 * 16}: &scenes.Town2Scene{
				MapScene: town2,
			},
		},
	}, state)
	ebiten.SetWindowSize(common.ScreenWidth, common.ScreenHeight)
	if err := ebiten.RunGame(sm); err != nil {
		slog.Error("something went wrong", "err", err)
	}
}
