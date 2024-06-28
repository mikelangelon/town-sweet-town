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
	RulesApplier     *npc.RuleApplier
}

func (g GameLogic) GetRuler() npc.RuleApplier {
	return *g.RulesApplier
}

func (g GameLogic) CreateHouse(id string, typ int) house.House {
	return house.House{ID: id, House: *g.HouseFactory.Houses[typ]}
}

func (g GameLogic) ChangePlayer(state scenes.State) scenes.State {
	char := g.CharFactory.NewChar(1, []int{10, 111, 304}, 16*6, 16*6)
	state.Player.Stuff = char.Stuff
	return state
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
		index := rand.Intn(len(positionAvailable))
		pos := positionAvailable[index]
		positionAvailable = append(positionAvailable[:index], positionAvailable[index+1:]...)
		return pos
	}

	switch day {
	case 1:
		char := g.CharFactory.NewChar(1, nil, 16*6, 16*6)

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
						g.NPCFactory.NewNPC(54, []int{11, 22}, getPositionAvailable(), npc.WithRandom(3)),
						g.NPCFactory.NewNPC(54, []int{11, 22}, getPositionAvailable(), npc.WithRandom(3)),
						g.NPCFactory.NewNPC(271, nil, getPositionAvailable(), npc.WithRandom(3)),
						g.NPCFactory.NewNPC(162, []int{389, 476, 312}, getPositionAvailable(), npc.WithRandom(4)),
					},
				},
			},
		}
	case 3:
		entrance.AddNPC(g.NPCFactory.NewNPC(1, []int{11, 101, 304}, getPositionAvailable(), npc.WithRandom(4)))
		entrance.AddNPC(g.NPCFactory.NewNPC(54, []int{12, 104, 561}, getPositionAvailable(), npc.WithRandom(4)))
		entrance.AddNPC(g.NPCFactory.NewNPC(109, []int{13, 300, 197}, getPositionAvailable(), npc.WithRandom(4)))
	case 5:
		entrance.AddNPC(g.NPCFactory.NewNPC(1, []int{11, 101}, getPositionAvailable(), npc.WithRandom(5)))
		entrance.AddNPC(g.NPCFactory.NewNPC(0, []int{12, 104, 478}, getPositionAvailable(), npc.WithRandom(5)))
	case 7:
		entrance.AddNPC(g.NPCFactory.NewNPC(1, []int{11, 101}, getPositionAvailable(), npc.WithRandom(5)))
		entrance.AddNPC(g.NPCFactory.NewNPC(1, []int{12, 104, 478}, getPositionAvailable(), npc.WithRandom(7)))
		entrance.AddNPC(g.NPCFactory.NewNPC(163, []int{12}, getPositionAvailable(), npc.WithRandom(8)))
	case 9:
		entrance.AddNPC(g.NPCFactory.NewNPC(487, []int{}, getPositionAvailable(), npc.WithRandom(5)))
		entrance.AddNPC(g.NPCFactory.NewNPC(108, []int{12, 104, 478}, getPositionAvailable(), npc.WithRandom(7)))
	case 11:
		entrance.AddNPC(g.NPCFactory.NewNPC(0, []int{11, 101}, getPositionAvailable(), npc.WithRandom(6)))
		entrance.AddNPC(g.NPCFactory.NewNPC(1, []int{12, 104, 478}, getPositionAvailable(), npc.WithRandom(7)))
		entrance.AddNPC(g.NPCFactory.NewNPC(163, []int{12}, getPositionAvailable(), npc.WithRandom(8)))
	case 13:
		entrance.AddNPC(g.NPCFactory.NewNPC(54, []int{11, 101}, getPositionAvailable(), npc.WithRandom(7)))
		entrance.AddNPC(g.NPCFactory.NewNPC(163, []int{12, 104, 478}, getPositionAvailable(), npc.WithRandom(8)))
	}
	g.RulesApplier.Rules = append(g.RulesApplier.Rules, npc.RandomRule())
	state.Day = day
	return state
}
