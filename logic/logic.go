package logic

import (
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/scenes"
	"github.com/mikelangelon/town-sweet-town/world/house"
	"github.com/mikelangelon/town-sweet-town/world/npc"
)

type GameLogic struct {
	HouseFactory     *graphics.HouseFactory
	NPCFactory       npc.NPCFactory
	TinyTownFactory  *graphics.CharFactory
	CharFactory      *graphics.CharFactory
	FancyTownFactory *graphics.CharFactory
}

func (g GameLogic) NextDay(state scenes.State) scenes.State {
	day := state.Day + 1
	switch day {
	case 1:
		char := g.CharFactory.NewChar(1, []int{10, 111, 304}, 16*6, 16*6)
		fire := g.FancyTownFactory.NewChar(470, nil, 16*6, 16*10)
		signal1 := g.TinyTownFactory.NewChar(83, nil, 16*4, 16*7)
		signal2 := g.TinyTownFactory.NewChar(83, nil, 16*19, 16*11)
		signal3 := g.TinyTownFactory.NewChar(83, nil, 16*17, 16*4)
		signal4 := g.TinyTownFactory.NewChar(83, nil, 16*11, 16*17)
		return scenes.State{
			GameLogic: g,
			Player:    char,
			Status:    scenes.InitialState,
			World: map[string]*scenes.SceneMap{
				"town1": {
					Houses: []*house.House{
						{ID: "House 1", House: g.HouseFactory.Houses[0]},
						{ID: "House 2", House: g.HouseFactory.Houses[1]},
					},
					Objects: []*graphics.Char{
						fire,
						signal1, signal2, signal3, signal4,
					},
				},
				"people": {
					NPCs: []*npc.NPC{
						g.NPCFactory.NewNPC(1, []int{10, 111, 304}, 16*30, 16*6,
							&common.Position{X: 16 * 18, Y: 16 * 6},
							npc.AddHappyCharacteristics(npc.Sports, npc.Cooking, npc.Animals)),
						g.NPCFactory.NewNPC(271, nil, 16*28, 16*11,
							&common.Position{X: 16 * 17, Y: 16 * 11},
							npc.AddHappyCharacteristics(npc.Extrovert, npc.Cooking, npc.Animals, npc.Reading)),
						g.NPCFactory.NewNPC(162, []int{389, 476, 312}, 16*31, 16*9,
							&common.Position{X: 16 * 19, Y: 16 * 6},
							npc.AddHappyCharacteristics(npc.Adventurous, npc.Music)),
					},
				},
			},
		}
	case 2:
		state.World["people"].AddNPC(g.NPCFactory.NewNPC(1, []int{11, 101, 304}, 16*29, 16*8,
			&common.Position{X: 16 * 18, Y: 20 * 6},
			npc.AddHappyCharacteristics(npc.Sports, npc.Cooking, npc.Animals)))
		state.World["people"].AddNPC(g.NPCFactory.NewNPC(2, []int{12, 104, 200}, 16*29, 16*5,
			&common.Position{X: 16 * 18, Y: 20 * 6},
			npc.AddHappyCharacteristics(npc.Sports, npc.Cooking, npc.Animals)))
		state.World["people"].AddNPC(g.NPCFactory.NewNPC(1, []int{13, 300, 400}, 16*29, 16*12,
			&common.Position{X: 16 * 18, Y: 20 * 6},
			npc.AddHappyCharacteristics(npc.Sports, npc.Cooking, npc.Animals)))
	}
	state.Day = day
	return state
}
