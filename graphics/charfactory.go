package graphics

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

// CharFactory same idea as MapScene, but for characters.
// TODO: Do I need this so similar structs?
type CharFactory struct {
	tileSet *TileSet
	image   *ebiten.Image
	scale   int
}

func NewCharFactory(mapImage, tsx []byte, scale int) (*CharFactory, error) {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(mapImage))
	if err != nil {
		return nil, err
	}

	ebitenImage := ebiten.NewImageFromImage(img)

	t, err := parseTileSet(tsx)
	if err != nil {
		return nil, err
	}

	return &CharFactory{
		tileSet: t,
		image:   ebitenImage,
		scale:   scale,
	}, nil
}

// TODO is almost as tileImage. Unify both functions, or extract as a new common one
func (c *CharFactory) CharImage(id int) *ebiten.Image {
	col := id % c.tileSet.Width()
	row := id / c.tileSet.Width()

	x0 := col * (c.tileSet.TileWidth + c.tileSet.Spacing)
	y0 := row * (c.tileSet.TileHeight + c.tileSet.Spacing)
	x1 := x0 + c.tileSet.TileWidth
	y1 := y0 + c.tileSet.TileHeight

	return c.image.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)
}

func (c *CharFactory) CharImages(ids []int) []*ebiten.Image {
	var images = make([]*ebiten.Image, len(ids))
	for i, id := range ids {
		images[i] = c.CharImage(id)
	}
	return images
}

func (c *CharFactory) NewChar(id int, withIDs []int, x, y int64) *Char {
	return &Char{
		Image:  c.CharImage(id),
		X:      x,
		Y:      y,
		ScaleX: float64(c.scale),
		ScaleY: float64(c.scale),
		Stuff:  c.CharImages(withIDs),
	}
}
