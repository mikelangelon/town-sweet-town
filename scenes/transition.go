package scenes

import (
	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/common"
)

type Transition struct {
	Position  common.Position
	Scene     stagehand.Scene[State]
	Direction stagehand.SlideDirection
}
