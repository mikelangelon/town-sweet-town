package graphics

import (
	"encoding/xml"
	"math"
)

type TileSet struct {
	XMLName    xml.Name `xml:"tileset"`
	Name       string   `xml:"name,attr"`
	TileWidth  int      `xml:"tilewidth,attr"`
	TileHeight int      `xml:"tileheight,attr"`
	Spacing    int      `xml:"spacing,attr"`
	TileCount  int      `xml:"tilecount,attr"`
	Columns    int      `xml:"columns,attr"`
	Image      Image    `xml:"image"`
	Tiles      []Tile   `xml:"tile"`
}

type Image struct {
	XMLName xml.Name `xml:"image"`
	Source  string   `xml:"source"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
}

type Tile struct {
	XMLName    xml.Name   `xml:"tile"`
	Id         int        `xml:"id,attr"`
	Type       string     `xml:"type,attr"`
	Properties Properties `xml:"properties>property"`
}

type Properties []Property
type Property struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

func (t TileSet) Width() int {
	return t.Columns
}

func (t TileSet) Height() int {
	return int(math.Ceil(float64(t.TileCount) / float64(t.Columns)))
}

func (tp Properties) HasProperty(property string) *string {
	for _, v := range tp {
		if v.Name == property {
			return &v.Value
		}
	}
	return nil
}

func (tp Properties) HasPropertyAs(property, value string) bool {
	if got := tp.HasProperty(property); got != nil {
		return value == *got
	}
	return false
}
