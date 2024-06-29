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

	positionAvailable []common.Position
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
func (g GameLogic) Reset() {
	g.positionAvailable = nil
	g.NPCFactory.Reset()
}
func (g GameLogic) NextDay(state scenes.State) scenes.State {
	entrance := state.World["people"]
	day := state.Day + 1

	const offsetX = 13
	const offsetY = 4
	if len(g.positionAvailable) == 0 {
		var positionAvailable []common.Position
		for x := 1; x < 10; x++ {
			for y := 1; y < 14; y++ {
				positionAvailable = append(positionAvailable, common.Position{X: int64(x+offsetX) * 16, Y: int64(y+offsetY) * 16})
			}
		}
		g.positionAvailable = positionAvailable
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
		sPeople := house.NewSignal(
			g.TinyTownFactory.NewChar(83, nil, 16*3, 16*4),
			"People",
			common.Position{X: 16 * 3, Y: 16 * 4})
		sPeople.Cost = 2

		stats := make(map[string]int)
		stats[npc.Money] = 10 * ratio(state.Difficulty)
		stats[npc.Happiness] = 0
		stats[npc.Food] = 15 * ratio(state.Difficulty) / 2
		stats[npc.Security] = 10
		stats[npc.Health] = 10
		stats[npc.Cultural] = 10
		options := []string{npc.Cultural, npc.Health, npc.Security}
		goals := []world.Goal{
			{
				Day:       5,
				Stat:      options[rand.Intn(3)],
				Value:     30 - 10*(ratio(state.Difficulty)-1)/2,
				GiftStat:  npc.Food,
				GiftValue: 15,
			},
			{
				Day:       10,
				Stat:      options[rand.Intn(3)],
				Value:     40 - 10*(ratio(state.Difficulty)-1)/2,
				GiftStat:  npc.Food,
				GiftValue: 30,
			},
			{
				Day:       14,
				Stat:      npc.Happiness,
				Value:     100,
				Mandatory: true,
			},
		}
		return scenes.State{
			GameLogic: g,
			Player:    char,
			Goals:     goals,
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
						g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(3)),
						g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(3)),
						g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(3)),
						g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(4)),
					},
					Objects: []world.Object{sPeople},
				},
			},
			Stats:         stats,
			Day:           day,
			MenuSong:      state.MenuSong,
			TownSillySong: state.TownSillySong,
		}
	case 3:
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(4)))
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(4)))
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(4)))
	case 5:
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(5)))
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(5)))
	case 7:
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(5)))
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(7)))
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(8)))
	case 9:
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(5)))
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(7)))
	case 11:
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(6)))
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(7)))
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(8)))
	case 13:
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(7)))
		entrance.AddNPC(g.NPCFactory.NewRandomGuy(g.GetPositionAvailable(), npc.WithRandom(8)))

	}
	g.RulesApplier.Rules = append(g.RulesApplier.Rules, npc.RandomRule())
	state.Day = day
	return state
}

func ratio(difficulty string) int {
	switch difficulty {
	case "Hard":
		return 1
	case "Normal":
		return 2
	case "Easy":
		return 3
	}
	return 0
}

func (g GameLogic) AddNPC() *npc.NPC {
	options := []int{270, 271, 324, 325, 378, 379, 432, 433, 486, 487, 540, 541, 594, 595}
	return g.NPCFactory.NewNPC(options[rand.Intn(len(options))], nil, g.GetPositionAvailable(), npc.WithRandom(5))
}

func (g GameLogic) GetPositionAvailable() common.Position {
	index := rand.Intn(len(g.positionAvailable))
	pos := g.positionAvailable[index]
	g.positionAvailable = append(g.positionAvailable[:index], g.positionAvailable[index+1:]...)
	return pos
}
