package npc

import (
	"fmt"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/world/house"
	"time"
)

type NPC struct {
	graphics.Char
	Move    *common.Position
	moving  bool
	Phrases []string
	Chars   Chars
	House   *house.House
	DayIn   int
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

func (n *NPC) Talk(day int) []string {
	var result = []string{
		fmt.Sprintf("My name is %s", n.ID),
	}
	for _, v := range n.Chars.AsPhrases() {
		result = append(result, v)
	}
	if n.House == nil {
		result = append(result, "Could I live in one house, please?")
	}
	if n.House != nil && n.DayIn != day {
		result = append(result, "How can I help you?")
	}
	return result
}
