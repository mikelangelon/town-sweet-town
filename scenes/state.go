package scenes

import "github.com/mikelangelon/town-sweet-town/graphics"

type State struct {
	Player *graphics.Char
	Status int
}

const (
	Menu = iota
	InitialState
	Playing
	Transitioning
)
