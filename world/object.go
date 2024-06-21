package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mikelangelon/town-sweet-town/common"
)

type Object interface {
	Draw(screen *ebiten.Image)
	Position() common.Position
}
