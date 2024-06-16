package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mikelangelon/town-sweet-town/npc"

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
	npcFactory := npc.NPCFactory{
		CharFactory: charFactory,
		Scale:       common.Scale,
	}
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	char := charFactory.NewChar(1, []int{10, 111, 304}, 16*6, 16*6)

	// set scenes
	state := scenes.State{
		Player: char,
		Status: scenes.InitialState,
	}
	people1Scene := &scenes.BaseScene{
		ID:       "people",
		MapScene: people1,
		NPCs: []*npc.NPC{
			npcFactory.NewNPC(1, []int{10, 111, 304}, 16*24, 16*6,
				common.Position{X: 16 * 18, Y: 16 * 6},
				npc.AddHappyCharacteristics(npc.Sports, npc.Cooking)),
			npcFactory.NewNPC(271, nil, 16*24, 16*11,
				common.Position{X: 16 * 17, Y: 16 * 11},
				npc.AddHappyCharacteristics(npc.Extrovert, npc.Cooking)),
			npcFactory.NewNPC(162, []int{389, 476, 312}, 16*26, 16*9,
				common.Position{X: 16 * 19, Y: 16 * 6},
				npc.AddHappyCharacteristics(npc.Adventurous, npc.Music)),
		},
	}
	town1Scene := &scenes.BaseScene{
		ID:       "town1",
		MapScene: town1,
	}
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
