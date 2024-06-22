package logic

import (
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/scenes"
	"github.com/mikelangelon/town-sweet-town/world"
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

func (g GameLogic) CreateHouse(id string, typ int) house.House {
	return house.House{ID: id, House: *g.HouseFactory.Houses[typ]}
}

func (g GameLogic) NextDay(state scenes.State) scenes.State {
	day := state.Day + 1
	switch day {
	case 1:
		char := g.CharFactory.NewChar(1, []int{10, 111, 304}, 16*6, 16*6)

		fire := world.Fire{Char: g.FancyTownFactory.NewChar(470, nil, 16*6, 16*10)}
		sWest := house.NewSignal(
			g.TinyTownFactory.NewChar(83, nil, 16*4, 16*7),
			"west",
			common.Position{X: 3 * 16, Y: 3 * 16})
		sEast := house.NewSignal(
			g.TinyTownFactory.NewChar(83, nil, 16*19, 16*11),
			"west",
			common.Position{X: 14 * 16, Y: 7 * 16})
		sNorth := house.NewSignal(
			g.TinyTownFactory.NewChar(83, nil, 16*17, 16*4),
			"west",
			common.Position{X: 16 * 12, Y: 16 * 0})
		sSouth := house.NewSignal(
			g.TinyTownFactory.NewChar(83, nil, 16*8, 16*17),
			"west",
			common.Position{X: 16 * 7, Y: 16 * 12})
		return scenes.State{
			GameLogic: g,
			Player:    char,
			Status:    scenes.InitialState,
			World: map[string]*scenes.SceneMap{
				"town1": {
					Houses: []*house.House{},
					Objects: []world.Object{
						fire,
						sWest, sEast, sNorth, sSouth,
					},
				},
				"people": {
					NPCs: []*npc.NPC{
						g.NPCFactory.NewNPC(54, []int{11, 22}, 16*30, 16*6,
							&common.Position{X: 16 * 18, Y: 16 * 6},
							npc.WithCharacteristic(npc.Sports, npc.Cooking, npc.Animals).WithRent(3)),
						g.NPCFactory.NewNPC(271, nil, 16*28, 16*11,
							&common.Position{X: 16 * 17, Y: 16 * 11},
							npc.WithCharacteristic(npc.Extrovert, npc.Cooking, npc.Animals, npc.Reading).WithRent(3)),
						g.NPCFactory.NewNPC(162, []int{389, 476, 312}, 16*31, 16*9,
							&common.Position{X: 16 * 19, Y: 16 * 6},
							npc.WithCharacteristic(npc.Adventurous, npc.Music, npc.Extrovert, npc.Stuff, npc.Reading).WithRent(4)),
					},
				},
			},
		}
	case 2:
		state.World["people"].AddNPC(g.NPCFactory.NewNPC(1, []int{11, 101, 304}, 16*29, 16*8,
			&common.Position{X: 16 * 18, Y: 20 * 6},
			npc.WithCharacteristic(npc.Sports, npc.Cooking, npc.Animals)))
		state.World["people"].AddNPC(g.NPCFactory.NewNPC(2, []int{12, 104, 200}, 16*29, 16*5,
			&common.Position{X: 16 * 18, Y: 20 * 6},
			npc.WithCharacteristic(npc.Sports, npc.Music, npc.Sports)))
		state.World["people"].AddNPC(g.NPCFactory.NewNPC(1, []int{13, 300, 400}, 16*29, 16*12,
			&common.Position{X: 16 * 18, Y: 20 * 6},
			npc.WithCharacteristic(npc.Sports, npc.Cooking, npc.Extrovert)))
	}
	state.Day = day
	return state
}
