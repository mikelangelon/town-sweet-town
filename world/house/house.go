package house

import (
	"fmt"
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
	Name string
	Cost int
	Type int
}

func (h HouseInfo) Text() string {
	return fmt.Sprintf("%d coins --> %s", h.Cost, h.Name)
}

type MapHouseInfo map[string]HouseInfo

var MapHouseBulding MapHouseInfo = map[string]HouseInfo{
	"10 coins --> Normal House": {
		Name: "Normal House",
		Cost: 10,
		Type: 0,
	},
	"11 coins --> Red House": {
		Name: "Red House",
		Cost: 11,
		Type: 1,
	},
	"15 coins --> Big House": {
		Name: "Big House",
		Cost: 15,
		Type: 2,
	},
	"8 coins --> Tend": {
		Name: "Tend",
		Cost: 8,
		Type: 3,
	},
	"17 coins --> Fashion House": {
		Name: "Fashion House",
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
