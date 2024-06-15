package graphics

import (
	"encoding/xml"
)

func parseMap(bytes []byte) (*Map, error) {
	var m Map
	err := xml.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	err = m.validateSupported()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func parseTileSet(bytes []byte) (*TileSet, error) {
	var t TileSet
	err := xml.Unmarshal(bytes, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
