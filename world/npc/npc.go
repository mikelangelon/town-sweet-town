package npc

import (
	"fmt"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/world/house"
	"time"
)

const groupPhases = 3

type NPC struct {
	graphics.Char
	Move          *common.Position
	moving        bool
	Chars         Chars
	House         *house.House
	DayIn         int
	alreadyMet    bool
	NitPicky      int
	NitPickyLevel int
	Wishes        []Wish
}

func (n *NPC) SetHouse(house *house.House, dayIn int) {
	n.House = house
	n.DayIn = dayIn
}

func (n *NPC) Update() error {
	if n.moving || n.Move == nil {
		return nil
	}
	if n.Char.X < n.Move.X {
		n.Char.X += 16
	} else if n.Char.X > n.Move.X {
		n.Char.X -= 16
	} else {
		n.Move = nil
	}
	n.moving = true
	timer := time.NewTimer(400 * time.Millisecond)
	go func() {
		<-timer.C
		n.moving = false
	}()
	return nil
}

func (n *NPC) Sentences() []string {
	var result []string
	if !n.alreadyMet {
		result = append(result, fmt.Sprintf("My name is %s", n.ID))
		n.alreadyMet = true
	}
	var group string
	for i, v := range n.Chars.AsPhrases() {
		group += v + "\n"
		if (i+1)%(groupPhases+1) == 0 {
			result = append(result, group)
			group = ""
			continue
		}
	}
	if len(n.Wishes) > 0 {
		result = append(result, n.Wishes[0].IamText())
	}
	if len(group) > 0 {
		result = append(result, group)
	}
	return result
}
func (n *NPC) Talk(day int) []string {
	result := n.Sentences()
	if n.House == nil {
		result = append(result, "Could I live in one house, please?")
	}
	if n.House != nil && n.DayIn != day {
		result = append(result, "How can I help you?")
	}
	return result
}

func (n *NPC) AdaptNitpicky() int {
	if n.House == nil {
		return n.NitPickyLevel
	}
	var substract int
	switch n.House.Type {
	case 0:
		substract = 0
	case 1:
		substract = 4
	case 2:
		substract = 10
	case 3:
		substract = -5
	case 4:
		substract = 20
	}
	adapted := n.NitPicky - substract
	if adapted < 0 {
		return 0
	}
	return adapted
}
