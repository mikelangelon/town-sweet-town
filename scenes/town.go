package scenes

import "github.com/mikelangelon/town-sweet-town/graphics"

type Town struct {
	BaseScene
}

func NewTown(id string, mapScene *graphics.MapScene) *Town {
	return &Town{
		BaseScene{ID: id, MapScene: mapScene},
	}
}

func (t *Town) Update() error {
	if t.endOfDay != nil {
		t.endOfDay.Update()
		if t.endOfDay.done {
			t.endOfDay = nil
			t.state.Status = DayStarting
		}
	}
	return t.BaseScene.Update()
}
