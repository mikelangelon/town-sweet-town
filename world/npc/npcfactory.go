package npc

import (
	"github.com/icrowley/fake"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"math/rand"
	"slices"
)

type NPCFactory struct {
	CharFactory *graphics.CharFactory
	Scale       int
	UsedIDs     []string
}

func (n *NPCFactory) NewRandomGuy(pos common.Position, chs Chars) *NPC {
	bodyOptions := []int{0, 1, 54, 55, 108, 109, 162, 163}
	body := bodyOptions[rand.Intn(len(bodyOptions))]
	clothes := clothes()
	mainClothes := clothes[rand.Intn(len(clothes))]
	pants := pants()[rand.Intn(len(pants()))]
	accessory := others()[rand.Intn(len(others()))]
	return n.NewNPC(body, []int{mainClothes, pants, accessory}, pos, chs)
}
func (n *NPCFactory) NewNPC(id int, withIDs []int, pos common.Position, chs Chars) *NPC {
	return &NPC{
		Char: graphics.Char{
			ID:     n.generateID(id),
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

func (n *NPCFactory) Reset() {
	n.UsedIDs = []string{}
}

func (n *NPCFactory) generateID(bodyID int) string {
	name := fake.MaleFirstName()
	if slices.Contains([]int{270, 271, 378, 379, 486}, bodyID) {
		name = fake.FemaleFirstName()
	}
	if slices.Contains(n.UsedIDs, name) {
		return n.generateID(bodyID)
	}
	return name
}
func clothes() []int {
	start := []int{6, 60, 114, 168, 222, 276, 330, 384}

	var total []int
	for _, v := range start {
		for i := 0; i < 12; i++ {
			total = append(total, v+i)
		}
	}
	return total
}

func pants() []int {
	return []int{3, 57, 111, 165, 273, 327, 381, 435}
}

func others() []int {
	return []int{27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 22, 131, 26, 242, 294, 402, 460, 461, 462, 463, 366, 367, 368, 421, 532, 308, 415, 357}
}
