package graphics

import (
	"bytes"
	"errors"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mikelangelon/town-sweet-town/common"
	"image"
	"strconv"
	"strings"
)

// MapScene is a helper that allows drawing maps using Tiled
type MapScene struct {
	Map            *Map
	TileSet        *TileSet
	Image          *ebiten.Image
	FormattedLayer []formattedLayer

	screenWidth, screenHeight, scale int

	Offset common.Position
	Child  []*MapScene
}

type formattedLayer []int64

func NewMapScene(mapImage, tmx, tsx []byte, screenWidth, screenHeight, scale int) (*MapScene, error) {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(mapImage))
	if err != nil {
		return nil, err
	}

	ebitenImage := ebiten.NewImageFromImage(img)

	m, err := parseMap(tmx)
	if err != nil {
		return nil, err
	}

	t, err := parseTileSet(tsx)
	if err != nil {
		return nil, err
	}

	formattedLayers, err := formatLayers(m)
	if err != nil {
		return nil, err
	}
	return &MapScene{
		Map:            m,
		TileSet:        t,
		Image:          ebitenImage,
		FormattedLayer: formattedLayers,
		screenWidth:    screenWidth,
		screenHeight:   screenHeight,
		scale:          scale,
	}, nil
}

func (g *MapScene) Update() error {
	return nil
}

func (g *MapScene) Draw(screen *ebiten.Image) {
	scaleX, scaleY := g.Scale()
	for _, layer := range g.FormattedLayer {
		for i, id := range layer {
			op := &ebiten.DrawImageOptions{}
			tx := int64((i % g.Map.Width) * g.Map.TileWidth)
			ty := int64((i / g.Map.Width) * g.Map.TileHeight)
			op.GeoM.Translate(float64(tx+g.Offset.X), float64(ty+g.Offset.Y))
			op.GeoM.Scale(scaleX, scaleY)

			screen.DrawImage(g.tileImage(int(id-1)), op)
		}
	}
	for _, v := range g.Child {
		v.Draw(screen)
	}
}

func (g *MapScene) Scale() (float64, float64) {
	return float64(g.scale), float64(g.scale)
}

func (g *MapScene) SetOffset(offset common.Position) {
	g.Offset = offset
}

func (g *MapScene) tileImage(id int) *ebiten.Image {
	col := id % g.TileSet.Width()
	row := id / g.TileSet.Width()

	x0 := col * (g.TileSet.TileWidth + g.TileSet.Spacing)
	y0 := row * (g.TileSet.TileHeight + g.TileSet.Spacing)
	x1 := x0 + g.TileSet.TileWidth
	y1 := y0 + g.TileSet.TileHeight

	return g.Image.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)
}

func (g *MapScene) TileForPos(x, y int) Tile {
	if x < 0 || y < 0 || x > g.Map.Width*g.Map.TileWidth || y > g.Map.Height*g.Map.TileHeight {
		return Tile{}
	}
	rx := x / g.Map.TileWidth
	ry := y / g.Map.TileHeight

	i := rx + ry*g.Map.Width
	if i >= len(g.FormattedLayer[0]) {
		return Tile{}
	}
	if i < 0 {
		return Tile{}
	}
	id := g.FormattedLayer[0][i] - 1
	for _, v := range g.TileSet.Tiles {
		if v.Id == int(id) {
			return v
		}
	}
	return Tile{}
}

func (g *MapScene) AnyPropertyTileAs(x, y int, property, value string) bool {
	t := g.TileForPos(x, y)
	if t.Properties.HasPropertyAs(property, value) {
		return true
	}
	for _, v := range g.Child {
		if value := v.AnyPropertyTileAs(x-int(v.Offset.X), y-int(v.Offset.Y), property, value); value {
			return true
		}
	}
	return false
}

func formatLayers(tmx *Map) ([]formattedLayer, error) {
	var result []formattedLayer
	for _, v := range tmx.Layers {
		var formattedLayer formattedLayer
		if v.Data.Encoding != "csv" {
			return result, errors.New("only csv file supported")
		}
		ids := strings.Split(string(v.Data.Raw), ",")
		for _, s := range ids {
			id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
			if err != nil {
				return result, err
			}
			if err != nil {
				return nil, err
			}
			formattedLayer = append(formattedLayer, id)
		}
		result = append(result, formattedLayer)
	}

	return result, nil
}
