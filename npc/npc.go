package npc

import (
	"fmt"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"time"
)

type NPC struct {
	graphics.Char
	Move    common.Position
	moving  bool
	Phrases []string
	Chars   Chars
}

func (n *NPC) Update() error {
	if n.moving {
		return nil
	}
	if n.Char.X < n.Move.X {
		n.Char.X += 16
	} else if n.Char.X > n.Move.X {
		n.Char.X -= 16
	}
	n.moving = true
	timer := time.NewTimer(800 * time.Millisecond)
	go func() {
		<-timer.C
		n.moving = false
	}()
	return nil
}

func (n *NPC) Talk() []string {
	var result = []string{
		fmt.Sprintf("My name is %s", n.ID),
	}
	for _, v := range n.Chars.AsPhrases() {
		result = append(result, v)
	}
	return result
}
