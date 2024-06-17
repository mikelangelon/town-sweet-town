package textbox

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/mikelangelon/town-sweet-town/assets"
	"image/color"
	"strings"
)

var (
	backgroundColor = color.RGBA{50, 50, 50, 150}
	selectedColor   = color.RGBA{150, 150, 150, 150}
)

const (
	NoResponse = "None"
	No         = "No"
)

type TextBox struct {
	visible  bool
	text     string
	nextText *string
	// has next phases. Also, for now defined when to do the question(if any)
	next []string

	Options        []string
	SelectedOption int
	answerFunc     func(answer string)
}

const (
	boxX float64 = 100
	boxY float64 = 400
	boxW int     = 500
	boxH int     = 150
)

func (c *TextBox) Update() error {
	if !c.visible {
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if !c.hasNext() && len(c.Options) > 0 {
			answer := c.Options[c.SelectedOption]
			c.answerFunc(answer)
			c.selectDefaultAnswer()
		}
		c.Next()
	}
	if !c.hasNext() {
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			c.SelectedOption--
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			c.SelectedOption++
		}
	}
	return nil
}

func (c *TextBox) Draw(screen *ebiten.Image) {
	if !c.visible {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(boxX, boxY)
	screen.DrawImage(c.drawBackground(), op)
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(boxX, boxY)

	_, h := drawText(screen, c.text, op2)
	font := guiFont(20)
	if !c.hasNext() {
		for i, v := range c.Options {
			h += 20
			op2 := &text.DrawOptions{}
			op2.GeoM.Translate(boxX+10, boxY+float64(h))
			if i == c.SelectedOption {
				op2.ColorScale.ScaleWithColor(selectedColor)
			}
			text.Draw(screen, v, font, op2)
		}
	}
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

func (c *TextBox) hasNext() bool {
	return len(c.next) > 0
}

func (c *TextBox) Next() {
	if c.hasNext() {
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

func (c *TextBox) ShowAndQuestion(text []string, options []string, answerFunc func(string)) {
	c.Show(text)
	c.Options = options
	c.answerFunc = answerFunc
	c.selectDefaultAnswer()
}

func (c *TextBox) selectDefaultAnswer() {
	c.SelectedOption = 0
	for i, v := range c.Options {
		if v == NoResponse {
			c.SelectedOption = i
		}
	}
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
