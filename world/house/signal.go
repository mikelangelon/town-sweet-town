package house

import (
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
)

type Signal struct {
	*graphics.Char
	HousePlace   common.Position
	House        *House
	HouseOptions *[]string

	Cost int
}

func NewSignal(char *graphics.Char, id string, housePlace common.Position) *Signal {
	s := Signal{
		Char:       char,
		HousePlace: housePlace,
	}
	s.ID = id
	return &s
}
