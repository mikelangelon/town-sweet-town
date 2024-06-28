package scenes

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/world"
	"github.com/mikelangelon/town-sweet-town/world/house"
	"github.com/mikelangelon/town-sweet-town/world/npc"
)

type State struct {
	Player        *graphics.Char
	Status        int
	World         map[string]*SceneMap
	Stats         map[string]int
	Day           int
	Goals         []world.Goal
	GameLogic     Brainer
	StatusEnded   *int
	TownSillySong *audio.Player
	MenuSong      *audio.Player
}

type Brainer interface {
	NextDay(state State) State
	CreateHouse(id string, typ int) house.House
	GetRuler() npc.RuleApplier
	ChangePlayer(state State) State
}

type SceneMap struct {
	NPCs    npc.NPCs
	Houses  []*house.House
	Objects []world.Object
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
	InitExplanation
	NoClothes
	Playing
	CheckWishes
	Transitioning
	DayEnding
	DayStarting
	Pause
	HappyEnd
	GameOver
	GoingToMenu
)
