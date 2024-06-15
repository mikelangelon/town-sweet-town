package graphics

import (
	"github.com/mikelangelon/town-sweet-town/common"
	"time"
)

type NPC struct {
	Char
	Move   common.Position
	moving bool
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
