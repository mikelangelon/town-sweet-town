package scenes

import "github.com/mikelangelon/town-sweet-town/graphics"

type Entrance struct {
	BaseScene
}

func NewEntrance(id string, mapScene *graphics.MapScene) *Entrance {
	return &Entrance{
		BaseScene{ID: id, MapScene: mapScene},
	}
}
