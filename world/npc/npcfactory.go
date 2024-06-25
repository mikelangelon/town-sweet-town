package npc

import (
	"github.com/icrowley/fake"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"math/rand"
)

type NPCFactory struct {
	CharFactory *graphics.CharFactory
	Scale       int
}

func (n *NPCFactory) NewNPC(id int, withIDs []int, pos common.Position, chs Chars) *NPC {
	return &NPC{
		Char: graphics.Char{
			ID:     fake.FirstName(),
			Image:  n.CharFactory.CharImage(id),
			X:      pos.X,
			Y:      pos.Y,
			ScaleX: float64(n.Scale),
			ScaleY: float64(n.Scale),
			Stuff:  n.CharFactory.CharImages(withIDs),
		},
		Move: &common.Position{
			X: pos.X - 7*16,
			Y: pos.Y,
		},
		Chars:         chs,
		NitPicky:      rand.Intn(100),
		NitPickyLevel: rand.Intn(8),
	}
}
