package house

import (
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
)

type House struct {
	ID    string
	Owner *string
	House *graphics.MapScene
}

func (h House) DoorPosition() common.Position {
	return common.Position{
		X: h.House.Offset.X + 3*16,
		Y: h.House.Offset.Y + 3*16,
	}
}
