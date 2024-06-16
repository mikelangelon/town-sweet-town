package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Char struct {
	ID    string
	Image *ebiten.Image
	X, Y  int64

	ScaleX, ScaleY float64
	Stuff          []*ebiten.Image
}

func (c *Char) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.X), float64(c.Y))
	op.GeoM.Scale(c.ScaleX, c.ScaleY)

	screen.DrawImage(c.Image, op)

	for _, v := range c.Stuff {
		screen.DrawImage(v, op)
	}
}
