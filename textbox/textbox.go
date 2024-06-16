package textbox

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/mikelangelon/town-sweet-town/assets"
	"image/color"
	"strings"
)

var backgroundColor = color.RGBA{50, 50, 50, 150}

type TextBox struct {
	visible  bool
	text     string
	nextText *string
	next     []string
}

const (
	boxX float64 = 100
	boxY float64 = 400
	boxW int     = 500
	boxH int     = 150
)

func (c *TextBox) Draw(screen *ebiten.Image) {
	if !c.visible {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(boxX, boxY)
	screen.DrawImage(c.drawBackground(), op)
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(boxX, boxY)
	drawText(screen, c.text, op2)

}

func (c *TextBox) drawBackground() *ebiten.Image {
	bg := ebiten.NewImage(boxW, boxH)
	bg.Fill(backgroundColor)
	return bg
}

func drawText(img *ebiten.Image, txt string, op *text.DrawOptions) (int, int) {
	font := guiFont(20)
	w, h := text.Measure(txt, font, 0)
	for _, v := range strings.Split(txt, "\n") {
		op.GeoM.Translate(10, h)
		text.Draw(img, v, font, op)
		h += 10
	}

	return int(w), int(h)
}

func (c *TextBox) Visible() bool {
	return c.visible
}

func (c *TextBox) Next() {
	if len(c.next) > 0 {
		c.Show(c.next)
		return
	}
	c.visible = false
}

func (c *TextBox) Show(text []string) {
	if len(text) == 0 {
		return
	}
	c.visible = true
	c.text = text[0]
	c.next = text[1:]
}

func guiFont(size float64) *text.GoTextFace {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(assets.HolsteinFont))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}
}
