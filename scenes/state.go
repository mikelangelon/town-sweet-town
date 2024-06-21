package scenes

import (
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/world/house"
	"github.com/mikelangelon/town-sweet-town/world/npc"
)

type State struct {
	Player    *graphics.Char
	Status    int
	World     map[string]*SceneMap
	Stats     map[string]int
	Day       int
	GameLogic Brainer
}

type Brainer interface {
	NextDay(state State) State
}

type SceneMap struct {
	NPCs    npc.NPCs
	Houses  []*house.House
	Objects []*graphics.Char
}

func (s *SceneMap) AddNPC(npc *npc.NPC) {
	s.NPCs = append(s.NPCs, npc)
}
func (s *SceneMap) RemoveNPC(ID string) {
	var npcs []*npc.NPC
	for _, v := range s.NPCs {
		if v.ID == ID {
			continue
		}
		npcs = append(npcs, v)
	}
	s.NPCs = npcs
}

const (
	Menu = iota
	InitialState
	Playing
	Transitioning
	DayEnding
	DayStarting
)
