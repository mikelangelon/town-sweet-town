package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mikelangelon/town-sweet-town/world/house"
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
	npcFactory := npc.NPCFactory{
		CharFactory: charFactory,
		Scale:       common.Scale,
	}
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	char := charFactory.NewChar(1, []int{10, 111, 304}, 16*6, 16*6)
	fire := fancyTownFactory.NewChar(470, nil, 16*6, 16*10)
	// set scenes
	state := scenes.State{
		Player: char,
		Status: scenes.InitialState,
		World: map[string]*scenes.SceneMap{
			"town1": {
				Houses: []*house.House{
					{ID: "House 1", Position: common.Position{X: 6 * 16, Y: 6 * 16}},
					{ID: "House 2", Position: common.Position{X: 17 * 16, Y: 11 * 16}},
				},
				Objects: []*graphics.Char{
					fire,
				},
			},
			"people": {
				NPCs: []*npc.NPC{
					npcFactory.NewNPC(1, []int{10, 111, 304}, 16*30, 16*6,
						&common.Position{X: 16 * 18, Y: 16 * 6},
						npc.AddHappyCharacteristics(npc.Sports, npc.Cooking, npc.Animals)),
					npcFactory.NewNPC(271, nil, 16*28, 16*11,
						&common.Position{X: 16 * 17, Y: 16 * 11},
						npc.AddHappyCharacteristics(npc.Extrovert, npc.Cooking, npc.Animals, npc.Reading)),
					npcFactory.NewNPC(162, []int{389, 476, 312}, 16*31, 16*9,
						&common.Position{X: 16 * 19, Y: 16 * 6},
						npc.AddHappyCharacteristics(npc.Adventurous, npc.Music)),
				},
			},
		},
	}
	people1Scene := scenes.NewEntrance("people", people1)
	town1Scene := scenes.NewTown("town1", town1)
	town1Scene.SetupUI()
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
