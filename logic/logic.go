package logic

import (
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/scenes"
	"github.com/mikelangelon/town-sweet-town/world"
	"github.com/mikelangelon/town-sweet-town/world/house"
	"github.com/mikelangelon/town-sweet-town/world/npc"
	"math/rand"
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
	entrance := state.World["people"]
	day := state.Day + 1

	const offsetX = 13
	const offsetY = 4
	var positionAvailable []common.Position
	for x := 1; x < 10; x++ {
		for y := 1; y < 14; y++ {
			positionAvailable = append(positionAvailable, common.Position{X: int64(x+offsetX) * 16, Y: int64(y+offsetY) * 16})
		}
	}

	getPositionAvailable := func() common.Position {
		index := rand.Intn(len(positionAvailable) + 1)
		pos := positionAvailable[index]
		positionAvailable = append(positionAvailable[:index], positionAvailable[index+1:]...)
		return pos
	}

	switch day {
	case 1:
		char := g.CharFactory.NewChar(1, []int{10, 111, 304}, 16*6, 16*6)

		fire := world.Fire{Char: g.FancyTownFactory.NewChar(470, nil, 16*6, 16*10)}
		sWest := house.NewSignal(
			g.TinyTownFactory.NewChar(83, nil, 16*4, 16*7),
			"West",
			common.Position{X: 3 * 16, Y: 3 * 16})
		sEast := house.NewSignal(
			g.TinyTownFactory.NewChar(83, nil, 16*19, 16*11),
			"East",
			common.Position{X: 14 * 16, Y: 7 * 16})
		sNorth := house.NewSignal(
			g.TinyTownFactory.NewChar(83, nil, 16*17, 16*4),
			"North",
			common.Position{X: 16 * 12, Y: 16 * 0})
		sSouth := house.NewSignal(
			g.TinyTownFactory.NewChar(83, nil, 16*8, 16*17),
			"South",
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
						g.NPCFactory.NewNPC(54, []int{11, 22}, getPositionAvailable(),
							npc.WithCharacteristic(npc.Sports, npc.Cooking, npc.Animals).WithRent(3)),
						g.NPCFactory.NewNPC(271, nil, getPositionAvailable(),
							npc.WithCharacteristic(npc.Extrovert, npc.Cooking, npc.Animals, npc.Reading).WithRent(3)),
						g.NPCFactory.NewNPC(162, []int{389, 476, 312}, getPositionAvailable(),
							npc.WithCharacteristic(npc.Adventurous, npc.Music, npc.Extrovert, npc.Stuff, npc.Reading).WithRent(4)),
					},
				},
			},
		}
	case 3:
		entrance.AddNPC(g.NPCFactory.NewNPC(1, []int{11, 101, 304}, getPositionAvailable(),
			npc.WithCharacteristic(npc.Sports, npc.Cooking, npc.Animals)))
		entrance.AddNPC(g.NPCFactory.NewNPC(2, []int{12, 104, 200}, getPositionAvailable(),
			npc.WithCharacteristic(npc.Sports, npc.Music, npc.Sports)))
		entrance.AddNPC(g.NPCFactory.NewNPC(1, []int{13, 300, 400}, getPositionAvailable(),
			npc.WithCharacteristic(npc.Sports, npc.Cooking, npc.Extrovert)))
	case 5:
		entrance.AddNPC(g.NPCFactory.NewNPC(1, []int{11, 101, 304}, getPositionAvailable(),
			npc.WithCharacteristic(npc.Reading, npc.Competitive, npc.Animals)))
		entrance.AddNPC(g.NPCFactory.NewNPC(2, []int{12, 104, 200}, getPositionAvailable(),
			npc.WithCharacteristic(npc.Money, npc.Animals, npc.Workaholic)))
	}
	state.Day = day
	return state
}
