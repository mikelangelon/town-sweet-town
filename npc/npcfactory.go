package npc

import (
	"github.com/icrowley/fake"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
)

type NPCFactory struct {
	CharFactory *graphics.CharFactory
	Scale       int
}

func (n *NPCFactory) NewNPC(id int, withIDs []int, x, y int64, move common.Position) *NPC {
	return &NPC{
		Char: graphics.Char{
			ID:     fake.FirstName(),
			Image:  n.CharFactory.CharImage(id),
			X:      x,
			Y:      y,
			ScaleX: float64(n.Scale),
			ScaleY: float64(n.Scale),
			Stuff:  n.CharFactory.CharImages(withIDs),
		},
		Move: move,
	}
}
