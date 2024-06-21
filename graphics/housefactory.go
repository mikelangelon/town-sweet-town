package graphics

import (
	"github.com/mikelangelon/town-sweet-town/assets"
	"github.com/mikelangelon/town-sweet-town/common"
)

type HouseFactory struct {
	Houses []*MapScene
}

func NewHouseFactory() (*HouseFactory, error) {
	house, err := NewMapScene(assets.TileMapPacked, assets.House, assets.TileMapPackedTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		return nil, err
	}
	houseFancy, err := NewMapScene(assets.FancyTown, assets.House4, assets.FancyTownTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		return nil, err
	}
	house2, err := NewMapScene(assets.TileMapPacked, assets.House2, assets.TileMapPackedTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		return nil, err
	}
	house3, err := NewMapScene(assets.TileMapPacked, assets.House3, assets.TileMapPackedTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		return nil, err
	}
	tend, err := NewMapScene(assets.FancyTown, assets.Tend, assets.FancyTownTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		return nil, err
	}
	return &HouseFactory{
		Houses: []*MapScene{house, house2, house3, tend, houseFancy},
	}, nil
}
