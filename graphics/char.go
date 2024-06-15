package graphics

import "github.com/hajimehoshi/ebiten/v2"

type Char struct {
	Image *ebiten.Image
	X, Y  int64

	scaleX, scaleY float64
	Stuff          []*ebiten.Image
}

func (c *Char) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.X), float64(c.Y))
	op.GeoM.Scale(c.scaleX, c.scaleY)

	screen.DrawImage(c.Image, op)

	for _, v := range c.Stuff {
		screen.DrawImage(v, op)
	}
}

func (c *Char) Talk() string {
	return "Welcome<NEXT>I hope you like this place<NEXT>Good luck!"
}
