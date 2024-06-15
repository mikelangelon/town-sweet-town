package scenes

import (
	"github.com/joelschutz/stagehand"
	"image"
)

type BaseScene struct {
	bounds image.Rectangle
	state  State
	sm     *stagehand.SceneManager[State]
}

func (s *BaseScene) Layout(w, h int) (int, int) {
	s.bounds = image.Rect(0, 0, w, h)
	return w, h
}

func (s *BaseScene) Load(st State, sm stagehand.SceneController[State]) {
	s.state = st
	s.sm = sm.(*stagehand.SceneManager[State])
}

func (s *BaseScene) Unload() State {
	return s.state
}
