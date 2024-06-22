package house

import (
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
)

const (
	Tend    = "tend"
	Red     = "red"
	Boring  = "boring"
	Big     = "big"
	Fashion = "fashion"
)

type House struct {
	ID    string
	Owner *string
	House graphics.MapScene
}

func (h House) DoorPosition() common.Position {
	return common.Position{
		X: h.House.Offset.X + 3*16,
		Y: h.House.Offset.Y + 3*16,
	}
}

type HouseInfo struct {
	Text string
	Cost int
	Type int
}

type MapHouseInfo map[string]HouseInfo

var MapHouseBulding MapHouseInfo = map[string]HouseInfo{
	"10 coins --> Normal House": {
		Text: "10 coins --> Normal House",
		Cost: 10,
		Type: 0,
	},
	"11 coins --> Red House": {
		Text: "11 coins --> Red House",
		Cost: 11,
		Type: 1,
	},
	"15 coins --> Big House": {
		Text: "15 coins --> Big House",
		Cost: 15,
		Type: 2,
	},
	"8 coins --> Tend": {
		Text: "8 coins --> Tend",
		Cost: 8,
		Type: 3,
	},
	"17 coins --> Fashion House": {
		Text: "17 coins --> Fashion House",
		Cost: 17,
		Type: 4,
	},
}

func (m MapHouseInfo) GiveMeThree() []string {
	var three []string
	var counter = 0
	for k, _ := range m {
		counter++
		three = append(three, k)
		if counter >= 3 {
			break
		}
	}
	return three
}
