package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mikelangelon/town-sweet-town/logic"
	"github.com/mikelangelon/town-sweet-town/world/npc"

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
	houseFactory, err := graphics.NewHouseFactory()
	people1, err := graphics.NewMapScene(assets.TileMapPacked, assets.People1, assets.TileMapPackedTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	charFactory, err := graphics.NewCharFactory(assets.Characters, assets.CharactersTSX, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	fancyTownFactory, err := graphics.NewCharFactory(assets.FancyTown, assets.FancyTownTSX, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	tinyTownFactory, err := graphics.NewCharFactory(assets.TileMapPacked, assets.TileMapPackedTSX, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	npcFactory := npc.NPCFactory{
		CharFactory: charFactory,
		Scale:       common.Scale,
	}
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	gameLogic := logic.GameLogic{
		NPCFactory:       npcFactory,
		CharFactory:      charFactory,
		FancyTownFactory: fancyTownFactory,
		TinyTownFactory:  tinyTownFactory,
		HouseFactory:     houseFactory,
	}
	state := gameLogic.NextDay(scenes.State{})
	people1Scene := scenes.NewEntrance("people", people1)
	town1Scene := scenes.NewTown("town1", town1)
	town1Scene.TransitionPoints = scenes.Transition{
		Position:  common.Position{X: 24 * 16, Y: 6 * 16},
		Scene:     people1Scene,
		Direction: stagehand.RightToLeft,
	}
	people1Scene.TransitionPoints = scenes.Transition{
		Position:  common.Position{X: 0 * 16, Y: 6 * 16},
		Scene:     town1Scene,
		Direction: stagehand.LeftToRight,
	}
	sm := stagehand.NewSceneManager[scenes.State](town1Scene, state)
	ebiten.SetWindowSize(common.ScreenWidth, common.ScreenHeight)
	if err := ebiten.RunGame(sm); err != nil {
		slog.Error("something went wrong", "err", err)
	}
}
