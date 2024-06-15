package graphics

import "encoding/xml"

type Map struct {
	XMLName     xml.Name `xml:"map"`
	Width       int      `xml:"width,attr"`
	Height      int      `xml:"height,attr"`
	TileWidth   int      `xml:"tilewidth,attr"`
	TileHeight  int      `xml:"tileheight,attr"`
	Orientation string   `xml:"orientation,attr"`
	RenderOrder string   `xml:"renderorder,attr"`
	Layers      []Layer  `xml:"layer"`
}

type Layer struct {
	XMLName xml.Name `xml:"layer"`
	ID      string   `xml:"id,attr"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
	Data    Data     `xml:"data"`
}

// Data represents the data inside a Layer
type Data struct {
	XMLName  xml.Name `xml:"data"`
	Encoding string   `xml:"encoding,attr"`
	Raw      []byte   `xml:",innerxml"`
}

func (t Map) validateSupported() error {
	return nil
}
