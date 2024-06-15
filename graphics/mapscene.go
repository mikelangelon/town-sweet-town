package graphics

import (
	"bytes"
	"errors"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"strconv"
	"strings"
)

// MapScene is a helper that allows drawing maps using Tiled
type MapScene struct {
	Map            *Map
	TileSet        *TileSet
	Image          *ebiten.Image
	FormattedLayer []int64

	screenWidth, screenHeight, scale int
}

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

	formattedLayer, err := formattedLayer(m)
	if err != nil {
		return nil, err
	}
	return &MapScene{
		Map:            m,
		TileSet:        t,
		Image:          ebitenImage,
		FormattedLayer: formattedLayer,
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
	for i, id := range g.FormattedLayer {
		op := &ebiten.DrawImageOptions{}
		tx := (i % g.Map.Width) * g.Map.TileWidth
		ty := (i / g.Map.Width) * g.Map.TileHeight
		op.GeoM.Translate(float64(tx), float64(ty))
		op.GeoM.Scale(scaleX, scaleY)

		screen.DrawImage(g.tileImage(int(id-1)), op)
	}
}

func (g *MapScene) Scale() (float64, float64) {
	return float64(g.scale), float64(g.scale)
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
	rx := x / g.Map.TileWidth
	ry := y / g.Map.TileHeight

	i := rx + ry*g.Map.Width
	id := g.FormattedLayer[i] - 1
	for _, v := range g.TileSet.Tiles {
		if v.Id == int(id) {
			return v
		}
	}
	return Tile{}
}

func formattedLayer(tmx *Map) ([]int64, error) {
	if len(tmx.Layers) > 1 {
		return nil, errors.New("only 1 layer is supported")
	}
	var formattedLayer []int64
	if tmx.Layers[0].Data.Encoding != "csv" {
		return formattedLayer, errors.New("only csv file supported")
	}
	ids := strings.Split(string(tmx.Layers[0].Data.Raw), ",")
	for _, s := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			return formattedLayer, err
		}
		if err != nil {
			return nil, err
		}
		formattedLayer = append(formattedLayer, id)
	}
	return formattedLayer, nil
}
