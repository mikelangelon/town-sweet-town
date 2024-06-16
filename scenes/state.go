package scenes

import (
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/npc"
)

type State struct {
	Player *graphics.Char
	Status int
	NPCs   map[string][]*npc.NPC
}

const (
	Menu = iota
	InitialState
	Playing
	Transitioning
)
